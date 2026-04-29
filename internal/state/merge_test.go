package state

import (
	"testing"
	"time"

	"github.com/mesa-studio/codemod-tally/internal/detector"
)

func makeMatch(file string, line int, context []string) detector.Match {
	return detector.Match{File: file, Line: line, Content: "console.log()", Context: context}
}

func TestMergeNewOccurrences(t *testing.T) {
	existing := &State{Recipe: "test"}
	ctx := []string{"before", "console.log()", "after"}
	current := []detector.Match{makeMatch("src/a.js", 5, ctx)}
	now := time.Now()

	result := Merge(existing, current, now)

	if len(result.Occurrences) != 1 {
		t.Fatalf("expected 1 occurrence, got %d", len(result.Occurrences))
	}
	if result.Occurrences[0].Status != StatusTodo {
		t.Errorf("expected todo, got %s", result.Occurrences[0].Status)
	}
}

func TestMergeMarksDone(t *testing.T) {
	ctx := []string{"before", "console.log()", "after"}
	fp := Compute("src/a.js", ctx)
	now := time.Now()

	existing := &State{
		Recipe: "test",
		Occurrences: []Occurrence{
			{Fingerprint: fp, File: "src/a.js", Line: 5, Status: StatusTodo, FirstSeen: now},
		},
	}

	result := Merge(existing, nil, now)

	if result.Occurrences[0].Status != StatusDone {
		t.Errorf("expected done, got %s", result.Occurrences[0].Status)
	}
	if result.Occurrences[0].ResolvedAt == nil {
		t.Error("expected ResolvedAt to be set")
	}
}

func TestMergePreservesAlreadyDone(t *testing.T) {
	now := time.Now()
	resolved := now.Add(-time.Hour)
	existing := &State{
		Recipe: "test",
		Occurrences: []Occurrence{
			{Fingerprint: "abc123", File: "src/a.js", Status: StatusDone, FirstSeen: now, ResolvedAt: &resolved},
		},
	}

	result := Merge(existing, nil, now)

	if result.Occurrences[0].Status != StatusDone {
		t.Errorf("expected done to be preserved, got %s", result.Occurrences[0].Status)
	}
}

func TestMergeStillPresent(t *testing.T) {
	ctx := []string{"before", "console.log()", "after"}
	match := makeMatch("src/a.js", 5, ctx)
	fp := ComputeMatch(match)
	now := time.Now()

	existing := &State{
		Recipe: "test",
		Occurrences: []Occurrence{
			{Fingerprint: fp, File: "src/a.js", Line: 5, Status: StatusTodo, FirstSeen: now},
		},
	}

	current := []detector.Match{match}
	result := Merge(existing, current, now)

	if result.Occurrences[0].Status != StatusTodo {
		t.Errorf("expected todo (still present), got %s", result.Occurrences[0].Status)
	}
}
