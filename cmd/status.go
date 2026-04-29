package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/mesa-studio/codemod-tally/internal/state"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <name>",
	Short: "Show progress from cached state (no detector run)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		repoID, err := state.RepoID(workDir)
		if err != nil {
			return err
		}

		cachePath := filepath.Join(stateDir, repoID, name, ".scan-cache.json")
		s, err := state.Load(cachePath)
		if err != nil {
			return err
		}
		if len(s.Occurrences) == 0 {
			fmt.Printf("No state for %q. Run: codemod-tally scan %s\n", name, name)
			return nil
		}

		total := len(s.Occurrences)
		done := s.Done()
		pct := 0
		if total > 0 {
			pct = done * 100 / total
		}

		bar := progressBar(pct, 25)
		fmt.Printf("%s (last scan: %s)\n", name, s.LastScan.Format("2006-01-02 15:04"))
		fmt.Printf("  %d/%d done (%d%%)  %s  %d remaining\n",
			done, total, pct, bar, s.Remaining())
		return nil
	},
}

func progressBar(pct, width int) string {
	filled := pct * width / 100
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
