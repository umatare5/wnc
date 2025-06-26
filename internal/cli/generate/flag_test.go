package subcommand

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
	"github.com/urfave/cli/v3"
)

func TestRegisterUsernameFlag(t *testing.T) {
	t.Run("registers username flag with correct properties", func(t *testing.T) {
		flags := registerUsernameFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if flag.Name != config.UsernameFlagName {
			t.Errorf("Expected name %s, got %s", config.UsernameFlagName, flag.Name)
		}

		if flag.Value != "" {
			t.Errorf("Expected empty default value, got %s", flag.Value)
		}

		if len(flag.Aliases) == 0 || flag.Aliases[0] != "u" {
			t.Error("Expected alias 'u'")
		}

		if flag.Usage == "" {
			t.Error("Expected non-empty usage")
		}
	})
}

func TestRegisterPasswordFlag(t *testing.T) {
	t.Run("registers password flag with correct properties", func(t *testing.T) {
		flags := registerPasswordFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if flag.Name != config.PasswordFlagName {
			t.Errorf("Expected name %s, got %s", config.PasswordFlagName, flag.Name)
		}

		if flag.Value != "" {
			t.Errorf("Expected empty default value, got %s", flag.Value)
		}

		if len(flag.Aliases) == 0 || flag.Aliases[0] != "p" {
			t.Error("Expected alias 'p'")
		}

		if flag.Usage == "" {
			t.Error("Expected non-empty usage")
		}
	})
}

func TestGenerateFlagJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "flag configurations",
			data: map[string]interface{}{
				"usernameFlag": len(registerUsernameFlag()),
				"passwordFlag": len(registerPasswordFlag()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("Failed to marshal flag data: %v", err)
			}

			var unmarshaled map[string]interface{}
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal flag data: %v", err)
			}
		})
	}
}

func TestGenerateFlagFailFast(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "flag registration should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Flag registration panicked: %v", r)
					}
				}()

				_ = registerUsernameFlag()
				_ = registerPasswordFlag()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestGenerateFlagTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		flagFunc func() []cli.Flag
		expected int
	}{
		{
			name:     "username flag returns one flag",
			flagFunc: registerUsernameFlag,
			expected: 1,
		},
		{
			name:     "password flag returns one flag",
			flagFunc: registerPasswordFlag,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := tt.flagFunc()
			if len(flags) != tt.expected {
				t.Errorf("Expected %d flags, got %d", tt.expected, len(flags))
			}
		})
	}
}

func TestGenerateFlagValidation(t *testing.T) {
	tests := []struct {
		name     string
		flagFunc func() []cli.Flag
		flagName string
		alias    string
	}{
		{
			name:     "username flag has correct name and alias",
			flagFunc: registerUsernameFlag,
			flagName: config.UsernameFlagName,
			alias:    "u",
		},
		{
			name:     "password flag has correct name and alias",
			flagFunc: registerPasswordFlag,
			flagName: config.PasswordFlagName,
			alias:    "p",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := tt.flagFunc()

			if len(flags) != 1 {
				t.Errorf("Expected 1 flag, got %d", len(flags))
				return
			}

			flag, ok := flags[0].(*cli.StringFlag)
			if !ok {
				t.Error("Expected StringFlag")
				return
			}

			if flag.Name != tt.flagName {
				t.Errorf("Expected name %s, got %s", tt.flagName, flag.Name)
			}

			if len(flag.Aliases) == 0 || flag.Aliases[0] != tt.alias {
				t.Errorf("Expected alias %s, got %v", tt.alias, flag.Aliases)
			}
		})
	}
}
