package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mesa-studio/codemod-tally/internal/detector"
	"github.com/mesa-studio/codemod-tally/internal/state"
)

type fakeDetector struct {
	matches []detector.Match
}

func (f *fakeDetector) Run(_ string) ([]detector.Match, error) {
	return f.matches, nil
}

func TestScanWritesProgress(t *testing.T) {
	dir := t.TempDir()
	jsFile := filepath.Join(dir, "index.js")
	os.WriteFile(jsFile, []byte("function foo() {\n  console.log('hello');\n}\n"), 0644)

	stateDir := t.TempDir()

	d := &fakeDetector{
		matches: []detector.Match{
			{File: jsFile, Line: 2, Content: "  console.log('hello');"},
		},
	}

	cfg := &ScanConfig{
		RecipeName: "test-recipe",
		RepoRoot:   dir,
		StateDir:   stateDir,
		Detector:   d,
	}

	result, err := Scan(cfg)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Total != 1 {
		t.Errorf("expected 1 total, got %d", result.Total)
	}
	if result.Remaining != 1 {
		t.Errorf("expected 1 remaining, got %d", result.Remaining)
	}

	progressPath := filepath.Join(stateDir, "test-recipe", "progress.md")
	data, err := os.ReadFile(progressPath)
	if err != nil {
		t.Fatalf("progress.md not written: %v", err)
	}
	if len(data) == 0 {
		t.Error("progress.md is empty")
	}
}

func TestScanSecondRunMarksDone(t *testing.T) {
	dir := t.TempDir()
	jsFile := filepath.Join(dir, "index.js")
	os.WriteFile(jsFile, []byte("function foo() {\n  console.log('hello');\n}\n"), 0644)
	stateDir := t.TempDir()

	// First scan — finds one match
	d1 := &fakeDetector{
		matches: []detector.Match{
			{File: jsFile, Line: 2, Content: "  console.log('hello');"},
		},
	}
	Scan(&ScanConfig{RecipeName: "test-recipe", RepoRoot: dir, StateDir: stateDir, Detector: d1})

	// Second scan — finds nothing (pattern gone)
	d2 := &fakeDetector{matches: nil}
	result, err := Scan(&ScanConfig{RecipeName: "test-recipe", RepoRoot: dir, StateDir: stateDir, Detector: d2})
	if err != nil {
		t.Fatalf("second Scan failed: %v", err)
	}

	if result.Done != 1 {
		t.Errorf("expected 1 done, got %d", result.Done)
	}
	if result.Remaining != 0 {
		t.Errorf("expected 0 remaining, got %d", result.Remaining)
	}
}

func TestScanKeepsAdjacentMatchesSeparate(t *testing.T) {
	dir := t.TempDir()
	jsFile := filepath.Join(dir, "index.js")
	os.WriteFile(jsFile, []byte("function boot() {\n  console.log('boot');\n  console.log('ready');\n}\n"), 0644)
	stateDir := t.TempDir()

	d := &fakeDetector{
		matches: []detector.Match{
			{File: jsFile, Line: 2, Content: "  console.log('boot');"},
			{File: jsFile, Line: 3, Content: "  console.log('ready');"},
		},
	}

	result, err := Scan(&ScanConfig{RecipeName: "test-recipe", RepoRoot: dir, StateDir: stateDir, Detector: d})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Total != 2 {
		t.Fatalf("expected 2 total, got %d", result.Total)
	}
	if result.Remaining != 2 {
		t.Errorf("expected 2 remaining, got %d", result.Remaining)
	}
}

func TestScanFiltersIncludedGlobs(t *testing.T) {
	dir := t.TempDir()
	jsFile := filepath.Join(dir, "index.js")
	mdFile := filepath.Join(dir, "README.md")
	os.WriteFile(jsFile, []byte("console.log('app');\n"), 0644)
	os.WriteFile(mdFile, []byte("docs mention console.log('doc');\n"), 0644)
	stateDir := t.TempDir()

	d := &fakeDetector{
		matches: []detector.Match{
			{File: jsFile, Line: 1, Content: "console.log('app');"},
			{File: mdFile, Line: 1, Content: "docs mention console.log('doc');"},
		},
	}

	result, err := Scan(&ScanConfig{
		RecipeName:   "test-recipe",
		RepoRoot:     dir,
		StateDir:     stateDir,
		Detector:     d,
		IncludeGlobs: []string{"**/*.js"},
	})
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("expected 1 total, got %d", result.Total)
	}

	s, err := state.Load(filepath.Join(stateDir, "test-recipe", ".scan-cache.json"))
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	if len(s.Occurrences) != 1 {
		t.Fatalf("expected 1 occurrence, got %d", len(s.Occurrences))
	}
	if s.Occurrences[0].File != jsFile {
		t.Errorf("expected only js file, got %s", s.Occurrences[0].File)
	}
}

func TestScanPartialFixDoesNotDuplicateRemainingMatch(t *testing.T) {
	dir := t.TempDir()
	jsFile := filepath.Join(dir, "index.js")
	os.WriteFile(jsFile, []byte("function boot() {\n  console.log('boot');\n  console.log('ready');\n}\n"), 0644)
	stateDir := t.TempDir()

	first := &fakeDetector{
		matches: []detector.Match{
			{File: jsFile, Line: 2, Content: "  console.log('boot');"},
			{File: jsFile, Line: 3, Content: "  console.log('ready');"},
		},
	}
	if _, err := Scan(&ScanConfig{RecipeName: "test-recipe", RepoRoot: dir, StateDir: stateDir, Detector: first}); err != nil {
		t.Fatalf("first Scan failed: %v", err)
	}

	os.WriteFile(jsFile, []byte("function boot() {\n  logger.info('boot');\n  console.log('ready');\n}\n"), 0644)
	second := &fakeDetector{
		matches: []detector.Match{
			{File: jsFile, Line: 3, Content: "  console.log('ready');"},
		},
	}
	result, err := Scan(&ScanConfig{RecipeName: "test-recipe", RepoRoot: dir, StateDir: stateDir, Detector: second})
	if err != nil {
		t.Fatalf("second Scan failed: %v", err)
	}

	if result.Total != 2 {
		t.Fatalf("expected 2 total, got %d", result.Total)
	}
	if result.Done != 1 {
		t.Errorf("expected 1 done, got %d", result.Done)
	}
	if result.Remaining != 1 {
		t.Errorf("expected 1 remaining, got %d", result.Remaining)
	}

	s, err := state.Load(filepath.Join(stateDir, "test-recipe", ".scan-cache.json"))
	if err != nil {
		t.Fatalf("load state: %v", err)
	}
	var readyRemaining, readyDone bool
	for _, occ := range s.Occurrences {
		if occ.Content != "  console.log('ready');" {
			continue
		}
		if occ.Status == state.StatusTodo {
			readyRemaining = true
		}
		if occ.Status == state.StatusDone {
			readyDone = true
		}
	}
	if !readyRemaining {
		t.Error("expected ready match to remain todo")
	}
	if readyDone {
		t.Error("ready match must not appear as both done and todo")
	}
}
