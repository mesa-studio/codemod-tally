package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mesa-studio/codemod-tally/internal/state"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean <name>",
	Short: "Delete scan state for current repo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		repoID, err := state.RepoID(workDir)
		if err != nil {
			return err
		}
		target := filepath.Join(stateDir, repoID, name)

		if _, err := os.Stat(target); os.IsNotExist(err) {
			fmt.Printf("No state found for %q in current repo.\n", name)
			return nil
		}

		fmt.Printf("This will delete: %s\n", target)
		fmt.Print("Confirm? [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		resp, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(resp)) != "y" {
			fmt.Println("Cancelled.")
			return nil
		}

		if err := os.RemoveAll(target); err != nil {
			return err
		}
		fmt.Println("State cleaned.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
