package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mesa-studio/codemod-tally/internal/state"
	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt <name>",
	Short: "Print agent prompt block for clipboard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		recDir := filepath.Join(recipeDir, name)
		repoID, err := state.RepoID(workDir)
		if err != nil {
			return err
		}
		recipeStateDir := filepath.Join(stateDir, repoID, name)

		recipeMd := filepath.Join(recDir, "recipe.md")
		recipeContent, err := os.ReadFile(recipeMd)
		if err != nil {
			return fmt.Errorf("read recipe.md: %w", err)
		}

		cachePath := filepath.Join(recipeStateDir, ".scan-cache.json")
		s, err := state.Load(cachePath)
		if err != nil {
			return err
		}

		sep := strings.Repeat("─", 60)
		fmt.Println(sep)
		fmt.Printf("SWEEP TASK: %s\n", name)
		fmt.Println(sep)
		fmt.Printf("Recipe:   %s\n", recipeMd)
		fmt.Printf("Examples: %s\n", filepath.Join(recDir, "examples"))
		fmt.Printf("Progress: %s\n", filepath.Join(recipeStateDir, "progress.md"))
		fmt.Printf("Journal:  %s\n", filepath.Join(recipeStateDir, "journal.md"))
		fmt.Println()

		if len(s.Occurrences) > 0 {
			fmt.Printf("Progress: %d/%d done — %d remaining\n\n",
				s.Done(), len(s.Occurrences), s.Remaining())
		}

		fmt.Println("INSTRUCTIONS:")
		fmt.Println(strings.TrimSpace(string(recipeContent)))
		fmt.Println()
		fmt.Println("RULES:")
		fmt.Println("- Work through the Remaining items in TODO, top to bottom")
		fmt.Println("- Do NOT edit progress.md — it is updated by `codemod-tally scan`")
		fmt.Println("- If a case is not covered by the recipe, write a note in journal.md and skip it")
		fmt.Println("- Stop when context gets long — the user will run `codemod-tally scan` and resume")
		fmt.Println(sep)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(promptCmd)
}
