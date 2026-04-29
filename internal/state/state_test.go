package state

import (
	"path/filepath"
	"testing"
	"time"
)

func TestLoadSaveRoundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".scan-cache.json")

	now := time.Now().Truncate(time.Second)
	s := &State{
		Recipe:   "console-to-logger",
		RepoRoot: "/home/user/project",
		LastScan: now,
		Occurrences: []Occurrence{
			{
				Fingerprint: "abc123def456abcd",
				File:        "src/index.js",
				Line:        5,
				Content:     "  console.log('hello');",
				Status:      StatusTodo,
				FirstSeen:   now,
			},
		},
	}

	if err := Save(path, s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Recipe != s.Recipe {
		t.Errorf("Recipe mismatch: %s vs %s", loaded.Recipe, s.Recipe)
	}
	if len(loaded.Occurrences) != 1 {
		t.Fatalf("expected 1 occurrence, got %d", len(loaded.Occurrences))
	}
	if loaded.Occurrences[0].File != "src/index.js" {
		t.Errorf("File mismatch: %s", loaded.Occurrences[0].File)
	}
	if loaded.Occurrences[0].Status != StatusTodo {
		t.Errorf("Status mismatch: %s", loaded.Occurrences[0].Status)
	}
}

func TestLoadMissing(t *testing.T) {
	s, err := Load("/nonexistent/.scan-cache.json")
	if err != nil {
		t.Fatalf("Load of missing file should return empty state, got error: %v", err)
	}
	if s == nil {
		t.Fatal("expected empty state, got nil")
	}
	if len(s.Occurrences) != 0 {
		t.Errorf("expected 0 occurrences, got %d", len(s.Occurrences))
	}
}
