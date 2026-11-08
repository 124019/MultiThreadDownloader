package lib

import (
	"os"
	"fmt"
	"MultiThreadDownloader/src/core"
)

type FixedSizeFile struct {
	F        *os.File
	FileSize int64
}

func CreatePlaceholderFile(filePath string, fileSize int) (*FixedSizeFile, error) {
	fileSize64 := int64(fileSize) // Ensure fileSize is an int64
	if fileSize <= 0 {
		return nil, fmt.Errorf("invalid file size: %d", fileSize64)
	}
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to create placeholder file: %w", err)
	}
	if err := f.Truncate(fileSize64); err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to truncate placeholder file: %w", err)
	}
	return &FixedSizeFile{F: f, FileSize: fileSize64}, nil
}

func (f *FixedSizeFile) WritePlaceholder(c core.CHUNK_Range, data []byte) error {
	length := c.Ed - c.St + 1
	if len(data) != length {
		return fmt.Errorf("data length does not match chunk range: %d(data) != %d(want)", len(data), length)
	}
	_, err := f.F.WriteAt(data, int64(c.St))
	return err
}