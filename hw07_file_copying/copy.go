package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	sourceInfo, err := os.Stat(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source file does not exist: %s", fromPath)
		}
		return fmt.Errorf("error accessing source file: %w", err)
	}

	if !sourceInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := sourceInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	remainingBytes := fileSize - offset

	bytesToCopy := remainingBytes
	if limit > 0 && limit < remainingBytes {
		bytesToCopy = limit
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}

	defer sourceFile.Close()

	destFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)

	}
	defer destFile.Close()

	_, err = sourceFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error seeking source file: %w", err)
	}
	// прогресс бар
	bar := pb.New64(bytesToCopy)
	bar.SetRefreshRate(time.Millisecond * 100)
	bar.Start()
	reader := bar.NewProxyReader(sourceFile)
	_, err = io.CopyN(destFile, reader, bytesToCopy)
	if err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}
	bar.Finish()
	return nil
}
