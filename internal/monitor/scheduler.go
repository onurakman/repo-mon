package monitor

import (
	"encoding/json"
	"sync"
	"time"
)

// StatusSaver is called after each status update to persist to DB.
type StatusSaver func(repoID uint, statusJSON string)

// EventEmitter is called to notify frontend of status changes.
type EventEmitter func(event string, repoID uint)

type repoTicker struct {
	ticker *time.Ticker
	stop   chan struct{}
}

type Scheduler struct {
	mu           sync.Mutex
	tickers      map[uint]*repoTicker
	statuses     sync.Map
	statusSaver  StatusSaver
	eventEmitter EventEmitter
}

func NewScheduler(saver StatusSaver, emitter EventEmitter) *Scheduler {
	return &Scheduler{
		tickers:      make(map[uint]*repoTicker),
		statusSaver:  saver,
		eventEmitter: emitter,
	}
}

// LoadCachedStatus loads a previously persisted status into memory.
func (s *Scheduler) LoadCachedStatus(repoID uint, statusJSON string) {
	if statusJSON == "" {
		return
	}
	var status RepoStatus
	if err := json.Unmarshal([]byte(statusJSON), &status); err != nil {
		return
	}
	s.statuses.Store(repoID, &status)
}

func (s *Scheduler) persistStatus(repoID uint, status *RepoStatus) {
	if s.statusSaver == nil {
		return
	}
	data, err := json.Marshal(status)
	if err != nil {
		return
	}
	s.statusSaver(repoID, string(data))
}

// checkRepo runs local status immediately, stores it, then runs remote in background.
func (s *Scheduler) SetEventEmitter(emitter EventEmitter) {
	s.eventEmitter = emitter
}

func (s *Scheduler) emitEvent(event string, repoID uint) {
	if s.eventEmitter != nil {
		s.eventEmitter(event, repoID)
	}
}

func (s *Scheduler) checkRepo(repoID uint, repoPath string) {
	prev := s.GetStatus(repoID)
	status := ComputeLocalStatus(repoID, repoPath, prev)
	s.statuses.Store(repoID, status)
	s.persistStatus(repoID, status)

	// If local check already failed or no remotes, skip remote and emit done
	if !status.CheckingRemote {
		s.emitEvent("repo:checked", repoID)
		return
	}

	s.emitEvent("repo:checking", repoID)

	// Remote check in background
	go func() {
		ComputeRemoteStatus(status, repoPath)
		s.statuses.Store(repoID, status)
		s.persistStatus(repoID, status)
		s.emitEvent("repo:checked", repoID)
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
