package log

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		name           string
		logLevel       string
		expectedLevel  logrus.Level
		expectWarnLogs bool
	}{
		{
			name:          "set_warn_level",
			logLevel:      "warn",
			expectedLevel: logrus.WarnLevel,
		},
		{
			name:          "set_error_level",
			logLevel:      "error",
			expectedLevel: logrus.ErrorLevel,
		},
		{
			name:          "set_debug_level",
			logLevel:      "debug",
			expectedLevel: logrus.DebugLevel,
		},
		{
			name:          "set_info_level_default",
			logLevel:      "info",
			expectedLevel: logrus.InfoLevel,
		},
		{
			name:          "set_invalid_level_defaults_to_info",
			logLevel:      "invalid",
			expectedLevel: logrus.InfoLevel,
		},
		{
			name:          "empty_level_defaults_to_info",
			logLevel:      "",
			expectedLevel: logrus.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the log level
			SetLogLevel(tt.logLevel)

			// Check that the level was set correctly
			if logger.GetLevel() != tt.expectedLevel {
				t.Errorf("Expected log level %v, got %v", tt.expectedLevel, logger.GetLevel())
			}
		})
	}
}

func TestLoggerFunctions(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.DebugLevel) // Set to debug to capture all levels

	tests := []struct {
		name     string
		logFunc  func()
		contains string
	}{
		{
			name: "Info_function",
			logFunc: func() {
				Info("test info message")
			},
			contains: "test info message",
		},
		{
			name: "Infof_function",
			logFunc: func() {
				Infof("test %s message", "infof")
			},
			contains: "test infof message",
		},
		{
			name: "Warnf_function",
			logFunc: func() {
				Warnf("test %s message", "warnf")
			},
			contains: "test warnf message",
		},
		{
			name: "Errorf_function",
			logFunc: func() {
				Errorf("test %s message", "errorf")
			},
			contains: "test errorf message",
		},
		{
			name: "Debugf_function",
			logFunc: func() {
				Debugf("test %s message", "debugf")
			},
			contains: "test debugf message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()

			output := buf.String()
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Expected log output to contain '%s', got: %s", tt.contains, output)
			}
		})
	}
}

func TestLoggerJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		config struct {
			LogLevel string `json:"log_level"`
		}
	}{
		{
			name: "log_config_serialization",
			config: struct {
				LogLevel string `json:"log_level"`
			}{
				LogLevel: "debug",
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
			var config struct {
				LogLevel string `json:"log_level"`
			}
			err = json.Unmarshal(jsonData, &config)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if config.LogLevel != tt.config.LogLevel {
				t.Errorf("Expected LogLevel %s, got %s", tt.config.LogLevel, config.LogLevel)
			}
		})
	}
}

func TestLoggerFailFast(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
	}{
		{
			name:     "valid_log_level_should_not_panic",
			logLevel: "debug",
		},
		{
			name:     "invalid_log_level_should_not_panic",
			logLevel: "invalid",
		},
		{
			name:     "empty_log_level_should_not_panic",
			logLevel: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Logger operation panicked: %v", r)
				}
			}()

			// Test that setting log level doesn't panic
			SetLogLevel(tt.logLevel)

			// Test that logging functions don't panic
			var buf bytes.Buffer
			logger.SetOutput(&buf)

			Info("test message")
			Infof("test %s", "message")
			Warnf("test %s", "message")
			Errorf("test %s", "message")
			Debugf("test %s", "message")
		})
	}
}

func TestLoggerTableDriven(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel logrus.Level
		message       string
		shouldLog     bool
	}{
		{
			name:          "debug_level_logs_debug",
			level:         "debug",
			expectedLevel: logrus.DebugLevel,
			message:       "debug message",
			shouldLog:     true,
		},
		{
			name:          "info_level_logs_info",
			level:         "info",
			expectedLevel: logrus.InfoLevel,
			message:       "info message",
			shouldLog:     true,
		},
		{
			name:          "warn_level_logs_warn",
			level:         "warn",
			expectedLevel: logrus.WarnLevel,
			message:       "warn message",
			shouldLog:     true,
		},
		{
			name:          "error_level_logs_error",
			level:         "error",
			expectedLevel: logrus.ErrorLevel,
			message:       "error message",
			shouldLog:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger.SetOutput(&buf)

			SetLogLevel(tt.level)

			if logger.GetLevel() != tt.expectedLevel {
				t.Errorf("Expected level %v, got %v", tt.expectedLevel, logger.GetLevel())
			}

			// Test appropriate logging function
			switch tt.level {
			case "debug":
				Debugf("%s", tt.message)
			case "info":
				Infof("%s", tt.message)
			case "warn":
				Warnf("%s", tt.message)
			case "error":
				Errorf("%s", tt.message)
			}

			output := buf.String()
			containsMessage := strings.Contains(output, tt.message)

			if tt.shouldLog && !containsMessage {
				t.Errorf("Expected log output to contain '%s'", tt.message)
			}
		})
	}
}
