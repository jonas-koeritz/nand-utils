package cmd

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"os"

	nand "github.com/jonas-koeritz/nand-utils/lib"
	"github.com/spf13/cobra"
)

var detectEccCmd = &cobra.Command{
	Use:   "detect-ecc",
	Short: "Tries to identify the ECC/CRC algorithm used and where the information is stored in the spare area",
	Run:   detectEcc,
}

func init() {
	detectEccCmd.Flags().IntVarP(&pageSize, "page-size", "p", 0, "The chips page size (in Bytes)")
	detectEccCmd.MarkFlagRequired("page-size")

	detectEccCmd.Flags().IntVarP(&spareSize, "spare-size", "s", 0, "The chips spare size per page (in Bytes)")
	detectEccCmd.MarkFlagRequired("spare-size")

	rootCmd.AddCommand(detectEccCmd)
}

func detectEcc(cmd *cobra.Command, args []string) {
	i, err := nand.OpenImageFile(imageFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	i.PageSize = pageSize
	i.SpareSize = spareSize

	pageData, spareData, err := i.ReadPage(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	offset, match := tryCRC32(pageData, spareData)
	if match {
		fmt.Printf("CRC32 of page data at offset %d in spare data\n", offset)
		return
	}

	offset, match = tryCRC64(pageData, spareData)
	if match {
		fmt.Printf("CRC64 of page data at offset %d in spare data\n", offset)
		return
	}

	offset, match = tryBCH(pageData, spareData)
	if match {
		fmt.Printf("BCH of page data at offset %d in spare data\n", offset)
		return
	}

	fmt.Printf("Could not detect the ECC/CRC algorithm\n")
}

func tryCRC32(pageData, spareData []byte) (int, bool) {
	checksum := crc32.ChecksumIEEE(pageData)
	needle := make([]byte, 4)
	binary.BigEndian.PutUint32(needle, checksum)

	fmt.Printf("Trying CRC32 big endian:\t\t%s\n", hex.EncodeToString(needle))

	offset := bytes.Index(needle, spareData)
	if offset >= 0 {
		return offset, true
	}

	binary.LittleEndian.PutUint32(needle, checksum)
	fmt.Printf("Trying CRC32 little endian:\t\t%s\n", hex.EncodeToString(needle))

	offset = bytes.Index(needle, spareData)
	if offset >= 0 {
		return offset, true
	}

	return 0, false
}

func tryCRC64(pageData, spareData []byte) (int, bool) {
	checksum := crc64.Checksum(pageData, crc64.MakeTable(crc64.ECMA))
	needle := make([]byte, 8)
	binary.BigEndian.PutUint64(needle, checksum)

	fmt.Printf("Trying ECMA CRC64 big endian:\t\t%s\n", hex.EncodeToString(needle))

	offset := bytes.Index(needle, spareData)
	if offset >= 0 {
		return offset, true
	}

	binary.LittleEndian.PutUint64(needle, checksum)
	fmt.Printf("Trying ECMA CRC64 little endian:\t%s\n", hex.EncodeToString(needle))

	offset = bytes.Index(needle, spareData)
	if offset >= 0 {
		return offset, true
	}

	checksum = crc64.Checksum(pageData, crc64.MakeTable(crc64.ISO))
	binary.BigEndian.PutUint64(needle, checksum)

	fmt.Printf("Trying ISO CRC64 big endian:\t\t%s\n", hex.EncodeToString(needle))

	offset = bytes.Index(needle, spareData)
	if offset >= 0 {
		return offset, true
	}

	binary.LittleEndian.PutUint64(needle, checksum)
	fmt.Printf("Trying ISO CRC64 little endian:\t\t%s\n", hex.EncodeToString(needle))

	offset = bytes.Index(needle, spareData)
	if offset >= 0 {
		return offset, true
	}

	return 0, false
}

func tryBCH(pageData, spareData []byte) (int, bool) {
	// TODO
	return 0, false
}
