package recipe

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mesa-studio/codemod-tally/internal/detector"
	"gopkg.in/yaml.v3"
)

type ScopeConfig struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

type PostScanCheck struct {
	Command string `yaml:"command"`
}

type Config struct {
	Name          string         `yaml:"name"`
	Description   string         `yaml:"description"`
	Detector      string         `yaml:"detector"`
	Recipe        string         `yaml:"recipe"`
	ExamplesDir   string         `yaml:"examples_dir"`
	Scope         ScopeConfig    `yaml:"scope"`
	PostScanCheck *PostScanCheck `yaml:"post_scan_check"`

	Dir            string
	DetectorConfig DetectorConfig
}

type DetectorConfig struct {
	Type     string           `yaml:"type"`
	Command  string           `yaml:"command"`
	Parser   string           `yaml:"parser"`
	Pattern  string           `yaml:"pattern"`
	Flags    []string         `yaml:"flags"`
	Rules    []map[string]any `yaml:"rules"`
	Rule     map[string]any   `yaml:"rule"`
	Language string           `yaml:"language"`
}

// Load reads config.yaml and the referenced detector.yaml from dir.
func Load(dir string) (*Config, error) {
	configPath := filepath.Join(dir, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config.yaml: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config.yaml: %w", err)
	}
	cfg.Dir = dir

	detectorPath := filepath.Join(dir, cfg.Detector)
	detData, err := os.ReadFile(detectorPath)
	if err != nil {
		return nil, fmt.Errorf("read detector config %s: %w", cfg.Detector, err)
	}
	if err := yaml.Unmarshal(detData, &cfg.DetectorConfig); err != nil {
		return nil, fmt.Errorf("parse detector config: %w", err)
	}

	return &cfg, nil
}

// NewDetector constructs a Detector from the loaded config.
func NewDetector(cfg *Config) (detector.Detector, error) {
	dc := cfg.DetectorConfig
	switch dc.Type {
	case "ripgrep":
		return &detector.RipgrepDetector{Pattern: dc.Pattern, Flags: dc.Flags}, nil
	case "semgrep":
		return &detector.SemgrepDetector{Rules: dc.Rules}, nil
	case "astgrep":
		return &detector.AstGrepDetector{Rule: dc.Rule, Language: dc.Language}, nil
	case "shell":
		return &detector.ShellDetector{Command: dc.Command, Parser: dc.Parser}, nil
	default:
		return nil, fmt.Errorf("unknown detector type: %s", dc.Type)
	}
}
