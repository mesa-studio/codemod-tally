package state

import (
	"os/exec"
	"testing"
)

func TestRepoIDFromGitRoot(t *testing.T) {
	dir := t.TempDir()
	exec.Command("git", "init", dir).Run()

	id, err := RepoID(dir)
	if err != nil {
		t.Fatalf("RepoID failed: %v", err)
	}
	if len(id) != 8 {
		t.Errorf("expected 8-char ID, got %d chars: %s", len(id), id)
	}
}

func TestRepoIDStable(t *testing.T) {
	dir := t.TempDir()
	exec.Command("git", "init", dir).Run()

	id1, _ := RepoID(dir)
	id2, _ := RepoID(dir)
	if id1 != id2 {
		t.Errorf("RepoID must be stable: %s vs %s", id1, id2)
	}
}

func TestRepoIDNoGit(t *testing.T) {
	dir := t.TempDir()
	id, err := RepoID(dir)
	if err != nil {
		t.Fatalf("RepoID failed for non-git dir: %v", err)
	}
	if len(id) != 8 {
		t.Errorf("expected 8-char ID, got %s", id)
	}
}
