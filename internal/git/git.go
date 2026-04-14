package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func run(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func runWithTimeout(repoPath string, timeout time.Duration, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout after %s", timeout)
	}
	return strings.TrimSpace(string(out)), err
}

func IsGitRepo(path string) bool {
	_, err := run(path, "rev-parse", "--git-dir")
	return err == nil
}

func CurrentBranch(repoPath string) (string, error) {
	return run(repoPath, "rev-parse", "--abbrev-ref", "HEAD")
}

func StatusCounts(repoPath string) (modified, staged, untracked int, err error) {
	out, err := run(repoPath, "status", "--porcelain")
	if err != nil {
		return 0, 0, 0, err
	}
	if out == "" {
		return 0, 0, 0, nil
	}
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 2 {
			continue
		}
		x := line[0]
		y := line[1]
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
	out, err := run(repoPath, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 2 {
			continue
		}
		x := line[0]
		y := line[1]
		if (x == 'U' || y == 'U') || (x == 'A' && y == 'A') || (x == 'D' && y == 'D') {
			return true, nil
		}
	}
	return false, nil
}

func StashCount(repoPath string) (int, error) {
	out, err := run(repoPath, "stash", "list")
	if err != nil {
		return 0, err
	}
	if out == "" {
		return 0, nil
	}
	return len(strings.Split(out, "\n")), nil
}

func FetchRemote(repoPath string, remoteName string) error {
	_, err := runWithTimeout(repoPath, remoteTimeout, "fetch", remoteName)
	return err
}

func RemoteNames(repoPath string) ([]string, error) {
	out, err := run(repoPath, "remote")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func RemoteURL(repoPath string, remoteName string) (string, error) {
	return run(repoPath, "remote", "get-url", remoteName)
}

func AheadBehind(repoPath string, branch string, remote string) (ahead, behind int, err error) {
	upstream := remote + "/" + branch
	out, err := run(repoPath, "rev-list", "--left-right", "--count", branch+"..."+upstream)
	if err != nil {
		return 0, 0, err
	}
	parts := strings.Fields(out)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output: %s", out)
	}
	ahead, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	behind, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return ahead, behind, nil
}
