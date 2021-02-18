package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var imageFilePath string
var pageSize int
var spareSize int
var blockSize int

var rootCmd = &cobra.Command{
	Use:   "nand",
	Short: "nand is a utility to aid with the analysis and manipulation of NAND flash image files",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&imageFilePath, "image", "i", "", "Path to the flash image file (required)")
	rootCmd.MarkPersistentFlagRequired("image")
}

// Execute the main program
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
