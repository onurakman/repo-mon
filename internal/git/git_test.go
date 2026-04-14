package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("setup %v failed: %s %v", args, out, err)
		}
	}
	f := filepath.Join(dir, "README.md")
	os.WriteFile(f, []byte("# test"), 0644)
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "init")
	cmd.Dir = dir
	cmd.Run()
	return dir
}

func TestIsGitRepo(t *testing.T) {
	repo := setupTestRepo(t)
	if !IsGitRepo(repo) {
		t.Error("expected true for git repo")
	}
	if IsGitRepo(t.TempDir()) {
		t.Error("expected false for non-git dir")
	}
}

func TestCurrentBranch(t *testing.T) {
	repo := setupTestRepo(t)
	branch, err := CurrentBranch(repo)
	if err != nil {
		t.Fatal(err)
	}
	if branch != "master" && branch != "main" {
		t.Errorf("unexpected branch: %s", branch)
	}
}

func TestStatusCounts(t *testing.T) {
	repo := setupTestRepo(t)

	mod, staged, untracked, err := StatusCounts(repo)
	if err != nil {
		t.Fatal(err)
	}
	if mod != 0 || staged != 0 || untracked != 0 {
		t.Errorf("expected clean, got mod=%d staged=%d untracked=%d", mod, staged, untracked)
	}

	os.WriteFile(filepath.Join(repo, "new.txt"), []byte("new"), 0644)
	mod, staged, untracked, err = StatusCounts(repo)
	if err != nil {
		t.Fatal(err)
	}
	if untracked != 1 {
		t.Errorf("expected 1 untracked, got %d", untracked)
	}

	cmd := exec.Command("git", "add", "new.txt")
	cmd.Dir = repo
	cmd.Run()
	mod, staged, untracked, err = StatusCounts(repo)
	if err != nil {
		t.Fatal(err)
	}
	if staged != 1 {
		t.Errorf("expected 1 staged, got %d", staged)
	}
}

func TestStashCount(t *testing.T) {
	repo := setupTestRepo(t)
	count, err := StashCount(repo)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0 stashes, got %d", count)
	}
}

func TestHasConflicts(t *testing.T) {
	repo := setupTestRepo(t)
	conflicts, err := HasConflicts(repo)
	if err != nil {
		t.Fatal(err)
	}
	if conflicts {
		t.Error("expected no conflicts")
	}
}
