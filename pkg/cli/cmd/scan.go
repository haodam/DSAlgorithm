package cmd

import (
	"fmt"

	"internal/scanner"

	"github.com/spf13/cobra"
)

var (
	scanPath  string
	highlight float64
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a directory and display file sizes",
	Long: `Scan recursively through a directory and display
the size of every file.

Example:

  getsize scan --path ./myproject
  getsize scan --path ./myproject --highlight 50`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runScan()
	},
}

func runScan() error {
	result, err := scanner.Scan(scanPath)
	if err != nil {
		return err
	}

	fmt.Printf("Scanning: %s\n\n", scanPath)

	fmt.Println("Files:")
	fmt.Println("---------------------------------------------")

	for _, file := range result.Files {
		if highlight > 0 &&
			float64(file.Size)/(1024*1024) >= highlight {

			fmt.Printf(
				"⚠️  %-60s %10s\n",
				file.Path,
				scanner.FormatSize(file.Size),
			)

			continue
		}

		fmt.Printf(
			"   %-60s %10s\n",
			file.Path,
			scanner.FormatSize(file.Size),
		)
	}

	fmt.Println("---------------------------------------------")
	fmt.Printf("Total files : %d\n", len(result.Files))
	fmt.Printf("Total size  : %s\n", scanner.FormatSize(result.TotalSize))

	return nil
}

func init() {
	scanCmd.Flags().StringVarP(
		&scanPath,
		"path",
		"p",
		".",
		"path of the directory to scan",
	)

	scanCmd.Flags().Float64Var(
		&highlight,
		"highlight",
		0,
		"highlight files larger than or equal to this size in MB",
	)

	rootCmd.AddCommand(scanCmd)
}
