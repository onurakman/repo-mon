package monitor

import (
	"sync"
	"time"
)

type repoTicker struct {
	ticker *time.Ticker
	stop   chan struct{}
}

type Scheduler struct {
	mu       sync.Mutex
	tickers  map[uint]*repoTicker
	statuses sync.Map
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tickers: make(map[uint]*repoTicker),
	}
}

// checkRepo runs local status immediately, stores it, then runs remote in background.
func (s *Scheduler) checkRepo(repoID uint, repoPath string) {
	prev := s.GetStatus(repoID)
	status := ComputeLocalStatus(repoID, repoPath, prev)
	s.statuses.Store(repoID, status)

	// Remote check in background
	go func() {
		ComputeRemoteStatus(status, repoPath)
		s.statuses.Store(repoID, status)
	}()
}

func (s *Scheduler) Start(repoID uint, repoPath string, intervalSec int) {
	s.Stop(repoID)

	s.mu.Lock()
	defer s.mu.Unlock()

	interval := time.Duration(intervalSec) * time.Second
	if interval < 5*time.Second {
		interval = 5 * time.Second
	}

	rt := &repoTicker{
		ticker: time.NewTicker(interval),
		stop:   make(chan struct{}),
	}
	s.tickers[repoID] = rt

	// Initial check
	go s.checkRepo(repoID, repoPath)

	// Polling loop
	go func() {
		for {
			select {
			case <-rt.ticker.C:
				s.checkRepo(repoID, repoPath)
			case <-rt.stop:
				rt.ticker.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) Stop(repoID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if rt, ok := s.tickers[repoID]; ok {
		close(rt.stop)
		delete(s.tickers, repoID)
	}
	s.statuses.Delete(repoID)
}

func (s *Scheduler) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, rt := range s.tickers {
		close(rt.stop)
		delete(s.tickers, id)
	}
}

func (s *Scheduler) UpdateInterval(repoID uint, repoPath string, intervalSec int) {
	s.Start(repoID, repoPath, intervalSec)
}

func (s *Scheduler) Refresh(repoID uint, repoPath string) {
	s.checkRepo(repoID, repoPath)
}

func (s *Scheduler) GetStatus(repoID uint) *RepoStatus {
	val, ok := s.statuses.Load(repoID)
	if !ok {
		return nil
	}
	return val.(*RepoStatus)
}

func (s *Scheduler) GetAllStatuses() map[uint]*RepoStatus {
	result := make(map[uint]*RepoStatus)
	s.statuses.Range(func(key, value any) bool {
		result[key.(uint)] = value.(*RepoStatus)
		return true
	})
	return result
}
