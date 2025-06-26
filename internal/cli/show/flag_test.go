package subcommand

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/umatare5/wnc/internal/config"
	"github.com/urfave/cli/v3"
)

func TestRegisterControllersFlag(t *testing.T) {
	t.Run("registers controllers flag with correct properties", func(t *testing.T) {
		flags := registerControllersFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if flag.Name != config.ControllersFlagName {
			t.Errorf("Expected name %s, got %s", config.ControllersFlagName, flag.Name)
		}

		if !flag.Required {
			t.Error("Expected flag to be required")
		}

		if len(flag.Aliases) == 0 || flag.Aliases[0] != "c" {
			t.Error("Expected alias 'c'")
		}
	})
}

func TestRegisterPrintFormatFlag(t *testing.T) {
	t.Run("registers print format flag with correct properties", func(t *testing.T) {
		flags := registerPrintFormatFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if flag.Name != config.PrintFormatFlagName {
			t.Errorf("Expected name %s, got %s", config.PrintFormatFlagName, flag.Name)
		}

		if flag.Value != config.PrintFormatTable {
			t.Errorf("Expected default value %s, got %s", config.PrintFormatTable, flag.Value)
		}

		if !strings.Contains(flag.Usage, config.PrintFormatJSON) {
			t.Error("Expected usage to contain JSON format")
		}
	})
}

func TestRegisterTimeoutFlag(t *testing.T) {
	t.Run("registers timeout flag with correct properties", func(t *testing.T) {
		flags := registerTimeoutFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.IntFlag)
		if !ok {
			t.Error("Expected IntFlag")
		}

		if flag.Name != config.TimeoutFlagName {
			t.Errorf("Expected name %s, got %s", config.TimeoutFlagName, flag.Name)
		}

		if flag.Value != 60 {
			t.Errorf("Expected default value 60, got %d", flag.Value)
		}
	})
}

func TestRegisterRadioFlag(t *testing.T) {
	t.Run("registers radio flag with correct properties", func(t *testing.T) {
		flags := registerRadioFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if flag.Name != config.RadioFlagName {
			t.Errorf("Expected name %s, got %s", config.RadioFlagName, flag.Name)
		}

		if !strings.Contains(flag.Usage, "2.4GHz") {
			t.Error("Expected usage to contain radio frequency information")
		}
	})
}

func TestRegisterSSIDFlag(t *testing.T) {
	t.Run("registers SSID flag with correct properties", func(t *testing.T) {
		flags := registerSSIDFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if flag.Name != config.SSIDFlagName {
			t.Errorf("Expected name %s, got %s", config.SSIDFlagName, flag.Name)
		}
	})
}

func TestRegisterInsecureFlag(t *testing.T) {
	t.Run("registers insecure flag with correct properties", func(t *testing.T) {
		flags := registerInsecureFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.BoolFlag)
		if !ok {
			t.Error("Expected BoolFlag")
		}

		if flag.Name != config.AllowInsecureAccessFlagName {
			t.Errorf("Expected name %s, got %s", config.AllowInsecureAccessFlagName, flag.Name)
		}

		if flag.Value != false {
			t.Error("Expected default value false")
		}
	})
}

func TestFlagJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "flag configurations",
			data: map[string]interface{}{
				"controllersFlag": len(registerControllersFlag()),
				"printFormatFlag": len(registerPrintFormatFlag()),
				"timeoutFlag":     len(registerTimeoutFlag()),
				"radioFlag":       len(registerRadioFlag()),
				"ssidFlag":        len(registerSSIDFlag()),
				"insecureFlag":    len(registerInsecureFlag()),
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

func TestFlagFailFast(t *testing.T) {
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

				_ = registerControllersFlag()
				_ = registerPrintFormatFlag()
				_ = registerTimeoutFlag()
				_ = registerRadioFlag()
				_ = registerSSIDFlag()
				_ = registerInsecureFlag()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestFlagTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		flagFunc func() []cli.Flag
		expected int
	}{
		{
			name:     "controllers flag returns one flag",
			flagFunc: registerControllersFlag,
			expected: 1,
		},
		{
			name:     "print format flag returns one flag",
			flagFunc: registerPrintFormatFlag,
			expected: 1,
		},
		{
			name:     "timeout flag returns one flag",
			flagFunc: registerTimeoutFlag,
			expected: 1,
		},
		{
			name:     "radio flag returns one flag",
			flagFunc: registerRadioFlag,
			expected: 1,
		},
		{
			name:     "ssid flag returns one flag",
			flagFunc: registerSSIDFlag,
			expected: 1,
		},
		{
			name:     "insecure flag returns one flag",
			flagFunc: registerInsecureFlag,
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

func TestSortByFlags(t *testing.T) {
	t.Run("client sort by flag contains expected options", func(t *testing.T) {
		flags := registerClientSortByFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if !strings.Contains(flag.Usage, config.ShowClientHeaderHostname) {
			t.Error("Expected usage to contain hostname option")
		}
	})

	t.Run("overview sort by flag contains expected options", func(t *testing.T) {
		flags := registerOverviewSortByFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if !strings.Contains(flag.Usage, config.ShowCommonHeaderApName) {
			t.Error("Expected usage to contain AP name option")
		}
	})

	t.Run("sort order flag contains expected options", func(t *testing.T) {
		flags := registerSortOrderFlag()

		if len(flags) != 1 {
			t.Errorf("Expected 1 flag, got %d", len(flags))
		}

		flag, ok := flags[0].(*cli.StringFlag)
		if !ok {
			t.Error("Expected StringFlag")
		}

		if !strings.Contains(flag.Usage, config.OrderByAscending) {
			t.Error("Expected usage to contain ascending option")
		}

		if !strings.Contains(flag.Usage, config.OrderByDescending) {
			t.Error("Expected usage to contain descending option")
		}
	})
}
