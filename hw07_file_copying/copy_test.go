package main

import (
	"errors"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	fileName := "testdata/input.txt"
	to := "testdata/input_copy.txt"
	err := Copy(fileName, to, 0, 10)
	if err != nil {
		t.Errorf("Copy failed: %v", err)
	}
	os.Remove(to)
}

const inputCopy = "testdata/input_copy.txt"

func TestOffsetIsGreaterThanFileSize(t *testing.T) {
	fileName := "testdata/input.txt"
	to := inputCopy
	err := Copy(fileName, to, 1000000000000000, 10)
	if !errors.Is(err, ErrOffsetExceedsFileSize) {
		t.Fatalf("expected ErrOffsetExceedsFileSize, got %v", err)
	}
	os.Remove(to)
}

func TestUnsupportedFile(t *testing.T) {
	fileName := "testdata/"
	to := inputCopy
	err := Copy(fileName, to, 0, 10)
	if !errors.Is(err, ErrUnsupportedFile) {
		t.Fatalf("expected ErrUnsupportedFile, got %v", err)
	}
	os.Remove(to)
}

func TestSourceFileDoesNotExist(t *testing.T) {
	fileName := "testdata/not_exists.txt"
	to := inputCopy
	err := Copy(fileName, to, 0, 10)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got %v", err)
	}
	os.Remove(to)
}
