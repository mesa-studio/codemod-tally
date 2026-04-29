package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mesa-studio/codemod-tally/internal/recipe"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available recipes",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := os.ReadDir(recipeDir)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No recipes found. Run: codemod-tally init <name>")
				return nil
			}
			return err
		}

		if len(entries) == 0 {
			fmt.Println("No recipes found. Run: codemod-tally init <name>")
			return nil
		}

		fmt.Println("Available recipes:")
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			dir := filepath.Join(recipeDir, e.Name())
			cfg, err := recipe.Load(dir)
			desc := ""
			if err == nil {
				desc = cfg.Description
			}
			if desc != "" {
				fmt.Printf("  %-24s %s\n", e.Name(), desc)
			} else {
				fmt.Printf("  %s\n", e.Name())
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
