package main

import (
	"fmt"
	"os"
	"strings"

	"bufio"
	"bytes"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("can't read dir %s: %w", dir, err)
	}

	env := make(Environment)
	for _, entry := range entries {
		if entry.IsDir() || strings.Contains(entry.Name(), "=") {
			continue
		}

		envValue, err := readEnvFile(dir + "/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("can't read env file %s: %w", entry.Name(), err)
		}

		env[entry.Name()] = envValue
	}
	return env, nil
}

func readEnvFile(path string) (EnvValue, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return EnvValue{}, fmt.Errorf("can't read file %s: %w", path, err)
	}

	if len(content) == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	contentScanner := bufio.NewScanner(bytes.NewReader(content))
	if !contentScanner.Scan() {
		return EnvValue{NeedRemove: true}, nil
	}

	firstLine := contentScanner.Bytes()

	firstLine = bytes.Replace(firstLine, []byte{0}, []byte{'\n'}, -1)

	value := strings.TrimRight(string(firstLine), " \t")

	if value == "" {
		return EnvValue{NeedRemove: true}, nil
	}

	return EnvValue{Value: value}, nil
}
