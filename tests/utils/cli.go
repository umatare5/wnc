package utils

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// CLIResult holds the result of a CLI command execution
type CLIResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// RunCLI executes a CLI command and captures stdout/stderr
func RunCLI(t *testing.T, args ...string) *CLIResult {
	t.Helper()

	projectRoot := findProjectRoot(t)
	binaryPath := filepath.Join(projectRoot, "tmp", "wnc")

	// Build the binary if it doesn't exist
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/")
		buildCmd.Dir = projectRoot
		if output, err := buildCmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, string(output))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, args...)
	cmd.Dir = projectRoot

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	return &CLIResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Error:    err,
	}
}

// RunCLIWithEnv executes a CLI command with custom environment variables
func RunCLIWithEnv(t *testing.T, env map[string]string, args ...string) *CLIResult {
	t.Helper()

	// Save original environment
	originalEnv := make(map[string]string)
	for key := range env {
		originalEnv[key] = os.Getenv(key)
		if err := os.Setenv(key, env[key]); err != nil {
			t.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}

	defer func() {
		// Restore original environment
		for key, value := range originalEnv {
			if value == "" {
				_ = os.Unsetenv(key) // Best effort
			} else {
				_ = os.Setenv(key, value) // Best effort
			}
		}
	}()

	return RunCLI(t, args...)
}

// ExpectSuccessfulCLI runs CLI and expects successful execution
func ExpectSuccessfulCLI(t *testing.T, args ...string) *CLIResult {
	t.Helper()

	result := RunCLI(t, args...)
	if result.ExitCode != 0 {
		t.Fatalf("CLI command failed with exit code %d\nStdout: %s\nStderr: %s",
			result.ExitCode, result.Stdout, result.Stderr)
	}

	return result
}

// ExpectFailedCLI runs CLI and expects it to fail
func ExpectFailedCLI(t *testing.T, args ...string) *CLIResult {
	t.Helper()

	result := RunCLI(t, args...)
	if result.ExitCode == 0 {
		t.Fatalf("Expected CLI command to fail, but it succeeded\nStdout: %s\nStderr: %s",
			result.Stdout, result.Stderr)
	}

	return result
}

// AssertOutputContains checks if the output contains expected strings
func (r *CLIResult) AssertOutputContains(t *testing.T, expected ...string) {
	t.Helper()

	combined := r.Stdout + r.Stderr
	for _, exp := range expected {
		if !strings.Contains(combined, exp) {
			t.Errorf("Expected output to contain %q\nActual output: %s", exp, combined)
		}
	}
}

// AssertStdoutContains checks if stdout contains expected strings
func (r *CLIResult) AssertStdoutContains(t *testing.T, expected ...string) {
	t.Helper()

	for _, exp := range expected {
		if !strings.Contains(r.Stdout, exp) {
			t.Errorf("Expected stdout to contain %q\nActual stdout: %s", exp, r.Stdout)
		}
	}
}

// AssertStderrContains checks if stderr contains expected strings
func (r *CLIResult) AssertStderrContains(t *testing.T, expected ...string) {
	t.Helper()

	for _, exp := range expected {
		if !strings.Contains(r.Stderr, exp) {
			t.Errorf("Expected stderr to contain %q\nActual stderr: %s", exp, r.Stderr)
		}
	}
}

// AssertExitCode checks if the exit code matches expected value
func (r *CLIResult) AssertExitCode(t *testing.T, expectedCode int) {
	t.Helper()

	if r.ExitCode != expectedCode {
		t.Errorf("Expected exit code %d, got %d\nStdout: %s\nStderr: %s",
			expectedCode, r.ExitCode, r.Stdout, r.Stderr)
	}
}

// BuildCLI builds the CLI binary for testing
func BuildCLI(t *testing.T) string {
	t.Helper()

	binaryPath := "./tmp/wnc"
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/")
	buildCmd.Dir = findProjectRoot(t)

	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, string(output))
	}

	return binaryPath
}

// findProjectRoot finds the project root directory containing go.mod
func findProjectRoot(t *testing.T) string {
	t.Helper()

	// Start from current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Look for go.mod file walking up the directory tree
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the root directory
			break
		}
		dir = parent
	}

	// If go.mod not found, assume current directory
	return wd
}
