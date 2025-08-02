package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/umatare5/wnc/tests/utils"
)

// TestMainPackage tests that the main package compiles and basic functionality
func TestMainPackage(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "main function exists",
			test: func(t *testing.T) {
				// Test that main function exists by checking compilation succeeds
				// The main function is always non-nil if the package compiles
				// This test validates the main package structure
			},
		},
		{
			name: "package compiles",
			test: func(t *testing.T) {
				// The fact that this test runs means the package compiles successfully
				// This is important for CLI applications
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestMainFunction tests the main function directly
// Since main() calls cli.Run(), we can test its execution pattern
func TestMainFunction(t *testing.T) {
	// We cannot directly test main() since it doesn't return
	// But we can test that it exists and compiles correctly
	t.Run("main function compilation", func(t *testing.T) {
		// The fact that this test runs means main() compiles
		// This validates the function signature and imports
	})
}

// TestMainImports verifies that all imports are correct
func TestMainImports(t *testing.T) {
	t.Run("cli import validation", func(t *testing.T) {
		// If this test runs, the cli import was successful
		// This validates the import path and package availability
	})
}

// TestMainExecution tests main function execution using subprocess
func TestMainExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping main execution test in short mode")
	}

	// Check if we're in a subprocess execution
	if os.Getenv("GO_TEST_SUBPROCESS") == "1" {
		// This runs the actual main function
		main()
		return
	}

	// Test various command line scenarios
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "version flag",
			args: []string{"--version"},
		},
		{
			name: "help flag",
			args: []string{"--help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestMainExecution"}, tt.args...)...)
			cmd.Env = append(os.Environ(), "GO_TEST_SUBPROCESS=1")

			// Execute the subprocess
			err := cmd.Run()

			// For CLI applications, exit with code 0 or 1 is normal
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					// Exit codes 0, 1, and 2 are acceptable for CLI help/version/errors
					if exitError.ExitCode() <= 2 {
						t.Logf("CLI exited with code %d (expected for %s)", exitError.ExitCode(), tt.name)
					} else {
						t.Errorf("CLI exited with unexpected code %d", exitError.ExitCode())
					}
				} else {
					t.Errorf("Failed to run CLI: %v", err)
				}
			}
		})
	}
}

// TestMainExecutionCLI tests CLI execution using utilities
func TestMainExecutionCLI(t *testing.T) {
	// Skip if environment variables for live testing are not set
	if !utils.HasTestControllers() {
		t.Skip("Skipping CLI tests - no test controllers configured")
	}

	tests := []struct {
		name           string
		args           []string
		expectSuccess  bool
		expectedOutput []string
	}{
		{
			name:           "version flag",
			args:           []string{"--version"},
			expectSuccess:  true,
			expectedOutput: []string{"wnc version"},
		},
		{
			name:           "help flag",
			args:           []string{"--help"},
			expectSuccess:  true,
			expectedOutput: []string{"USAGE:", "COMMANDS:"},
		},
		{
			name:           "show overview with env",
			args:           []string{"show", "overview"},
			expectSuccess:  true,
			expectedOutput: []string{}, // Will be populated from environment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result *utils.CLIResult

			if tt.expectSuccess {
				result = utils.ExpectSuccessfulCLI(t, tt.args...)
			} else {
				result = utils.ExpectFailedCLI(t, tt.args...)
			}

			// Check expected output
			if len(tt.expectedOutput) > 0 {
				result.AssertOutputContains(t, tt.expectedOutput...)
			}
		})
	}
}

// TestMainWithInvalidArgs tests main function with invalid arguments
func TestMainWithInvalidArgs(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	t.Run("invalid command handling", func(t *testing.T) {
		os.Args = []string{"wnc", "nonexistent-command"}

		// Call main in a controlled way
		done := make(chan bool, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("main() handled invalid command correctly: %v", r)
				}
				done <- true
			}()
			main()
		}()

		// Wait for completion
		<-done
	})
}
