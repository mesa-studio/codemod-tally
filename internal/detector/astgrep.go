package detector

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// AstGrepDetector runs ast-grep with --json and parses results.
type AstGrepDetector struct {
	Rule     map[string]any
	Language string
}

func (d *AstGrepDetector) Run(dir string) ([]Match, error) {
	args, err := d.args(dir)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("ast-grep", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}
	return parseAstGrepJSON(out)
}

func (d *AstGrepDetector) args(dir string) ([]string, error) {
	pattern, ok := d.Rule["pattern"].(string)
	if !ok || pattern == "" {
		return nil, fmt.Errorf("ast-grep rule.pattern is required")
	}

	args := []string{"run", "--json=compact", "--pattern", pattern}
	if d.Language != "" {
		args = append(args, "--lang", d.Language)
	}
	return append(args, dir), nil
}

type astGrepResult struct {
	Text  string `json:"text"`
	Range struct {
		Start struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"start"`
	} `json:"range"`
	File string `json:"file"`
}

func parseAstGrepJSON(data []byte) ([]Match, error) {
	var results []astGrepResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	matches := make([]Match, 0, len(results))
	for _, r := range results {
		matches = append(matches, Match{
			File:    r.File,
			Line:    r.Range.Start.Line + 1, // ast-grep is 0-indexed
			Content: r.Text,
		})
	}
	return matches, nil
}
