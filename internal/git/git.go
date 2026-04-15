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
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const remoteTimeout = 10 * time.Second

// skipDirs contains directory names that never hold git repo roots.
var skipDirs = map[string]bool{
	"node_modules": true, "vendor": true, "__pycache__": true,
	".cache": true, "dist": true, "build": true, "target": true,
	".gradle": true, ".idea": true, ".vscode": true,
}

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
		name := d.Name()
		// Skip hidden dirs (except .git check)
		if len(name) > 0 && name[0] == '.' && name != "." {
			return filepath.SkipDir
		}
		// Skip known non-repo directories
		if skipDirs[name] {
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

// OpenRepo opens a git repository at the given path. Callers should reuse the
// returned handle across multiple operations to avoid redundant I/O.
func OpenRepo(repoPath string) (*gogit.Repository, error) {
	return gogit.PlainOpen(repoPath)
}

func IsGitRepo(path string) bool {
	_, err := OpenRepo(path)
	return err == nil
}

func CurrentBranch(repo *gogit.Repository) (string, error) {
	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	return head.Name().Short(), nil
}

// WorktreeStatus computes file counts and conflict status in a single w.Status() call.
// repoPath is needed to load .git/info/exclude and global gitignore patterns
// that go-git does not load by default.
func WorktreeStatus(repo *gogit.Repository, repoPath string) (modified, staged, untracked int, hasConflicts bool, err error) {
	w, err := repo.Worktree()
	if err != nil {
		return 0, 0, 0, false, err
	}
	// Load exclude patterns that go-git misses (global gitignore, .git/info/exclude)
	w.Excludes = append(w.Excludes, loadExcludePatterns(repoPath)...)
	status, err := w.Status()
	if err != nil {
		return 0, 0, 0, false, err
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
		if (x == 'U' || y == 'U') || (x == 'A' && y == 'A') || (x == 'D' && y == 'D') {
			hasConflicts = true
		}
	}
	return modified, staged, untracked, hasConflicts, nil
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

func FetchRemote(repo *gogit.Repository, remoteName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), remoteTimeout)
	defer cancel()

	err := repo.FetchContext(ctx, &gogit.FetchOptions{
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

func RemoteNames(repo *gogit.Repository) ([]string, error) {
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

func RemoteURL(repo *gogit.Repository, remoteName string) (string, error) {
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

// maxCommitWalk caps commit iteration to prevent OOM on huge repos.
const maxCommitWalk = 10000

func AheadBehind(repo *gogit.Repository, branch string, remote string) (ahead, behind int, err error) {
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

	localCommit, err := object.GetCommit(repo.Storer, localHash)
	if err != nil {
		return 0, 0, err
	}
	remoteCommit, err := object.GetCommit(repo.Storer, remoteHash)
	if err != nil {
		return 0, 0, err
	}

	// Find merge base to limit walk scope
	bases, err := localCommit.MergeBase(remoteCommit)
	stopAt := plumbing.ZeroHash
	if err == nil && len(bases) > 0 {
		stopAt = bases[0].Hash
	}

	ahead = countCommitsUntil(repo, localHash, stopAt)
	behind = countCommitsUntil(repo, remoteHash, stopAt)
	return ahead, behind, nil
}

// countCommitsUntil walks from `from` counting commits until it reaches `stopAt` or the cap.
func countCommitsUntil(repo *gogit.Repository, from, stopAt plumbing.Hash) int {
	iter, err := repo.Log(&gogit.LogOptions{From: from})
	if err != nil {
		return 0
	}
	defer iter.Close()
	count := 0
	for {
		c, err := iter.Next()
		if err != nil {
			break
		}
		if c.Hash == stopAt {
			break
		}
		count++
		if count >= maxCommitWalk {
			break
		}
	}
	return count
}

// loadExcludePatterns loads gitignore patterns from .git/info/exclude and
// global gitignore files that go-git does not read by default.
func loadExcludePatterns(repoPath string) []gitignore.Pattern {
	var patterns []gitignore.Pattern

	// .git/info/exclude
	patterns = append(patterns, readPatternsFromFile(filepath.Join(repoPath, ".git", "info", "exclude"))...)

	// Global gitignore: prefer core.excludesFile from git config, then fallbacks
	home, err := os.UserHomeDir()
	if err != nil {
		return patterns
	}

	if configPath := readCoreExcludesFile(home); configPath != "" {
		patterns = append(patterns, readPatternsFromFile(configPath)...)
	} else {
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg == "" {
			xdg = filepath.Join(home, ".config")
		}
		fallbacks := []string{
			filepath.Join(xdg, "git", "ignore"),
			filepath.Join(home, ".gitignore_global"),
			filepath.Join(home, ".gitignore"),
		}
		for _, p := range fallbacks {
			patterns = append(patterns, readPatternsFromFile(p)...)
		}
	}

	return patterns
}

// readCoreExcludesFile parses ~/.gitconfig and XDG git config for core.excludesFile.
func readCoreExcludesFile(home string) string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		xdg = filepath.Join(home, ".config")
	}
	candidates := []string{
		filepath.Join(xdg, "git", "config"),
		filepath.Join(home, ".gitconfig"),
	}
	for _, cfg := range candidates {
		if path := parseExcludesFileFromConfig(cfg); path != "" {
			// Expand ~ prefix
			if strings.HasPrefix(path, "~/") {
				path = filepath.Join(home, path[2:])
			}
			return path
		}
	}
	return ""
}

// parseExcludesFileFromConfig does a minimal parse of a git config file for core.excludesFile.
func parseExcludesFileFromConfig(configPath string) string {
	f, err := os.Open(configPath)
	if err != nil {
		return ""
	}
	defer f.Close()

	inCore := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			inCore = strings.EqualFold(strings.TrimRight(line, "]"), "[core")
			continue
		}
		if inCore {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[0]) == "excludesFile" || strings.TrimSpace(parts[0]) == "excludesfile" {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func readPatternsFromFile(path string) []gitignore.Pattern {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var patterns []gitignore.Pattern
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, gitignore.ParsePattern(line, nil))
	}
	return patterns
}
