package cli

import (
	"os"
	"testing"
)

// Local test utilities
func assertNoPanic(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Function panicked: %v", r)
		}
	}()
	fn()
}

// TestRegisterSubCommands tests CLI subcommand registration (CLI test)
func TestRegisterSubCommands(t *testing.T) {
	tests := []struct {
		name            string
		wantMinCommands int
	}{
		{
			name:            "registers generate and show commands",
			wantMinCommands: 2, // At least generate and show commands
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := registerSubCommands()

			if len(got) < tt.wantMinCommands {
				t.Errorf("registerSubCommands() returned %d commands, want at least %d",
					len(got), tt.wantMinCommands)
			}

			// Check that commands are not nil
			for i, cmd := range got {
				if cmd == nil {
					t.Errorf("Command at index %d is nil", i)
				}
			}

			// Verify that expected command names exist
			commandNames := make(map[string]bool)
			for _, cmd := range got {
				if cmd != nil {
					commandNames[cmd.Name] = true
				}
			}

			expectedCommands := []string{"generate", "show"}
			for _, expectedCmd := range expectedCommands {
				if !commandNames[expectedCmd] {
					t.Errorf("Expected command %q not found in registered commands", expectedCmd)
				}
			}
		})
	}
}

// TestGetVersion tests the version function
func TestGetVersion(t *testing.T) {
	t.Run("getVersion_returns_non_empty_string", func(t *testing.T) {
		version := getVersion()
		if len(version) == 0 {
			t.Error("getVersion() returned empty string")
		}
	})
}

// TestRun tests the main CLI run function
func TestRun(t *testing.T) {
	t.Run("run_with_help_flag", func(t *testing.T) {
		// Test with help flag
		assertNoPanic(t, func() {
			// Save original args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set args to help
			os.Args = []string{"wnc", "--help"}

			// This should not panic
			Run()
		})
	})

	t.Run("run_with_version_flag", func(t *testing.T) {
		// Test with version flag
		assertNoPanic(t, func() {
			// Save original args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set args to version
			os.Args = []string{"wnc", "--version"}

			// This should not panic
			Run()
		})
	})

	t.Run("run_with_invalid_command", func(t *testing.T) {
		// Test with invalid command
		assertNoPanic(t, func() {
			// Save original args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set args to invalid command
			os.Args = []string{"wnc", "invalid-command"}

			// This should not panic
			Run()
		})
	})

	t.Run("run_with_no_args", func(t *testing.T) {
		// Test with no arguments (should show help)
		assertNoPanic(t, func() {
			// Save original args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set args to just program name
			os.Args = []string{"wnc"}

			// This should not panic
			Run()
		})
	})
}

// TestCLICommandStructure tests the overall CLI command structure
func TestCLICommandStructure(t *testing.T) {
	t.Run("cli_commands_have_proper_structure", func(t *testing.T) {
		commands := registerSubCommands()

		for _, cmd := range commands {
			if cmd == nil {
				continue
			}

			// Verify command has a name
			if len(cmd.Name) == 0 {
				t.Errorf("Command has empty name")
			}

			// Verify command has description (some commands may have empty descriptions)
			if len(cmd.Description) == 0 && len(cmd.Usage) == 0 {
				t.Logf("Command %s has both empty description and usage", cmd.Name)
			}

			// Verify command has action or subcommands
			if cmd.Action == nil && len(cmd.Commands) == 0 {
				t.Errorf("Command %s has neither action nor subcommands", cmd.Name)
			}
		}
	})
}
