package cmd

import (
	"fmt"
	"os"

	nand "github.com/jonas-koeritz/nand-utils/lib"
	"github.com/spf13/cobra"
)

var dumpPageIndex int

var dumpPageCmd = &cobra.Command{
	Use:   "dump-page",
	Short: "Dump a single page",
	Run:   dumpPage,
}

func init() {
	dumpPageCmd.Flags().IntVarP(&pageSize, "page-size", "p", 0, "The chips page size (in Bytes)")
	dumpPageCmd.MarkFlagRequired("page-size")

	dumpPageCmd.Flags().IntVarP(&spareSize, "spare-size", "s", 0, "The chips spare size per page (in Bytes)")
	dumpPageCmd.MarkFlagRequired("spare-size")

	dumpPageCmd.Flags().IntVar(&dumpPageIndex, "page-index", 0, "The index of the page to dump")
	dumpPageCmd.MarkFlagRequired("page-index")

	rootCmd.AddCommand(dumpPageCmd)
}

func dumpPage(cmd *cobra.Command, args []string) {
	i, err := nand.OpenImageFile(imageFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	i.PageSize = pageSize
	i.SpareSize = spareSize

	dump, err := i.DumpPage(dumpPageIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println(dump)
}
