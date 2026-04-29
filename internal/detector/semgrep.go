package detector

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

// SemgrepDetector runs semgrep with --json and parses results.
type SemgrepDetector struct {
	Rules []map[string]any
}

func (d *SemgrepDetector) Run(dir string) ([]Match, error) {
	ruleJSON, err := json.Marshal(map[string]any{"rules": d.Rules})
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("semgrep", "--json", "--config=-", dir)
	cmd.Stdin = bytes.NewReader(ruleJSON)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}
	return parseSemgrepJSON(out)
}

type semgrepOutput struct {
	Results []struct {
		Path  string `json:"path"`
		Start struct {
			Line int `json:"line"`
		} `json:"start"`
		Extra struct {
			Lines string `json:"lines"`
		} `json:"extra"`
	} `json:"results"`
}

func parseSemgrepJSON(data []byte) ([]Match, error) {
	var out semgrepOutput
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	matches := make([]Match, 0, len(out.Results))
	for _, r := range out.Results {
		matches = append(matches, Match{
			File:    r.Path,
			Line:    r.Start.Line,
			Content: r.Extra.Lines,
		})
	}
	return matches, nil
}
