package jobs

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sonda2208/grc/jobs/rules"

	"github.com/sonda2208/grc/model"
	"github.com/sonda2208/grc/store"
	"github.com/sonda2208/grc/util"
)

const (
	DefaultPollingInterval = 3000 // 3s

	LocalGitDirectory = "/tmp/gr/"

	numberOfWorkers = 8
)

type Rule interface {
	Name() string
	Metadata() model.FindingMetadata
	Scan(filePath string, fileBuf []byte) ([]*model.Finding, error)
}

type JobServer struct {
	s store.Store

	stop            chan bool
	stopped         chan bool
	pollingInterval int

	startOnce sync.Once
	locker    *util.NamedMutex

	rules map[string]Rule
}

func NewJobServer(s store.Store) (*JobServer, error) {
	js := &JobServer{
		s:               s,
		stop:            make(chan bool, 1),
		stopped:         make(chan bool, 1),
		pollingInterval: DefaultPollingInterval,
		locker:          util.NewNamedMutex(),
		rules:           map[string]Rule{},
		startOnce:       sync.Once{},
	}

	secretKeyRule, err := rules.NewSecretKeyRule()
	if err != nil {
		return nil, err
	}

	js.rules[secretKeyRule.Name()] = secretKeyRule
	return js, nil
}

func (js *JobServer) Start() {
	// random start time to avoid race condition if we run multiple instances
	rand.Seed(time.Now().UTC().UnixNano())
	<-time.After(time.Duration(rand.Intn(js.pollingInterval)) * time.Millisecond)

	js.startOnce.Do(func() {
		go js.startWatcher()
	})
}

func (js *JobServer) Stop() {
	close(js.stop)
	<-js.stopped
}

func (js *JobServer) pollAndNotify() {
	scans, err := js.s.Scan().GetByStatus(model.ScanStatusQueued)
	if err != nil {
		log.Fatalf("Error occurred when getting pending scans")
		return
	}

	for _, s := range scans {
		go js.doScan(s)
	}
}

func (js *JobServer) startWatcher() {
	defer func() {
		close(js.stopped)
	}()

	for {
		select {
		case <-js.stop:
			return
		case <-time.After(time.Duration(js.pollingInterval) * time.Millisecond):
			js.pollAndNotify()
		}
	}
}

func (js *JobServer) getScanByID(id int) (*model.Scan, error) {
	scan, err := js.s.Scan().Get(id)
	if err != nil {
		return nil, err
	}

	return scan, nil
}

func (js *JobServer) updateScan(s *model.Scan) error {
	err := js.s.Scan().Update(s)
	if err != nil {
		return err
	}

	return nil
}

func (js *JobServer) SetScanStatus(scanID int, status string, message string, findings ...*model.Finding) error {
	s, err := js.getScanByID(scanID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	switch status {
	case model.ScanStatusQueued:
		s.Findings = nil
	case model.ScanStatusInProgress:
		s.ScanningAt = &now
	case model.ScanStatusSuccess:
		s.Findings = findings
		fallthrough
	case model.ScanStatusFailure:
		s.FinishedAt = &now
	}

	s.Status = status
	s.Message = message
	err = js.updateScan(s)
	if err != nil {
		return err
	}

	return nil
}

func (js *JobServer) GetRepoByID(id int) (*model.Repository, error) {
	repo, err := js.s.Repository().Get(id)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (js *JobServer) doScan(s *model.Scan) {
	var err error
	defer func() {
		if err != nil {
			err := js.SetScanStatus(s.ID, model.ScanStatusFailure, err.Error())
			if err != nil {
				log.Println(err)
			}
		}
	}()

	err = js.SetScanStatus(s.ID, model.ScanStatusInProgress, "")
	if err != nil {
		return
	}

	r, err := js.GetRepoByID(s.RepoID)
	if err != nil {
		return
	}

	// only run 1 scan for a repo at same time
	js.locker.Lock(r.Name)
	defer js.locker.Unlock(r.Name)

	repoPath := path.Join(LocalGitDirectory, r.Name)
	err = js.fetchRepository(repoPath, r.URL, s)
	if err != nil {
		return
	}

	done := make(chan struct{})
	defer close(done)

	paths, walkErr := js.walk(done, repoPath)
	results := make(chan *model.ScanResult)
	var wg sync.WaitGroup
	wg.Add(numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		go func() {
			js.scanner(done, paths, results)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var findings []*model.Finding
	for res := range results {
		findings = append(findings, res.Findings...)
	}

	// remove local repo directory prefix
	for _, f := range findings {
		f.Location.Path = strings.TrimPrefix(f.Location.Path, LocalGitDirectory)
	}

	if err := <-walkErr; err != nil {
		return
	}

	err = js.SetScanStatus(s.ID, model.ScanStatusSuccess, "", findings...)
	if err != nil {
		return
	}
}

func (js *JobServer) fetchRepository(repoPath string, url string, s *model.Scan) error {
	//var (
	//	repo *git.Repository
	//	err  error
	//)
	//
	//cloneOpt := &git.CloneOptions{
	//	URL: url,
	//}
	//if s.Branch != "" {
	//	cloneOpt.ReferenceName = plumbing.NewBranchReferenceName(s.Branch)
	//}
	//
	//repo, err = git.PlainClone(repoPath, false, cloneOpt)
	//if err != nil {
	//	if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
	//		return err
	//	}
	//
	//	repo, err = git.PlainOpen(repoPath)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//wt, err := repo.Worktree()
	//if err != nil {
	//	return err
	//}
	//
	//err = wt.Checkout(&git.CheckoutOptions{
	//	Hash:  plumbing.NewHash(s.Commit),
	//	Force: true,
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}

func (js *JobServer) walk(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(paths)
		errs <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			select {
			case paths <- path:
			case <-done:
				return model.NewError("scanner.file_walk_canceled", "walk canceled")
			}
			return nil
		})
	}()

	return paths, errs
}

func (js *JobServer) scanner(done <-chan struct{}, paths <-chan string, results chan<- *model.ScanResult) {
	for p := range paths {
		fileBuf, err := ioutil.ReadFile(p)
		if err != nil {
			continue
		}

		for _, r := range js.rules {
			res, err := r.Scan(p, fileBuf)
			select {
			case results <- &model.ScanResult{
				Findings: res,
				Error:    err,
			}:
			case <-done:
				return
			}
		}
	}
}
