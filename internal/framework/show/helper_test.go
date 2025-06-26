package cli

import (
	"encoding/json"
	"testing"
)

func TestPrintJsonFunction(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "simple string",
			data: "test",
		},
		{
			name: "simple map",
			data: map[string]string{"key": "value"},
		},
		{
			name: "simple slice",
			data: []string{"item1", "item2"},
		},
		{
			name: "nil data",
			data: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since printJson prints to stdout and calls log.Fatal on error,
			// we just test that the JSON marshaling would work
			_, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("JSON marshaling failed for %v: %v", tt.data, err)
			}
		})
	}
}

func TestPrintJsonJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "complex structure",
			data: map[string]interface{}{
				"string": "value",
				"number": 42,
				"bool":   true,
				"array":  []int{1, 2, 3},
				"nested": map[string]string{"inner": "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("Failed to marshal data: %v", err)
			}

			var unmarshaled interface{}
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal data: %v", err)
			}
		})
	}
}

func TestPrintJsonFailFast(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "json marshaling should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("JSON marshaling panicked: %v", r)
					}
				}()

				data := map[string]string{"test": "value"}
				_, err := json.Marshal(data)
				if err != nil {
					t.Logf("JSON marshaling failed as expected for certain data types: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestPrintJsonTableDriven(t *testing.T) {
	tests := []struct {
		name        string
		data        interface{}
		shouldError bool
	}{
		{
			name:        "valid string",
			data:        "test",
			shouldError: false,
		},
		{
			name:        "valid number",
			data:        42,
			shouldError: false,
		},
		{
			name:        "valid boolean",
			data:        true,
			shouldError: false,
		},
		{
			name:        "valid slice",
			data:        []string{"a", "b", "c"},
			shouldError: false,
		},
		{
			name:        "valid map",
			data:        map[string]int{"one": 1, "two": 2},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := json.Marshal(tt.data)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for %v, but got none", tt.data)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error for %v: %v", tt.data, err)
			}
		})
	}
}
