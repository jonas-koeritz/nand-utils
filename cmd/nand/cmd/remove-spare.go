package cmd

import (
	"fmt"
	"os"

	nand "github.com/jonas-koeritz/nand-utils/lib"
	"github.com/spf13/cobra"
)

var removeSpareOutputFile string

var removeSpareCmd = &cobra.Command{
	Use:     "remove-spare",
	Aliases: []string{"remove-oob"},
	Short:   "Remove spare data from pages to prepare image for further analysis",
	Run:     removeSpare,
}

func init() {
	removeSpareCmd.Flags().IntVarP(&pageSize, "page-size", "p", 0, "The chips page size (in Bytes)")
	removeSpareCmd.MarkFlagRequired("page-size")

	removeSpareCmd.Flags().IntVarP(&spareSize, "spare-size", "s", 0, "The chips spare size per page (in Bytes)")
	removeSpareCmd.MarkFlagRequired("spare-size")

	removeSpareCmd.Flags().StringVarP(&removeSpareOutputFile, "output-file", "o", "", "The output file to write to")
	removeSpareCmd.MarkFlagRequired("output-file")

	rootCmd.AddCommand(removeSpareCmd)

}

func removeSpare(cmd *cobra.Command, args []string) {
	i, err := nand.OpenImageFile(imageFilePath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	outputFile, err := os.OpenFile(removeSpareOutputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	i.PageSize = pageSize
	i.SpareSize = spareSize

	currentPage := 0
	for ; ; currentPage++ {
		pageData, _, err := i.ReadPage(currentPage)
		if err != nil {
			break
		}
		outputFile.Write(pageData)
	}
	outputFile.Close()
}
