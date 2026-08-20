package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "getsize",
	Short: "A simple directory size checker CLI tool",
	Long: `getsize is a simple CLI tool for scanning
and displaying the size of files inside a directory.

Examples:

  getsize scan --path ./myproject
  getsize scan --path ./myproject --highlight 50
  getsize scan --path ./myproject --verbose`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("Verbose mode enabled")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"enable verbose output",
	)
}
