package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("git init: %v", err)
	}
	f := filepath.Join(dir, "README.md")
	os.WriteFile(f, []byte("# test"), 0644)
	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("worktree: %v", err)
	}
	if _, err := w.Add("README.md"); err != nil {
		t.Fatalf("add: %v", err)
	}
	if _, err := w.Commit("init", &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@test.com",
			When:  time.Now(),
		},
	}); err != nil {
		t.Fatalf("commit: %v", err)
	}
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

	r, _ := gogit.PlainOpen(repo)
	wt, _ := r.Worktree()
	wt.Add("new.txt")
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
