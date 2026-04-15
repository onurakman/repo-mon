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
	emitterMu    sync.RWMutex
	pauseMu      sync.RWMutex
	paused       bool
	wg           sync.WaitGroup
	tickers      map[uint]*repoTicker
	statuses     sync.Map
	checking     sync.Map // per-repo guard to prevent overlapping checks
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

func (s *Scheduler) SetEventEmitter(emitter EventEmitter) {
	s.emitterMu.Lock()
	defer s.emitterMu.Unlock()
	s.eventEmitter = emitter
}

func (s *Scheduler) emitEvent(event string, repoID uint) {
	s.emitterMu.RLock()
	emitter := s.eventEmitter
	s.emitterMu.RUnlock()
	if emitter != nil {
		emitter(event, repoID)
	}
}

func (s *Scheduler) checkRepo(repoID uint, repoPath string) {
	// Skip if a check is already running for this repo
	if _, running := s.checking.LoadOrStore(repoID, true); running {
		return
	}

	prev := s.GetStatus(repoID)
	status := ComputeLocalStatus(repoID, repoPath, prev)

	// Store a snapshot copy so readers never see in-flight mutations
	localCopy := *status
	s.statuses.Store(repoID, &localCopy)
	s.persistStatus(repoID, &localCopy)

	// If local check already failed or no remotes, skip remote and emit done
	if !status.CheckingRemote {
		s.checking.Delete(repoID)
		s.emitEvent("repo:checked", repoID)
		return
	}

	s.emitEvent("repo:checking", repoID)

	// Remote check in background on an independent copy
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer s.checking.Delete(repoID)
		remoteStatus := *status
		ComputeRemoteStatus(&remoteStatus, repoPath)
		s.statuses.Store(repoID, &remoteStatus)
		s.persistStatus(repoID, &remoteStatus)
		s.emitEvent("repo:checked", repoID)
	}()
}

func (s *Scheduler) Start(repoID uint, repoPath string, intervalSec int) {
	s.stopTicker(repoID)

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

	// Polling loop — respects pause; manual Refresh bypasses this
	go func() {
		for {
			select {
			case <-rt.ticker.C:
				s.pauseMu.RLock()
				paused := s.paused
				s.pauseMu.RUnlock()
				if !paused {
					s.checkRepo(repoID, repoPath)
				}
			case <-rt.stop:
				rt.ticker.Stop()
				return
			}
		}
	}()
}

// stopTicker stops the polling ticker for a repo without clearing its cached status.
func (s *Scheduler) stopTicker(repoID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if rt, ok := s.tickers[repoID]; ok {
		close(rt.stop)
		delete(s.tickers, repoID)
	}
}

// Stop stops monitoring and removes all cached state for a repo.
func (s *Scheduler) Stop(repoID uint) {
	s.stopTicker(repoID)
	s.statuses.Delete(repoID)
}

func (s *Scheduler) StopAll() {
	s.mu.Lock()
	for id, rt := range s.tickers {
		close(rt.stop)
		delete(s.tickers, id)
	}
	s.mu.Unlock()

	// Wait for in-flight remote checks to finish
	s.wg.Wait()
}

func (s *Scheduler) SetPaused(paused bool) {
	s.pauseMu.Lock()
	defer s.pauseMu.Unlock()
	s.paused = paused
}

func (s *Scheduler) UpdateInterval(repoID uint, repoPath string, intervalSec int) {
	s.Start(repoID, repoPath, intervalSec)
}

func (s *Scheduler) Refresh(repoID uint, repoPath string) {
	go s.checkRepo(repoID, repoPath)
}

func (s *Scheduler) GetStatus(repoID uint) *RepoStatus {
	val, ok := s.statuses.Load(repoID)
	if !ok {
		return nil
	}
	orig := val.(*RepoStatus)
	cp := *orig
	if orig.Remotes != nil {
		cp.Remotes = make([]RemoteInfo, len(orig.Remotes))
		copy(cp.Remotes, orig.Remotes)
	}
	return &cp
}

func (s *Scheduler) GetAllStatuses() map[uint]*RepoStatus {
	result := make(map[uint]*RepoStatus)
	s.statuses.Range(func(key, value any) bool {
		orig := value.(*RepoStatus)
		cp := *orig
		if orig.Remotes != nil {
			cp.Remotes = make([]RemoteInfo, len(orig.Remotes))
			copy(cp.Remotes, orig.Remotes)
		}
		result[key.(uint)] = &cp
		return true
	})
	return result
}
