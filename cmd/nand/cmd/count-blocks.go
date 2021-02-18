package cmd

import (
	"fmt"
	"os"

	nand "github.com/jonas-koeritz/nand-utils/lib"
	"github.com/spf13/cobra"
)

var countBlocksCmd = &cobra.Command{
	Use:   "count-blocks",
	Short: "Count the number of blocks in the given image file",
	Run:   countBlocks,
}

func init() {
	countBlocksCmd.Flags().IntVarP(&pageSize, "page-size", "p", 0, "The chips page size (in Bytes)")
	countBlocksCmd.MarkFlagRequired("page-size")

	countBlocksCmd.Flags().IntVarP(&spareSize, "spare-size", "s", 0, "The chips spare size per page (in Bytes)")
	countBlocksCmd.MarkFlagRequired("spare-size")

	countBlocksCmd.Flags().IntVarP(&blockSize, "block-size", "b", 0, "Number of pages that make up a block")
	countBlocksCmd.MarkFlagRequired("block-size")

	rootCmd.AddCommand(countBlocksCmd)
}

func countBlocks(cmd *cobra.Command, args []string) {
	i, err := nand.OpenImageFile(imageFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	i.PageSize = pageSize
	i.SpareSize = spareSize

	currentPage := 0
	for ; ; currentPage++ {
		_, _, err := i.ReadPage(currentPage)
		if err != nil {
			break
		}
	}

	blocks := currentPage / blockSize
	fmt.Printf("%d blocks (%d pages)\n", blocks, currentPage)
}
