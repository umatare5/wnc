package show

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TestApTagCli_JSON tests JSON serialization and deserialization
func TestApTagCli_JSON(t *testing.T) {
	tests := []struct {
		name string
		data ApTagCli
	}{
		{
			name: "valid ApTagCli struct",
			data: ApTagCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("JSON marshaling failed: %v", err)
				return
			}

			// Test JSON unmarshaling
			var unmarshaled ApTagCli
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("JSON unmarshaling failed: %v", err)
			}
		})
	}
}

// TestApTagCli_GetShowApTagTableHeaders tests the getShowApTagTableHeaders method
func TestApTagCli_GetShowApTagTableHeaders(t *testing.T) {
	tests := []struct {
		name     string
		cli      *ApTagCli
		expected []string
	}{
		{
			name: "valid headers",
			cli: &ApTagCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			expected: []string{
				"AP Name", "Config", "Policy Tag Name", "RF Tag Name", "Site Tag Name",
				"AP Profile", "Flex Profile", "Tag Source",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("getShowApTagTableHeaders panicked: %v", r)
				}
			}()

			result := tt.cli.getShowApTagTableHeaders()
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d headers, got %d", len(tt.expected), len(result))
				return
			}

			for i, header := range result {
				if header != tt.expected[i] {
					t.Errorf("expected header[%d] = %q, got %q", i, tt.expected[i], header)
				}
			}
		})
	}
}

// TestApTagCli_FormatShowApTagRow tests the formatShowApTagRow method
func TestApTagCli_FormatShowApTagRow(t *testing.T) {
	tests := []struct {
		name string
		cli  *ApTagCli
		ap   *application.ShowApTagData
	}{
		{
			name: "nil ap tag data",
			cli: &ApTagCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			ap: nil,
		},
		{
			name: "empty ap tag data",
			cli: &ApTagCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			ap: &application.ShowApTagData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("formatShowApTagRow panicked: %v", r)
				}
			}()

			if tt.ap == nil {
				return // Skip nil test case
			}

			row, err := tt.cli.formatShowApTagRow(tt.ap)
			if err != nil {
				t.Errorf("formatShowApTagRow returned error: %v", err)
				return
			}

			if len(row) == 0 {
				t.Error("formatShowApTagRow returned empty row")
			}
		})
	}
}

// TestApTagCli_SortShowApTagRow tests the sortShowApTagRow method
func TestApTagCli_SortShowApTagRow(t *testing.T) {
	tests := []struct {
		name   string
		cli    *ApTagCli
		apTags []*application.ShowApTagData
	}{
		{
			name: "nil slice",
			cli: &ApTagCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			apTags: nil,
		},
		{
			name: "empty slice",
			cli: &ApTagCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			apTags: []*application.ShowApTagData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("sortShowApTagRow panicked: %v", r)
				}
			}()

			tt.cli.sortShowApTagRow(tt.apTags)
		})
	}
}

// TestApTagCli_ConvertCapwapTagInfoIsApMisconfigurationToConfigCheck tests the convertCapwapTagInfoIsApMisconfigurationToConfigCheck method
func TestApTagCli_ConvertCapwapTagInfoIsApMisconfigurationToConfigCheck(t *testing.T) {
	tests := []struct {
		name            string
		cli             *ApTagCli
		isMisconfigured bool
		expected        string
	}{
		{
			name: "ap not misconfigured",
			cli: &ApTagCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			isMisconfigured: false,
			expected:        "  ✅️",
		},
		{
			name: "ap misconfigured",
			cli: &ApTagCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			isMisconfigured: true,
			expected:        "  ❌️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertCapwapTagInfoIsApMisconfigurationToConfigCheck panicked: %v", r)
				}
			}()

			result := tt.cli.convertCapwapTagInfoIsApMisconfigurationToConfigCheck(tt.isMisconfigured)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
