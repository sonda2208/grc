package scanners

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/sonda2208/guardrails-challenge/store"

	"github.com/sonda2208/guardrails-challenge/model"
)

const (
	DefaultPollingInterval = 3000 // 3s

	LocalGitDirectory = "/tmp/gr/"
)

type Scanner interface {
	Run()
	Stop()
	Channel() chan<- model.Scan
	Name() string
}

type JobServer struct {
	s store.Store

	stop            chan bool
	stopped         chan bool
	pollingInterval int

	startOnce sync.Once

	scanners map[string]Scanner
}

func NewJobServer(s store.Store) *JobServer {
	return &JobServer{
		s:               s,
		stop:            make(chan bool, 1),
		stopped:         make(chan bool, 1),
		pollingInterval: DefaultPollingInterval,
		startOnce:       sync.Once{},
		scanners:        map[string]Scanner{},
	}
}

func (js *JobServer) Start() {
	// random start time to avoid race condition if we run multiple instances
	rand.Seed(time.Now().UTC().UnixNano())
	<-time.After(time.Duration(rand.Intn(js.pollingInterval)) * time.Millisecond)

	js.startOnce.Do(func() {
		for _, w := range js.scanners {
			go w.Run()
		}

		go js.startWatcher()
	})
}

func (js *JobServer) Stop() {
	close(js.stop)
	<-js.stopped

	for _, w := range js.scanners {
		w.Stop()
	}
}

func (js *JobServer) AddScanner(w Scanner) {
	js.scanners[w.Name()] = w
}

func (js *JobServer) pollAndNotify() {
	scans, err := js.s.Scan().GetByStatus(model.ScanStatusQueued)
	if err != nil {
		log.Fatalf("Error occurred when getting pending scans")
		return
	}

	for _, s := range scans {
		w, ok := js.scanners[s.Type]
		if !ok {
			err := js.SetScanStatus(s.ID, model.ScanStatusFailure, "Unknown type")
			if err != nil {
				log.Println(err)
			}
			continue
		}

		select {
		case w.Channel() <- *s:
		default:
		}
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
