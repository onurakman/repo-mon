package git

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const remoteTimeout = 10 * time.Second

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
