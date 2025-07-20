package main

import (
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple echo command", func(t *testing.T) {
		env := Environment{
			"TEST_VAR": EnvValue{Value: "test_value"},
		}

		cmd := []string{"sh", "-c", "echo $TEST_VAR"}
		exitCode := RunCmd(cmd, env)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	t.Run("command with arguments", func(t *testing.T) {
		env := Environment{}
		cmd := []string{"echo", "hello", "world"}
		exitCode := RunCmd(cmd, env)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	t.Run("nonexistent command", func(t *testing.T) {
		env := Environment{}
		cmd := []string{"nonexistent_command_12345"}
		exitCode := RunCmd(cmd, env)

		if exitCode == 0 {
			t.Error("Expected non-zero exit code for nonexistent command")
		}
	})

	t.Run("command with non-zero exit", func(t *testing.T) {
		env := Environment{}
		cmd := []string{"sh", "-c", "exit 42"}
		exitCode := RunCmd(cmd, env)

		if exitCode != 42 {
			t.Errorf("Expected exit code 42, got %d", exitCode)
		}
	})

	t.Run("environment variable removal", func(t *testing.T) {
		os.Setenv("TEMP_TEST_VAR", "should_be_removed")
		defer os.Unsetenv("TEMP_TEST_VAR")

		env := Environment{
			"TEMP_TEST_VAR": EnvValue{NeedRemove: true},
		}

		cmd := []string{"sh", "-c", "[ -z \"$TEMP_TEST_VAR\" ] && echo 'removed' || echo 'not_removed'"}
		exitCode := RunCmd(cmd, env)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	t.Run("environment variable setting", func(t *testing.T) {
		env := Environment{
			"NEW_VAR": EnvValue{Value: "new_value"},
		}

		cmd := []string{"sh", "-c", "[ \"$NEW_VAR\" = \"new_value\" ] && exit 0 || exit 1"}
		exitCode := RunCmd(cmd, env)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})

	t.Run("empty command slice", func(t *testing.T) {
		env := Environment{}
		cmd := []string{}
		exitCode := RunCmd(cmd, env)

		if exitCode == 0 {
			t.Error("Expected non-zero exit code for empty command")
		}
	})

	t.Run("preserve existing environment", func(t *testing.T) {
		os.Setenv("EXISTING_VAR", "existing_value")
		defer os.Unsetenv("EXISTING_VAR")

		env := Environment{
			"NEW_VAR": EnvValue{Value: "new_value"},
		}

		cmd := []string{"sh", "-c", "[ \"$EXISTING_VAR\" = \"existing_value\" ] && [ \"$NEW_VAR\" = \"new_value\" ]"}
		exitCode := RunCmd(cmd, env)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
	})
}
