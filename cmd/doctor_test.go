package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRequiredDetectorExecutable(t *testing.T) {
	cases := map[string]string{
		"ripgrep": "rg",
		"semgrep": "semgrep",
		"astgrep": "ast-grep",
		"shell":   "sh",
	}

	for detectorType, want := range cases {
		got, ok := requiredDetectorExecutable(detectorType)
		if !ok {
			t.Fatalf("expected executable for %s", detectorType)
		}
		if got != want {
			t.Fatalf("expected %s for %s, got %s", want, detectorType, got)
		}
	}
}

func TestDoctorCheckExecutableReportsMissingRequired(t *testing.T) {
	result := checkExecutable("rg", true, func(string) (string, error) {
		return "", exec.ErrNotFound
	})

	if result.OK {
		t.Fatal("expected missing required executable to fail")
	}
	if !result.Required {
		t.Fatal("expected result to be required")
	}
}

func TestDoctorRecipeReportsMissingDetectorExecutable(t *testing.T) {
	dir := t.TempDir()
	files, err := renderTemplateFiles("js-migration", "semgrep-js")
	if err != nil {
		t.Fatalf("render template: %v", err)
	}
	if _, err := writeRecipeFiles(filepath.Join(dir, "js-migration"), files); err != nil {
		t.Fatalf("write recipe: %v", err)
	}

	results := checkRecipeDoctor(filepath.Join(dir, "js-migration"), func(name string) (string, error) {
		if name == "semgrep" {
			return "", exec.ErrNotFound
		}
		return "/usr/bin/" + name, nil
	})

	if !hasRequiredDoctorFailure(results) {
		t.Fatalf("expected required doctor failure, got %+v", results)
	}
}

func TestDoctorRecipeReportsUnknownDetectorType(t *testing.T) {
	dir := t.TempDir()
	files, err := renderTemplateFiles("bad-recipe", "blank")
	if err != nil {
		t.Fatalf("render template: %v", err)
	}
	files["detector.yaml"] = "type: nope\n"
	if _, err := writeRecipeFiles(filepath.Join(dir, "bad-recipe"), files); err != nil {
		t.Fatalf("write recipe: %v", err)
	}

	results := checkRecipeDoctor(filepath.Join(dir, "bad-recipe"), func(string) (string, error) {
		return "", errors.New("should not be called")
	})

	if !hasRequiredDoctorFailure(results) {
		t.Fatalf("expected required doctor failure, got %+v", results)
	}
}

func TestDoctorRecipeWarnsOnScaffoldPlaceholders(t *testing.T) {
	dir := t.TempDir()
	files, err := renderTemplateFiles("console-to-logger", "ripgrep-text")
	if err != nil {
		t.Fatalf("render template: %v", err)
	}
	if _, err := writeRecipeFiles(filepath.Join(dir, "console-to-logger"), files); err != nil {
		t.Fatalf("write recipe: %v", err)
	}

	results := checkRecipeDoctor(filepath.Join(dir, "console-to-logger"), func(name string) (string, error) {
		return "/usr/bin/" + name, nil
	})

	if hasRequiredDoctorFailure(results) {
		t.Fatalf("placeholder warnings must not fail doctor, got %+v", results)
	}
	if !hasDoctorWarning(results, "detector readiness", "placeholder") {
		t.Fatalf("expected detector placeholder warning, got %+v", results)
	}
	if !hasDoctorWarning(results, "recipe readiness", "scaffold") {
		t.Fatalf("expected recipe scaffold warning, got %+v", results)
	}
}

func TestDoctorRecipeDoesNotWarnAfterRecipeIsEdited(t *testing.T) {
	dir := t.TempDir()
	recipeDir := filepath.Join(dir, "console-to-logger")
	files, err := renderTemplateFiles("console-to-logger", "ripgrep-text")
	if err != nil {
		t.Fatalf("render template: %v", err)
	}
	if _, err := writeRecipeFiles(recipeDir, files); err != nil {
		t.Fatalf("write recipe: %v", err)
	}
	if err := os.WriteFile(filepath.Join(recipeDir, "detector.yaml"), []byte("type: ripgrep\npattern: 'console\\.log\\('\nflags: []\n"), 0644); err != nil {
		t.Fatalf("write detector: %v", err)
	}
	if err := os.WriteFile(filepath.Join(recipeDir, "recipe.md"), []byte("# console-to-logger\n\nReplace console.log calls with logger.info while leaving tests unchanged.\n"), 0644); err != nil {
		t.Fatalf("write recipe.md: %v", err)
	}

	results := checkRecipeDoctor(recipeDir, func(name string) (string, error) {
		return "/usr/bin/" + name, nil
	})

	if hasDoctorWarning(results, "detector readiness", "") || hasDoctorWarning(results, "recipe readiness", "") {
		t.Fatalf("expected no readiness warnings after edits, got %+v", results)
	}
}

func hasDoctorWarning(results []doctorResult, name string, messagePart string) bool {
	for _, result := range results {
		if result.Name != name || result.OK || result.Required {
			continue
		}
		if messagePart == "" || strings.Contains(result.Message, messagePart) {
			return true
		}
	}
	return false
}
