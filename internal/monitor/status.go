package monitor

import (
	"repo-mon/internal/git"
	"time"
)

type RemoteInfo struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Ahead  int    `json:"ahead"`
	Behind int    `json:"behind"`
}

type RepoStatus struct {
	RepoID              uint         `json:"repoId"`
	CurrentBranch       string       `json:"currentBranch"`
	UncommittedChanges  int          `json:"uncommittedChanges"`
	UntrackedFiles      int          `json:"untrackedFiles"`
	ModifiedFiles       int          `json:"modifiedFiles"`
	StagedFiles         int          `json:"stagedFiles"`
	UnpushedCommits     int          `json:"unpushedCommits"`
	UnpulledCommits     int          `json:"unpulledCommits"`
	StashCount          int          `json:"stashCount"`
	HasConflicts        bool         `json:"hasConflicts"`
	Remotes             []RemoteInfo `json:"remotes"`
	RemoteAccessible    bool         `json:"remoteAccessible"`
	LastChecked         time.Time    `json:"lastChecked"`
	LastSuccessfulCheck time.Time    `json:"lastSuccessfulCheck"`
	Error               string       `json:"error"`
	CheckingRemote      bool         `json:"checkingRemote"`
}

// ComputeLocalStatus runs only fast local git checks (no network).
func ComputeLocalStatus(repoID uint, repoPath string, previousStatus *RepoStatus) *RepoStatus {
	status := &RepoStatus{
		RepoID:         repoID,
		LastChecked:    time.Now(),
		CheckingRemote: true,
	}

	if previousStatus != nil {
		status.LastSuccessfulCheck = previousStatus.LastSuccessfulCheck
		status.RemoteAccessible = previousStatus.RemoteAccessible
		status.UnpushedCommits = previousStatus.UnpushedCommits
		status.UnpulledCommits = previousStatus.UnpulledCommits
		status.Error = previousStatus.Error
		// Deep copy Remotes slice to avoid sharing with concurrent goroutines
		if previousStatus.Remotes != nil {
			status.Remotes = make([]RemoteInfo, len(previousStatus.Remotes))
			copy(status.Remotes, previousStatus.Remotes)
		}
	}

	repo, err := git.OpenRepo(repoPath)
	if err != nil {
		status.Error = "not a git repository or git error: " + err.Error()
		status.CheckingRemote = false
		return status
	}

	branch, err := git.CurrentBranch(repo)
	if err != nil {
		status.Error = "not a git repository or git error: " + err.Error()
		status.CheckingRemote = false
		return status
	}
	status.CurrentBranch = branch

	modified, staged, untracked, hasConflicts, err := git.WorktreeStatus(repo, repoPath)
	if err != nil {
		status.Error = "status error: " + err.Error()
		status.CheckingRemote = false
		return status
	}
	status.ModifiedFiles = modified
	status.StagedFiles = staged
	status.UntrackedFiles = untracked
	status.UncommittedChanges = modified + untracked + staged
	status.HasConflicts = hasConflicts

	stashCount, _ := git.StashCount(repoPath)
	status.StashCount = stashCount

	return status
}

// ComputeRemoteStatus runs the slow remote checks (fetch + ahead/behind).
// It merges results into the provided local status.
func ComputeRemoteStatus(status *RepoStatus, repoPath string) {
	status.CheckingRemote = false

	repo, err := git.OpenRepo(repoPath)
	if err != nil {
		status.RemoteAccessible = false
		status.Remotes = nil
		return
	}

	remoteNames, err := git.RemoteNames(repo)
	if err != nil || len(remoteNames) == 0 {
		status.RemoteAccessible = false
		status.Remotes = nil
		return
	}

	remoteAccessible := false
	totalAhead := 0
	totalBehind := 0
	var remotes []RemoteInfo

	for _, remoteName := range remoteNames {
		ri := RemoteInfo{Name: remoteName}
		url, _ := git.RemoteURL(repo, remoteName)
		ri.URL = url

		fetchErr := git.FetchRemote(repo, remoteName)
		if fetchErr == nil {
			remoteAccessible = true
			ahead, behind, abErr := git.AheadBehind(repo, status.CurrentBranch, remoteName)
			if abErr == nil {
				ri.Ahead = ahead
				ri.Behind = behind
				totalAhead += ahead
				totalBehind += behind
			}
		}

		remotes = append(remotes, ri)
	}

	status.Remotes = remotes
	status.RemoteAccessible = remoteAccessible
	status.UnpushedCommits = totalAhead
	status.UnpulledCommits = totalBehind

	if remoteAccessible {
		status.LastSuccessfulCheck = time.Now()
		status.Error = ""
	} else {
		status.Error = "remote unreachable"
	}
}
