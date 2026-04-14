package git

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const remoteTimeout = 10 * time.Second

// ScanForRepos finds all git repositories under a directory (max 3 levels deep).
func ScanForRepos(root string) []string {
	var repos []string
	maxDepth := baseDepth(root) + 3

	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			return nil
		}
		// Skip hidden dirs (except .git check)
		if d.Name()[0] == '.' && d.Name() != "." {
			return filepath.SkipDir
		}
		// Limit depth
		if baseDepth(path) > maxDepth {
			return filepath.SkipDir
		}
		// Check if this dir is a git repo
		gitDir := filepath.Join(path, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			repos = append(repos, path)
			return filepath.SkipDir // Don't descend into git repos
		}
		return nil
	})
	return repos
}

func baseDepth(path string) int {
	return strings.Count(filepath.Clean(path), string(os.PathSeparator))
}

func openRepo(repoPath string) (*gogit.Repository, error) {
	return gogit.PlainOpen(repoPath)
}

func IsGitRepo(path string) bool {
	_, err := openRepo(path)
	return err == nil
}

func CurrentBranch(repoPath string) (string, error) {
	repo, err := openRepo(repoPath)
	if err != nil {
		return "", err
	}
	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	return head.Name().Short(), nil
}

func StatusCounts(repoPath string) (modified, staged, untracked int, err error) {
	repo, err := openRepo(repoPath)
	if err != nil {
		return 0, 0, 0, err
	}
	w, err := repo.Worktree()
	if err != nil {
		return 0, 0, 0, err
	}
	status, err := w.Status()
	if err != nil {
		return 0, 0, 0, err
	}
	for _, s := range status {
		x := byte(s.Staging)
		y := byte(s.Worktree)
		if x == '?' {
			untracked++
		} else {
			if x != ' ' && x != '?' {
				staged++
			}
			if y != ' ' && y != '?' {
				modified++
			}
		}
	}
	return modified, staged, untracked, nil
}

func HasConflicts(repoPath string) (bool, error) {
	repo, err := openRepo(repoPath)
	if err != nil {
		return false, err
	}
	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}
	status, err := w.Status()
	if err != nil {
		return false, err
	}
	for _, s := range status {
		x := byte(s.Staging)
		y := byte(s.Worktree)
		if (x == 'U' || y == 'U') || (x == 'A' && y == 'A') || (x == 'D' && y == 'D') {
			return true, nil
		}
	}
	return false, nil
}

func StashCount(repoPath string) (int, error) {
	stashLog := filepath.Join(repoPath, ".git", "logs", "refs", "stash")
	f, err := os.Open(stashLog)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	defer f.Close()

	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() != "" {
			count++
		}
	}
	return count, scanner.Err()
}

func FetchRemote(repoPath string, remoteName string) error {
	repo, err := openRepo(repoPath)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), remoteTimeout)
	defer cancel()

	err = repo.FetchContext(ctx, &gogit.FetchOptions{
		RemoteName: remoteName,
	})
	if err == gogit.NoErrAlreadyUpToDate {
		return nil
	}
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("timeout after %s", remoteTimeout)
	}
	return err
}

func RemoteNames(repoPath string) ([]string, error) {
	repo, err := openRepo(repoPath)
	if err != nil {
		return nil, err
	}
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, err
	}
	if len(remotes) == 0 {
		return nil, nil
	}
	names := make([]string, len(remotes))
	for i, r := range remotes {
		names[i] = r.Config().Name
	}
	return names, nil
}

func RemoteURL(repoPath string, remoteName string) (string, error) {
	repo, err := openRepo(repoPath)
	if err != nil {
		return "", err
	}
	remote, err := repo.Remote(remoteName)
	if err != nil {
		return "", err
	}
	urls := remote.Config().URLs
	if len(urls) == 0 {
		return "", fmt.Errorf("no URLs for remote %s", remoteName)
	}
	return urls[0], nil
}

func AheadBehind(repoPath string, branch string, remote string) (ahead, behind int, err error) {
	repo, err := openRepo(repoPath)
	if err != nil {
		return 0, 0, err
	}

	localRef, err := repo.Reference(plumbing.NewBranchReferenceName(branch), true)
	if err != nil {
		return 0, 0, fmt.Errorf("local branch %s: %w", branch, err)
	}

	remoteRefName := plumbing.NewRemoteReferenceName(remote, branch)
	remoteRef, err := repo.Reference(remoteRefName, true)
	if err != nil {
		return 0, 0, fmt.Errorf("remote ref %s/%s: %w", remote, branch, err)
	}

	localHash := localRef.Hash()
	remoteHash := remoteRef.Hash()

	if localHash == remoteHash {
		return 0, 0, nil
	}

	// Collect commits reachable from each side
	localSet := make(map[plumbing.Hash]struct{})
	li, err := repo.Log(&gogit.LogOptions{From: localHash})
	if err != nil {
		return 0, 0, err
	}
	li.ForEach(func(c *object.Commit) error {
		localSet[c.Hash] = struct{}{}
		return nil
	})

	remoteSet := make(map[plumbing.Hash]struct{})
	ri, err := repo.Log(&gogit.LogOptions{From: remoteHash})
	if err != nil {
		return 0, 0, err
	}
	ri.ForEach(func(c *object.Commit) error {
		remoteSet[c.Hash] = struct{}{}
		return nil
	})

	for h := range localSet {
		if _, ok := remoteSet[h]; !ok {
			ahead++
		}
	}
	for h := range remoteSet {
		if _, ok := localSet[h]; !ok {
			behind++
		}
	}

	return ahead, behind, nil
}
