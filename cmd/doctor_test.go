package cmd

import (
	"errors"
	"os/exec"
	"path/filepath"
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
