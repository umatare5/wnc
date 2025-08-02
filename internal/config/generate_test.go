package config

import (
	"encoding/json"
	"testing"
)

// TestGenerateCmdConfigJSONSerialization tests JSON serialization for GenerateCmdConfig (Unit test)
func TestGenerateCmdConfigJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		config GenerateCmdConfig
	}{
		{
			name: "empty_generate_config",
			config: GenerateCmdConfig{
				Username: "",
				Password: "",
			},
		},
		{
			name: "full_generate_config",
			config: GenerateCmdConfig{
				Username: "testuser",
				Password: "testpassword",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.config)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Test JSON unmarshaling
			var config GenerateCmdConfig
			err = json.Unmarshal(jsonData, &config)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Basic validation
			if config.Username != tt.config.Username {
				t.Errorf("Username mismatch: got %s, want %s", config.Username, tt.config.Username)
			}
			if config.Password != tt.config.Password {
				t.Errorf("Password mismatch: got %s, want %s", config.Password, tt.config.Password)
			}
		})
	}
}
