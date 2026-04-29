package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestInitTemplateNames(t *testing.T) {
	got := templateNames()
	want := []string{"blank", "ripgrep-text", "semgrep-js", "astgrep-js"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestRenderRipgrepTextTemplate(t *testing.T) {
	files, err := renderTemplateFiles("console-to-logger", "ripgrep-text")
	if err != nil {
		t.Fatalf("renderTemplateFiles failed: %v", err)
	}

	if !strings.Contains(files["config.yaml"], "name: console-to-logger") {
		t.Fatalf("config.yaml missing recipe name:\n%s", files["config.yaml"])
	}
	if !strings.Contains(files["detector.yaml"], "type: ripgrep") {
		t.Fatalf("detector.yaml missing ripgrep type:\n%s", files["detector.yaml"])
	}
	if !strings.Contains(files["detector.yaml"], "pattern:") {
		t.Fatalf("detector.yaml missing pattern:\n%s", files["detector.yaml"])
	}
	if !strings.Contains(files["recipe.md"], "## What to change") {
		t.Fatalf("recipe.md missing instructions scaffold:\n%s", files["recipe.md"])
	}
}

func TestWriteRecipeFilesDoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(existing, []byte("keep me\n"), 0644); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	files, err := renderTemplateFiles("console-to-logger", "blank")
	if err != nil {
		t.Fatalf("renderTemplateFiles failed: %v", err)
	}
	actions, err := writeRecipeFiles(dir, files)
	if err != nil {
		t.Fatalf("writeRecipeFiles failed: %v", err)
	}

	data, err := os.ReadFile(existing)
	if err != nil {
		t.Fatalf("read existing config: %v", err)
	}
	if string(data) != "keep me\n" {
		t.Fatalf("expected existing file to be preserved, got %q", string(data))
	}

	var skippedConfig bool
	for _, action := range actions {
		if filepath.Base(action.Path) == "config.yaml" && !action.Created {
			skippedConfig = true
		}
	}
	if !skippedConfig {
		t.Fatalf("expected config.yaml skip action, got %+v", actions)
	}
}

func TestInitNextStepsGuideAgentWorkflow(t *testing.T) {
	got := initNextSteps("console-to-logger")
	for _, want := range []string{
		"codemod-tally doctor console-to-logger",
		"codemod-tally scan console-to-logger",
		"codemod-tally prompt console-to-logger",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected init next steps to contain %q:\n%s", want, got)
		}
	}
}
