package secret_key_scanner

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sonda2208/guardrails-challenge/util"

	"github.com/sonda2208/guardrails-challenge/model"
	"github.com/sonda2208/guardrails-challenge/scanners"
)

const (
	SecretKeysPattern = `(?i)\s*(\bBEGIN\b).*((PRIVATE KEY)|(PUBLIC KEY)\b)\s*`

	numberOfWorkers = 8
)

type SecretKeyScanner struct {
	stop     chan bool
	stopped  chan bool
	scans    chan model.Scan
	server   *scanners.JobServer
	re       *regexp.Regexp
	locker   *util.NamedMutex
	metadata model.FindingMetadata
}

func NewWorker(s *scanners.JobServer) (scanners.Scanner, error) {
	re, err := regexp.Compile(SecretKeysPattern)
	if err != nil {
		return nil, err
	}

	return &SecretKeyScanner{
		stop:    make(chan bool, 1),
		stopped: make(chan bool, 1),
		scans:   make(chan model.Scan),
		server:  s,
		re:      re,
		locker:  util.NewNamedMutex(),
		metadata: model.FindingMetadata{
			Description: "Use of secret keys",
			Severity:    "HIGH",
		},
	}, nil
}

func (w SecretKeyScanner) Run() {
	defer func() {
		w.stopped <- true
	}()

	for {
		select {
		case <-w.stop:
			return
		case job := <-w.scans:
			go w.doScan(&job)
		}
	}
}

func (w SecretKeyScanner) Stop() {
	w.stop <- true
	<-w.stopped
}

func (w SecretKeyScanner) Channel() chan<- model.Scan {
	return w.scans
}

func (w SecretKeyScanner) Name() string {
	return model.ScanTypeSecretKey
}

func (w SecretKeyScanner) doScan(s *model.Scan) {
	var err error
	defer func() {
		if err != nil {
			err := w.server.SetScanStatus(s.ID, model.ScanStatusFailure, err.Error())
			if err != nil {
				log.Println(err)
			}
		}
	}()

	err = w.server.SetScanStatus(s.ID, model.ScanStatusInProgress, "")
	if err != nil {
		return
	}

	r, err := w.server.GetRepoByID(s.RepoID)
	if err != nil {
		return
	}

	w.locker.Lock(r.Name)
	defer w.locker.Unlock(r.Name)

	repoPath := path.Join(scanners.LocalGitDirectory, r.Name)
	err = w.fetchRepository(repoPath, r.URL, s)
	if err != nil {
		return
	}

	done := make(chan struct{})
	defer close(done)

	paths, errc := w.walk(done, repoPath)
	findings := make(chan *model.Finding)
	var wg sync.WaitGroup
	wg.Add(numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		go func() {
			w.scanner(done, paths, findings)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(findings)
	}()

	var res []*model.Finding
	for f := range findings {
		// remove local repo directory
		f.Location.Path = strings.TrimPrefix(f.Location.Path, scanners.LocalGitDirectory)

		res = append(res, f)
	}

	if err := <-errc; err != nil {
		return
	}

	err = w.server.SetScanStatus(s.ID, model.ScanStatusSuccess, "", res...)
	if err != nil {
		return
	}
}

func (w SecretKeyScanner) fetchRepository(repoPath string, url string, s *model.Scan) error {
	var (
		repo *git.Repository
		err  error
	)

	cloneOpt := &git.CloneOptions{
		URL: url,
	}
	if s.Branch != "" {
		cloneOpt.ReferenceName = plumbing.NewBranchReferenceName(s.Branch)
	}

	repo, err = git.PlainClone(repoPath, false, cloneOpt)
	if err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return err
		}

		repo, err = git.PlainOpen(repoPath)
		if err != nil {
			return err
		}
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(s.Commit),
		Force: true,
	})
	if err != nil {
		return err
	}

	return nil
}

func (w SecretKeyScanner) walk(done <-chan struct{}, root string) (<-chan string, <-chan error) {
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
				return model.NewError("secret_key_scanner.file_walk_canceled", "walk canceled")
			}
			return nil
		})
	}()

	return paths, errs
}

func (w SecretKeyScanner) scanner(done <-chan struct{}, paths <-chan string, findings chan<- *model.Finding) {
	for path := range paths {
		res, err := scanners.RegexSearch(w.re, path)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, r := range res {
			r.Metadata = w.metadata
			select {
			case findings <- r:
			case <-done:
				return
			}
		}
	}
}
