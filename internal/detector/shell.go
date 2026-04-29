package detector

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ShellDetector runs an arbitrary shell command and parses its output.
// Parser must be one of: "ripgrep", "semgrep", "astgrep", "lines".
type ShellDetector struct {
	Command string
	Parser  string
}

func (d *ShellDetector) Run(dir string) ([]Match, error) {
	cmd := exec.Command("sh", "-c", d.Command)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("detector command failed: %w", err)
	}
	return d.parse(out)
}

func (d *ShellDetector) parse(data []byte) ([]Match, error) {
	switch d.Parser {
	case "ripgrep":
		return parseRipgrepJSON(data)
	case "semgrep":
		return parseSemgrepJSON(data)
	case "astgrep":
		return parseAstGrepJSON(data)
	case "lines", "":
		return parseLines(data)
	default:
		return nil, fmt.Errorf("unknown parser: %s", d.Parser)
	}
}

// parseLines parses output in the format "file:line:content".
func parseLines(data []byte) ([]Match, error) {
	var matches []Match
	for _, line := range bytes.Split(data, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(string(line), ":", 3)
		if len(parts) < 2 {
			continue
		}
		lineNum, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		content := ""
		if len(parts) == 3 {
			content = parts[2]
		}
		matches = append(matches, Match{
			File:    parts[0],
			Line:    lineNum,
			Content: content,
		})
	}
	return matches, nil
}
