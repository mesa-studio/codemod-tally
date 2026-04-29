package recipe

import (
	"testing"
)

func TestLoadRecipe(t *testing.T) {
	cfg, err := Load("testdata")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Name != "console-to-logger" {
		t.Errorf("unexpected name: %s", cfg.Name)
	}
	if len(cfg.Scope.Exclude) != 2 {
		t.Errorf("expected 2 exclude patterns, got %d", len(cfg.Scope.Exclude))
	}
}

func TestNewDetectorRipgrep(t *testing.T) {
	cfg, _ := Load("testdata")
	d, err := NewDetector(cfg)
	if err != nil {
		t.Fatalf("NewDetector failed: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestLoadInvalidDir(t *testing.T) {
	_, err := Load("/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}
