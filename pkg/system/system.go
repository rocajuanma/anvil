/*
Copyright © 2022 Juanma Roca juanmaxroca@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// CommandResult represents the result of a command execution
type CommandResult struct {
	Command  string
	ExitCode int
	Output   string
	Error    string
	Success  bool
}

// RunCommand executes a system command with a default timeout of 5 minutes
func RunCommand(command string, args ...string) (*CommandResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	return RunCommandWithTimeout(ctx, command, args...)
}

// RunCommandWithTimeout executes a system command with the given context
func RunCommandWithTimeout(ctx context.Context, command string, args ...string) (*CommandResult, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	// Capture both stdout and stderr
	output, err := cmd.CombinedOutput()

	result := &CommandResult{
		Command: strings.Join(append([]string{command}, args...), " "),
		Output:  string(output),
		Success: err == nil,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		}
		result.Error = err.Error()
	}

	return result, nil
}

// RunCommandWithOutput executes a command and prints output in real-time
func RunCommandWithOutput(command string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	return RunCommandWithOutputTimeout(ctx, command, args...)
}

// RunCommandWithOutputTimeout executes a command with context and prints output in real-time
func RunCommandWithOutputTimeout(ctx context.Context, command string, args ...string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CommandExists checks if a command exists in the system PATH
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// GetCommandPath returns the full path to a command if it exists
func GetCommandPath(command string) (string, error) {
	return exec.LookPath(command)
}

// RunCommandInDirectory executes a command in a specific directory
func RunCommandInDirectory(dir, command string, args ...string) (*CommandResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	return RunCommandInDirectoryWithTimeout(ctx, dir, command, args...)
}

// RunCommandInDirectoryWithTimeout executes a command in a specific directory with context
func RunCommandInDirectoryWithTimeout(ctx context.Context, dir, command string, args ...string) (*CommandResult, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()

	result := &CommandResult{
		Command: fmt.Sprintf("cd %s && %s", dir, strings.Join(append([]string{command}, args...), " ")),
		Output:  string(output),
		Success: err == nil,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		}
		result.Error = err.Error()
	}

	return result, nil
}

// GetEnvironmentVariable gets an environment variable with a default value
func GetEnvironmentVariable(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetEnvironmentVariable sets an environment variable
func SetEnvironmentVariable(key, value string) error {
	return os.Setenv(key, value)
}
