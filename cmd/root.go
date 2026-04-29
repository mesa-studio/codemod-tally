package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	recipeDir string
	stateDir  string
	workDir   string
)

var rootCmd = &cobra.Command{
	Use:   "codemod-tally",
	Short: "Managed mass refactorings with LLM agents",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	home, _ := os.UserHomeDir()
	rootCmd.PersistentFlags().StringVar(&recipeDir, "recipe-dir", home+"/.codemod-tally/recipes", "recipes directory")
	rootCmd.PersistentFlags().StringVar(&stateDir, "state-dir", home+"/.codemod-tally/state", "state directory")
	rootCmd.PersistentFlags().StringVar(&workDir, "dir", ".", "target repository directory")
}
