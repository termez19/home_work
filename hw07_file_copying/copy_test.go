package main

import "testing"

func TestCopy(t *testing.T) {
	fileName := "testdata/input.txt"
	to := "testdata/input_copy.txt"
	err := Copy(fileName, to, 0, 10)
	if err != nil {
		t.Errorf("Copy failed: %v", err)
	}
}
func TestOffsetIsGreaterThanFileSize(t *testing.T) {
	fileName := "testdata/input.txt"
	to := "testdata/input_copy.txt"
	err := Copy(fileName, to, 1000000000000000, 10)
	if err != ErrOffsetExceedsFileSize {
		t.Errorf("expected ErrOffsetExceedsFileSize, got %v", err)
	}
}

func TestUnsupportedFile(t *testing.T) {
	fileName := "testdata/"
	to := "testdata/input_copy.txt"
	err := Copy(fileName, to, 0, 10)
	if err != ErrUnsupportedFile {
		t.Errorf("expected ErrUnsupportedFile, got %v", err)
	}
}

func TestSourceFileDoesNotExist(t *testing.T) {
	fileName := "testdata/"
	to := "testdata/input_copy.txt"
	err := Copy(fileName, to, 0, 10)
	if err != ErrUnsupportedFile {
		t.Errorf("expected ErrUnsupportedFile, got %v", err)
	}
}
