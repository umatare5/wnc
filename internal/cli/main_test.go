package cli

import (
	"testing"
)

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
