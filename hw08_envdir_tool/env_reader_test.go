package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("existing testdata", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		if err != nil {
			t.Fatalf("ReadDir failed: %v", err)
		}

		tests := []struct {
			name     string
			expected EnvValue
		}{
			{"FOO", EnvValue{Value: "   foo\nwith new line"}},
			{"BAR", EnvValue{Value: "bar"}},
			{"HELLO", EnvValue{Value: "\"hello\""}},
			{"EMPTY", EnvValue{NeedRemove: true}},
			{"UNSET", EnvValue{NeedRemove: true}},
		}

		for _, test := range tests {
			actual, exists := env[test.name]
			if !exists {
				t.Errorf("Variable %s not found", test.name)
				continue
			}

			if actual.NeedRemove != test.expected.NeedRemove {
				t.Errorf("Variable %s: NeedRemove = %v, want %v",
					test.name, actual.NeedRemove, test.expected.NeedRemove)
			}

			if actual.Value != test.expected.Value {
				t.Errorf("Variable %s: Value = %q, want %q",
					test.name, actual.Value, test.expected.Value)
			}
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		_, err := ReadDir("nonexistent")
		if err == nil {
			t.Error("Expected error for nonexistent directory")
		}
	})

	t.Run("custom test directory", func(t *testing.T) {
		tempDir := t.TempDir()

		createTestFile(t, tempDir, "SIMPLE", "value")
		createTestFile(t, tempDir, "WITH_SPACES", "  trimmed  ")
		createTestFile(t, tempDir, "MULTILINE", "first\nsecond\nthird")
		createTestFile(t, tempDir, "EMPTY_FILE", "")
		createTestFile(t, tempDir, "ONLY_SPACES", "   \t  ")
		createTestFile(t, tempDir, "WITH_NULLS", "hello\x00world")
		createTestFile(t, tempDir, "INVALID=NAME", "should_be_ignored")

		env, err := ReadDir(tempDir)
		if err != nil {
			t.Fatalf("ReadDir failed: %v", err)
		}

		tests := []struct {
			name     string
			expected EnvValue
		}{
			{"SIMPLE", EnvValue{Value: "value"}},
			{"WITH_SPACES", EnvValue{Value: "  trimmed"}},
			{"MULTILINE", EnvValue{Value: "first"}},
			{"EMPTY_FILE", EnvValue{NeedRemove: true}},
			{"ONLY_SPACES", EnvValue{NeedRemove: true}},
			{"WITH_NULLS", EnvValue{Value: "hello\nworld"}},
		}

		for _, test := range tests {
			actual, exists := env[test.name]
			if !exists {
				t.Errorf("Variable %s not found", test.name)
				continue
			}

			if actual.NeedRemove != test.expected.NeedRemove {
				t.Errorf("Variable %s: NeedRemove = %v, want %v",
					test.name, actual.NeedRemove, test.expected.NeedRemove)
			}

			if actual.Value != test.expected.Value {
				t.Errorf("Variable %s: Value = %q, want %q",
					test.name, actual.Value, test.expected.Value)
			}
		}

		if _, exists := env["INVALID=NAME"]; exists {
			t.Error("File with = in name should be ignored")
		}
	})
}

func createTestFile(t *testing.T, dir, name, content string) {
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", path, err)
	}
}
