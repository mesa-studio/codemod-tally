package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/mesa-studio/codemod-tally/internal/recipe"
	"github.com/mesa-studio/codemod-tally/internal/scanner"
	"github.com/mesa-studio/codemod-tally/internal/state"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan <name>",
	Short: "Run detector and update progress file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		recDir := filepath.Join(recipeDir, name)

		cfg, err := recipe.Load(recDir)
		if err != nil {
			return fmt.Errorf("load recipe %q: %w", name, err)
		}

		det, err := recipe.NewDetector(cfg)
		if err != nil {
			return fmt.Errorf("create detector: %w", err)
		}

		repoID, err := state.RepoID(workDir)
		if err != nil {
			return fmt.Errorf("identify repo: %w", err)
		}
		repoStateDir := filepath.Join(stateDir, repoID)

		fmt.Printf("Running detector (%s)...\n", cfg.DetectorConfig.Type)

		result, err := scanner.Scan(&scanner.ScanConfig{
			RecipeName:   name,
			RepoRoot:     workDir,
			StateDir:     repoStateDir,
			Detector:     det,
			IncludeGlobs: cfg.Scope.Include,
			ExcludeGlobs: cfg.Scope.Exclude,
		})
		if err != nil {
			return err
		}

		fmt.Println(scanSummaryLine(result.Total, result.New))
		fmt.Printf("  ✓ %d done\n", result.Done)
		fmt.Printf("  · %d remaining\n", result.Remaining)
		fmt.Printf("\nProgress: %s\n", filepath.Join(repoStateDir, name, "progress.md"))

		if result.Remaining == 0 && result.Total > 0 {
			fmt.Println("\n✓ All done!")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

func scanSummaryLine(total, new int) string {
	if new > 0 {
		return fmt.Sprintf("Tracking %d occurrences (%d new)", total, new)
	}
	return fmt.Sprintf("Tracking %d occurrences", total)
}
