package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/mesa-studio/codemod-tally/internal/detector"
	"github.com/mesa-studio/codemod-tally/internal/progress"
	"github.com/mesa-studio/codemod-tally/internal/state"
)

const contextRadius = 3

type ScanConfig struct {
	RecipeName   string
	RepoRoot     string
	StateDir     string
	Detector     detector.Detector
	IncludeGlobs []string
	ExcludeGlobs []string
}

type ScanResult struct {
	Total     int
	Done      int
	Remaining int
	New       int
}

// Scan runs the detector, merges with existing state, and writes progress.md.
func Scan(cfg *ScanConfig) (*ScanResult, error) {
	matches, err := cfg.Detector.Run(cfg.RepoRoot)
	if err != nil {
		return nil, fmt.Errorf("detector: %w", err)
	}

	if len(cfg.IncludeGlobs) > 0 {
		matches, err = filterIncluded(matches, cfg.RepoRoot, cfg.IncludeGlobs)
		if err != nil {
			return nil, fmt.Errorf("scope filter: %w", err)
		}
	}

	if len(cfg.ExcludeGlobs) > 0 {
		matches, err = filterExcluded(matches, cfg.RepoRoot, cfg.ExcludeGlobs)
		if err != nil {
			return nil, fmt.Errorf("scope filter: %w", err)
		}
	}

	enriched, err := enrichContext(matches, cfg.RepoRoot, contextRadius)
	if err != nil {
		return nil, fmt.Errorf("context enrichment: %w", err)
	}

	recipeStateDir := filepath.Join(cfg.StateDir, cfg.RecipeName)
	cachePath := filepath.Join(recipeStateDir, ".scan-cache.json")

	existing, err := state.Load(cachePath)
	if err != nil {
		return nil, fmt.Errorf("load state: %w", err)
	}
	if existing.Recipe == "" {
		existing.Recipe = cfg.RecipeName
		existing.RepoRoot = cfg.RepoRoot
	}

	prevTotal := len(existing.Occurrences)
	merged := state.Merge(existing, enriched, time.Now())

	if err := state.Save(cachePath, merged); err != nil {
		return nil, fmt.Errorf("save state: %w", err)
	}

	progressContent := progress.Generate(merged)
	progressPath := filepath.Join(recipeStateDir, "progress.md")
	if err := os.WriteFile(progressPath, []byte(progressContent), 0644); err != nil {
		return nil, fmt.Errorf("write progress.md: %w", err)
	}

	return &ScanResult{
		Total:     len(merged.Occurrences),
		Done:      merged.Done(),
		Remaining: merged.Remaining(),
		New:       len(merged.Occurrences) - prevTotal,
	}, nil
}

// filterIncluded removes matches whose relative path does not match any include glob.
func filterIncluded(matches []detector.Match, root string, globs []string) ([]detector.Match, error) {
	var result []detector.Match
	for _, m := range matches {
		rel, err := filepath.Rel(root, m.File)
		if err != nil {
			rel = m.File
		}
		if matchesAnyGlob(globs, rel) {
			result = append(result, m)
		}
	}
	return result, nil
}

// filterExcluded removes matches whose relative path matches any exclude glob.
func filterExcluded(matches []detector.Match, root string, globs []string) ([]detector.Match, error) {
	var result []detector.Match
	for _, m := range matches {
		rel, err := filepath.Rel(root, m.File)
		if err != nil {
			rel = m.File
		}
		if !matchesAnyGlob(globs, rel) {
			result = append(result, m)
		}
	}
	return result, nil
}

func matchesAnyGlob(globs []string, rel string) bool {
	for _, pattern := range globs {
		if matched, _ := doublestar.Match(pattern, rel); matched {
			return true
		}
	}
	return false
}

// enrichContext reads each matched file to populate Context field (±radius lines).
// Relative file paths are resolved against root.
func enrichContext(matches []detector.Match, root string, radius int) ([]detector.Match, error) {
	fileLines := make(map[string][]string)
	result := make([]detector.Match, 0, len(matches))

	for _, m := range matches {
		absFile := m.File
		if !filepath.IsAbs(absFile) {
			absFile = filepath.Join(root, absFile)
		}
		if _, ok := fileLines[absFile]; !ok {
			lines, err := readLines(absFile)
			if err != nil {
				return nil, fmt.Errorf("read %s: %w", absFile, err)
			}
			fileLines[absFile] = lines
		}
		m.File = absFile
		m.Context, m.ContextLine = extractContext(fileLines[absFile], m.Line, radius)
		result = append(result, m)
	}
	return result, nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

// extractContext returns lines [line-radius .. line+radius] and the matched
// line's zero-based index inside that context (line is 1-indexed).
func extractContext(lines []string, line, radius int) ([]string, int) {
	start := line - 1 - radius
	if start < 0 {
		start = 0
	}
	end := line - 1 + radius + 1
	if end > len(lines) {
		end = len(lines)
	}
	return lines[start:end], line - 1 - start
}
