package detector

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

// RipgrepDetector runs rg with --json output and parses results.
type RipgrepDetector struct {
	Pattern string
	Flags   []string
}

func (d *RipgrepDetector) Run(dir string) ([]Match, error) {
	args := append([]string{"--json", d.Pattern}, d.Flags...)
	args = append(args, dir)
	cmd := exec.Command("rg", args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}
	return parseRipgrepJSON(out)
}

type rgLine struct {
	Type string `json:"type"`
	Data struct {
		Path struct {
			Text string `json:"text"`
		} `json:"path"`
		Lines struct {
			Text string `json:"text"`
		} `json:"lines"`
		LineNumber int `json:"line_number"`
	} `json:"data"`
}

func parseRipgrepJSON(data []byte) ([]Match, error) {
	var matches []Match
	for _, line := range bytes.Split(data, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var rg rgLine
		if err := json.Unmarshal(line, &rg); err != nil {
			continue
		}
		if rg.Type != "match" {
			continue
		}
		matches = append(matches, Match{
			File:    rg.Data.Path.Text,
			Line:    rg.Data.LineNumber,
			Content: strings.TrimRight(rg.Data.Lines.Text, "\n\r"),
		})
	}
	return matches, nil
}
