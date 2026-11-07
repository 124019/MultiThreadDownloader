package lib

import (
	"os";
	"fmt"
)

type FixedSizeFile struct {
	F        *os.File
	FileSize int64
}

func CreatePlaceholderFile(filePath string, fileSize int64) (*FixedSizeFile, error) {
	if fileSize <= 0 {
		return nil, fmt.Errorf("invalid file size: %d", fileSize)
	}
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to create placeholder file: %w", err)
	}
	if err := f.Truncate(fileSize); err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to truncate placeholder file: %w", err)
	}
	return &FixedSizeFile{F: f, FileSize: fileSize}, nil
}