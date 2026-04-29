package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	initTemplateName  string
	initListTemplates bool
)

var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Create a new recipe skeleton",
	Args: func(cmd *cobra.Command, args []string) error {
		if initListTemplates {
			if len(args) != 0 {
				return errors.New("--list-templates does not accept a recipe name")
			}
			return nil
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if initListTemplates {
			for _, tmpl := range recipeTemplates {
				fmt.Printf("%-14s %s\n", tmpl.Name, tmpl.Description)
			}
			return nil
		}

		name := args[0]
		dir := filepath.Join(recipeDir, name)

		files, err := renderTemplateFiles(name, initTemplateName)
		if err != nil {
			return err
		}

		actions, err := writeRecipeFiles(dir, files)
		if err != nil {
			return err
		}
		for _, action := range actions {
			if action.Created {
				fmt.Printf("  created: %s\n", action.Path)
			} else {
				fmt.Printf("  skip (exists): %s\n", action.Path)
			}
		}
		fmt.Printf("\nRecipe %q created at %s\n", name, dir)
		fmt.Print(initNextSteps(name))
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initTemplateName, "template", "blank", "recipe template")
	initCmd.Flags().BoolVar(&initListTemplates, "list-templates", false, "list available recipe templates")
	rootCmd.AddCommand(initCmd)
}

type initRecipeTemplate struct {
	Name        string
	Description string
	Config      func(string) string
	Detector    string
	Recipe      func(string) string
}

type fileAction struct {
	Path    string
	Created bool
}

var recipeTemplates = []initRecipeTemplate{
	{
		Name:        "blank",
		Description: "empty recipe skeleton",
		Config:      configTemplate,
		Detector:    detectorTemplate,
		Recipe:      recipeTemplate,
	},
	{
		Name:        "ripgrep-text",
		Description: "text search with ripgrep",
		Config:      ripgrepTextConfigTemplate,
		Detector:    ripgrepTextDetectorTemplate,
		Recipe:      recipeTemplate,
	},
	{
		Name:        "semgrep-js",
		Description: "JavaScript/TypeScript AST search with semgrep",
		Config:      semgrepJSConfigTemplate,
		Detector:    semgrepJSDetectorTemplate,
		Recipe:      recipeTemplate,
	},
	{
		Name:        "astgrep-js",
		Description: "JavaScript AST search with ast-grep",
		Config:      astGrepJSConfigTemplate,
		Detector:    astGrepJSDetectorTemplate,
		Recipe:      recipeTemplate,
	},
}

func templateNames() []string {
	names := make([]string, 0, len(recipeTemplates))
	for _, tmpl := range recipeTemplates {
		names = append(names, tmpl.Name)
	}
	return names
}

func renderTemplateFiles(name string, templateName string) (map[string]string, error) {
	for _, tmpl := range recipeTemplates {
		if tmpl.Name == templateName {
			return map[string]string{
				"config.yaml":   tmpl.Config(name),
				"detector.yaml": tmpl.Detector,
				"recipe.md":     tmpl.Recipe(name),
			}, nil
		}
	}
	return nil, fmt.Errorf("unknown template %q (available: %s)", templateName, strings.Join(templateNames(), ", "))
}

func writeRecipeFiles(dir string, files map[string]string) ([]fileAction, error) {
	if err := os.MkdirAll(filepath.Join(dir, "examples"), 0755); err != nil {
		return nil, err
	}

	order := []string{"config.yaml", "detector.yaml", "recipe.md"}
	actions := make([]fileAction, 0, len(order))
	for _, fname := range order {
		content, ok := files[fname]
		if !ok {
			continue
		}
		path := filepath.Join(dir, fname)
		action := fileAction{Path: path}
		if _, err := os.Stat(path); err == nil {
			actions = append(actions, action)
			continue
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return nil, err
		}
		action.Created = true
		actions = append(actions, action)
	}
	return actions, nil
}

func initNextSteps(name string) string {
	return fmt.Sprintf(`
Next steps:
  1. Edit detector.yaml and recipe.md for the migration.
  2. Run: codemod-tally doctor %s
  3. Run: codemod-tally scan %s
  4. Run: codemod-tally prompt %s
`, name, name, name)
}

func configTemplate(name string) string {
	return fmt.Sprintf(`name: %s
description: ""
detector: detector.yaml
recipe: recipe.md
examples_dir: examples/
scope:
  include: []
  exclude: []
`, name)
}

func ripgrepTextConfigTemplate(name string) string {
	return fmt.Sprintf(`name: %s
description: Text pattern migration
detector: detector.yaml
recipe: recipe.md
examples_dir: examples/
scope:
  include: []
  exclude: ["**/node_modules/**", "**/.git/**"]
`, name)
}

func semgrepJSConfigTemplate(name string) string {
	return fmt.Sprintf(`name: %s
description: JavaScript/TypeScript AST migration
detector: detector.yaml
recipe: recipe.md
examples_dir: examples/
scope:
  include: ["**/*.js", "**/*.jsx", "**/*.ts", "**/*.tsx"]
  exclude: ["**/node_modules/**", "**/dist/**", "**/build/**"]
`, name)
}

func astGrepJSConfigTemplate(name string) string {
	return fmt.Sprintf(`name: %s
description: JavaScript AST migration
detector: detector.yaml
recipe: recipe.md
examples_dir: examples/
scope:
  include: ["**/*.js", "**/*.jsx", "**/*.ts", "**/*.tsx"]
  exclude: ["**/node_modules/**", "**/dist/**", "**/build/**"]
`, name)
}

const detectorTemplate = `# Choose one detector type and remove the others.

# ripgrep (text pattern search)
type: ripgrep
pattern: 'your\.pattern\('
flags: []

# shell (arbitrary command — output parsed as file:line:content)
# type: shell
# command: "rg -n 'pattern' --json"
# parser: ripgrep   # ripgrep | semgrep | astgrep | lines

# semgrep (AST-level, requires semgrep installed)
# type: semgrep
# rules:
#   - id: my-rule
#     pattern: foo(...)
#     languages: [javascript]

# ast-grep (AST-level, requires ast-grep installed)
# type: astgrep
# language: JavaScript
# rule:
#   pattern: foo($$$ARGS)
`

const ripgrepTextDetectorTemplate = `type: ripgrep
pattern: 'your\.pattern'
flags: []
`

const semgrepJSDetectorTemplate = `type: semgrep
rules:
  - id: codemod-tally-js-pattern
    pattern: console.log(...)
    languages: [javascript, typescript]
`

const astGrepJSDetectorTemplate = `type: astgrep
language: JavaScript
rule:
  pattern: console.log($$$ARGS)
`

func recipeTemplate(name string) string {
	return fmt.Sprintf(`# %s

## What to change

Describe the transformation here.

## Do NOT touch

- List exceptions

## Examples

See examples/ directory for reference diffs.
`, name)
}
