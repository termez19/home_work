package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	command := cmd[0]
	args := cmd[1:]

	execCmd := exec.Command(command, args...)
	currentEnv := os.Environ()

	modifiedEnv := applyEnvChanges(currentEnv, env)

	execCmd.Env = modifiedEnv

	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	err := execCmd.Run()

	returnCode = 0

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			returnCode = exitError.ExitCode()
		} else {
			returnCode = 1
		}
	}
	return returnCode
}

func applyEnvChanges(currentEnv []string, changes Environment) []string {

	envMap := make(map[string]string)

	for _, envVar := range currentEnv {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	for name, value := range changes {
		if value.NeedRemove {
			delete(envMap, name)
		} else {
			envMap[name] = value.Value
		}
	}

	result := make([]string, 0, len(envMap))
	for name, value := range envMap {
		result = append(result, name+"="+value)
	}

	return result
}
