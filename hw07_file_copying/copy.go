package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceInfo, err := os.Stat(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source file %s does not exist: %w", fromPath, err)
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

	// simple progress indicator (no external packages)
	const chunkSize = 1 // 32 * 1024  -  чтобы проверять работу прогресс-бара проще было, в оригинале я 32кб использовал
	buf := make([]byte, chunkSize)
	var copied int64

	for copied < bytesToCopy {
		remaining := bytesToCopy - copied
		if remaining < int64(len(buf)) {
			buf = buf[:remaining]
		}

		n, readErr := sourceFile.Read(buf)
		if n > 0 {
			w, writeErr := destFile.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("error writing to destination file: %w", writeErr)
			}
			copied += int64(w)

			// update progress bar
			printProgressBar(copied, bytesToCopy, 40)
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return fmt.Errorf("error reading from source file: %w", readErr)
		}
		time.Sleep(time.Millisecond * 30) // тормозит выполнение, но зато прогресс бар красиво наполняется
	}

	printProgressBar(bytesToCopy, bytesToCopy, 40)
	fmt.Println()
	return nil
}

func printProgressBar(current, total int64, width int) {
	if total == 0 {
		return
	}
	ratio := float64(current) / float64(total)
	filled := int(ratio * float64(width))
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat(" ", width-filled)
	fmt.Printf("\r[%s] %3.0f%%", bar, ratio*100)
}
