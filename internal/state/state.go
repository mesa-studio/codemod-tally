package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Status string

const (
	StatusTodo Status = "todo"
	StatusDone Status = "done"
)

type Occurrence struct {
	Fingerprint Fingerprint `json:"fingerprint"`
	File        string      `json:"file"`
	Line        int         `json:"line"`
	Content     string      `json:"content"`
	Status      Status      `json:"status"`
	FirstSeen   time.Time   `json:"first_seen"`
	ResolvedAt  *time.Time  `json:"resolved_at"`
}

type State struct {
	Recipe      string       `json:"recipe"`
	RepoRoot    string       `json:"repo_root"`
	LastScan    time.Time    `json:"last_scan"`
	Occurrences []Occurrence `json:"occurrences"`
}

func (s *State) Done() int {
	n := 0
	for _, o := range s.Occurrences {
		if o.Status == StatusDone {
			n++
		}
	}
	return n
}

func (s *State) Remaining() int {
	return len(s.Occurrences) - s.Done()
}

// Load reads a .scan-cache.json file. Returns an empty State if the file doesn't exist.
func Load(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &State{}, nil
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Save writes State to path, creating parent directories as needed.
func Save(path string, s *State) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
