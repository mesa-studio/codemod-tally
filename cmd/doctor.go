package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mesa-studio/codemod-tally/internal/recipe"
	"github.com/spf13/cobra"
)

type doctorResult struct {
	Name     string
	OK       bool
	Required bool
	Message  string
}

type lookPathFunc func(string) (string, error)

var doctorCmd = &cobra.Command{
	Use:   "doctor [recipe]",
	Short: "Check Codemod Tally dependencies and recipe readiness",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		results := checkEnvironmentDoctor(exec.LookPath)
		if len(args) == 1 {
			results = append(results, checkRecipeDoctor(filepath.Join(recipeDir, args[0]), exec.LookPath)...)
		}

		printDoctorResults(results)
		if hasRequiredDoctorFailure(results) {
			return fmt.Errorf("doctor found required failures")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func checkEnvironmentDoctor(lookPath lookPathFunc) []doctorResult {
	results := []doctorResult{
		checkExecutable("git", true, lookPath),
		checkExecutable("rg", true, lookPath),
		checkExecutable("semgrep", false, lookPath),
		checkExecutable("ast-grep", false, lookPath),
		checkDirectory("recipes directory", recipeDir, false),
		checkDirectory("state directory", stateDir, false),
	}
	return results
}

func checkRecipeDoctor(recipePath string, lookPath lookPathFunc) []doctorResult {
	cfg, err := recipe.Load(recipePath)
	if err != nil {
		return []doctorResult{{
			Name:     "recipe",
			OK:       false,
			Required: true,
			Message:  err.Error(),
		}}
	}

	results := []doctorResult{{
		Name:     "recipe",
		OK:       true,
		Required: true,
		Message:  recipePath,
	}}

	exe, ok := requiredDetectorExecutable(cfg.DetectorConfig.Type)
	if !ok {
		results = append(results, doctorResult{
			Name:     "detector",
			OK:       false,
			Required: true,
			Message:  fmt.Sprintf("unknown detector type %q", cfg.DetectorConfig.Type),
		})
		return results
	}

	results = append(results, checkExecutable(exe, true, lookPath))
	results = append(results, checkRecipeReadiness(cfg)...)
	return results
}

func checkRecipeReadiness(cfg *recipe.Config) []doctorResult {
	var results []doctorResult

	detectorPath := filepath.Join(cfg.Dir, cfg.Detector)
	detectorData, err := os.ReadFile(detectorPath)
	if err != nil {
		return []doctorResult{{
			Name:     "detector config",
			OK:       false,
			Required: true,
			Message:  err.Error(),
		}}
	}
	if containsAny(string(detectorData), []string{
		"your\\.pattern",
		"your.pattern",
		"foo(...)",
		"my-rule",
		"rg -n 'pattern'",
	}) {
		results = append(results, doctorResult{
			Name:     "detector readiness",
			OK:       false,
			Required: false,
			Message:  "detector.yaml still contains placeholder values; edit it before scanning a real repository",
		})
	}

	recipePath := filepath.Join(cfg.Dir, cfg.Recipe)
	recipeData, err := os.ReadFile(recipePath)
	if err != nil {
		return append(results, doctorResult{
			Name:     "recipe.md",
			OK:       false,
			Required: true,
			Message:  err.Error(),
		})
	}
	if containsAny(string(recipeData), []string{
		"Describe the transformation here.",
		"List exceptions",
		"See examples/ directory for reference diffs.",
	}) {
		results = append(results, doctorResult{
			Name:     "recipe readiness",
			OK:       false,
			Required: false,
			Message:  "recipe.md still contains scaffold text; replace it with agent instructions before scanning",
		})
	}

	return results
}

func containsAny(content string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(content, needle) {
			return true
		}
	}
	return false
}

func requiredDetectorExecutable(detectorType string) (string, bool) {
	switch detectorType {
	case "ripgrep":
		return "rg", true
	case "semgrep":
		return "semgrep", true
	case "astgrep":
		return "ast-grep", true
	case "shell":
		return "sh", true
	default:
		return "", false
	}
}

func checkExecutable(name string, required bool, lookPath lookPathFunc) doctorResult {
	path, err := lookPath(name)
	if err != nil {
		return doctorResult{
			Name:     name,
			OK:       false,
			Required: required,
			Message:  "not found in PATH",
		}
	}
	return doctorResult{
		Name:     name,
		OK:       true,
		Required: required,
		Message:  path,
	}
}

func checkDirectory(name string, path string, required bool) doctorResult {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return doctorResult{Name: name, OK: false, Required: required, Message: "does not exist yet: " + path}
		}
		return doctorResult{Name: name, OK: false, Required: required, Message: err.Error()}
	}
	if !info.IsDir() {
		return doctorResult{Name: name, OK: false, Required: required, Message: "not a directory: " + path}
	}
	return doctorResult{Name: name, OK: true, Required: required, Message: path}
}

func hasRequiredDoctorFailure(results []doctorResult) bool {
	for _, result := range results {
		if result.Required && !result.OK {
			return true
		}
	}
	return false
}

func printDoctorResults(results []doctorResult) {
	for _, result := range results {
		marker := "OK"
		if !result.OK && result.Required {
			marker = "FAIL"
		} else if !result.OK {
			marker = "WARN"
		}
		fmt.Printf("[%s] %-18s %s\n", marker, result.Name, result.Message)
	}
}
