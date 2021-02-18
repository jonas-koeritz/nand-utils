package nand

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// ImageFile provides functions to read from a NAND flash image file
type ImageFile struct {
	file      *os.File
	PageSize  int
	SpareSize int
}

// OpenImageFile opens a NAND flash file for reading
func OpenImageFile(path string) (*ImageFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	i := &ImageFile{
		file:      file,
		PageSize:  0,
		SpareSize: 0,
	}
	return i, nil
}

// ReadPage reads the next page from the image file
func (i *ImageFile) ReadPage(pageIndex int) (pageData []byte, spareData []byte, err error) {
	if i.PageSize == 0 || i.SpareSize == 0 {
		return nil, nil, errors.New("Unknown page size or spare size")
	}

	totalPageSize := i.PageSize + i.SpareSize
	offset := int64(pageIndex * totalPageSize)

	fullPageData := make([]byte, totalPageSize)

	_, err = i.file.ReadAt(fullPageData, offset)
	if err != nil {
		return nil, nil, err
	}

	pageData = fullPageData[:i.PageSize]
	spareData = fullPageData[i.PageSize:]

	return
}

// DumpPage returns a string that includes hexdumps for the page and spare contents
func (i *ImageFile) DumpPage(pageIndex int) (string, error) {
	pageData, spareData, err := i.ReadPage(pageIndex)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("Page Data:\n%s\nSpare Data:\n%s", hex.Dump(pageData), hex.Dump(spareData))
	return result, nil
}
