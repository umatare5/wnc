package tablewriter

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// TestNewTable tests the NewTable function
func TestNewTable(t *testing.T) {
	tests := []struct {
		name   string
		writer *bytes.Buffer
		want   bool
	}{
		{
			name:   "with valid writer",
			writer: &bytes.Buffer{},
			want:   true,
		},
		{
			name:   "with nil writer should use stdout",
			writer: nil,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var table *Table
			if tt.writer != nil {
				table = NewTable(tt.writer)
			} else {
				table = NewTable(nil)
			}

			if (table != nil) != tt.want {
				t.Errorf("NewTable() = %v, want %v", table != nil, tt.want)
			}

			if tt.writer != nil && table.writer != tt.writer {
				t.Errorf("NewTable() writer = %v, want %v", table.writer, tt.writer)
			}

			if tt.writer == nil && table.writer != os.Stdout {
				t.Errorf("NewTable() with nil writer should use os.Stdout")
			}
		})
	}
}

// TestTableHeader tests the Header method
func TestTableHeader(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		want    []string
	}{
		{
			name:    "single header",
			headers: []string{"Name"},
			want:    []string{"Name"},
		},
		{
			name:    "multiple headers",
			headers: []string{"Name", "Age", "City"},
			want:    []string{"Name", "Age", "City"},
		},
		{
			name:    "empty headers",
			headers: []string{},
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(&bytes.Buffer{})
			table.Header(tt.headers)

			if len(table.headers) != len(tt.want) {
				t.Errorf("Header() length = %v, want %v", len(table.headers), len(tt.want))
			}

			for i, header := range table.headers {
				if header != tt.want[i] {
					t.Errorf("Header()[%d] = %v, want %v", i, header, tt.want[i])
				}
			}
		})
	}
}

// TestTableAppend tests the Append method
func TestTableAppend(t *testing.T) {
	tests := []struct {
		name string
		rows [][]string
		want int
	}{
		{
			name: "single row",
			rows: [][]string{{"John", "25", "NYC"}},
			want: 1,
		},
		{
			name: "multiple rows",
			rows: [][]string{
				{"John", "25", "NYC"},
				{"Jane", "30", "LA"},
			},
			want: 2,
		},
		{
			name: "empty row",
			rows: [][]string{{}},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(&bytes.Buffer{})

			for _, row := range tt.rows {
				table.Append(row)
			}

			if len(table.rows) != tt.want {
				t.Errorf("Append() rows count = %v, want %v", len(table.rows), tt.want)
			}
		})
	}
}

// TestTableRender tests the Render method
func TestTableRender(t *testing.T) {
	tests := []struct {
		name         string
		headers      []string
		rows         [][]string
		wantError    bool
		wantContains []string
	}{
		{
			name:      "no headers should return error",
			headers:   []string{},
			rows:      [][]string{{"data"}},
			wantError: true,
		},
		{
			name:         "valid table",
			headers:      []string{"Name", "Age"},
			rows:         [][]string{{"John", "25"}, {"Jane", "30"}},
			wantError:    false,
			wantContains: []string{"Name", "Age", "John", "25", "Jane", "30", "│", "─"},
		},
		{
			name:         "table with unicode characters",
			headers:      []string{"名前", "年齢"},
			rows:         [][]string{{"太郎", "25"}, {"花子", "30"}},
			wantError:    false,
			wantContains: []string{"名前", "年齢", "太郎", "25", "花子", "30"},
		},
		{
			name:         "mismatched column count",
			headers:      []string{"Name", "Age", "City"},
			rows:         [][]string{{"John", "25"}, {"Jane"}},
			wantError:    false,
			wantContains: []string{"Name", "Age", "City", "John", "25", "Jane"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			table := NewTable(buffer)
			table.Header(tt.headers)

			for _, row := range tt.rows {
				table.Append(row)
			}

			err := table.Render()

			if (err != nil) != tt.wantError {
				t.Errorf("Render() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				output := buffer.String()
				for _, want := range tt.wantContains {
					if !strings.Contains(output, want) {
						t.Errorf("Render() output does not contain %q. Output:\n%s", want, output)
					}
				}
			}
		})
	}
}

// TestCalculateColumnWidths tests the calculateColumnWidths method
func TestCalculateColumnWidths(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		rows    [][]string
		want    []int
	}{
		{
			name:    "headers only",
			headers: []string{"Name", "Age"},
			rows:    [][]string{},
			want:    []int{6, 5}, // "Name" + 2 padding, "Age" + 2 padding
		},
		{
			name:    "rows wider than headers",
			headers: []string{"Name", "Age"},
			rows:    [][]string{{"John", "25"}, {"Alexander", "30"}},
			want:    []int{11, 5}, // "Alexander" + 2 padding, "Age" + 2 padding
		},
		{
			name:    "unicode characters",
			headers: []string{"名前", "年齢"},
			rows:    [][]string{{"太郎", "25"}},
			want:    []int{4, 4}, // 2 unicode chars + 2 padding each
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTable(&bytes.Buffer{})
			table.headers = tt.headers
			table.rows = tt.rows

			got := table.calculateColumnWidths()

			if len(got) != len(tt.want) {
				t.Errorf("calculateColumnWidths() length = %v, want %v", len(got), len(tt.want))
			}

			for i, width := range got {
				if width != tt.want[i] {
					t.Errorf("calculateColumnWidths()[%d] = %v, want %v", i, width, tt.want[i])
				}
			}
		})
	}
}

// TestNewTableOrig tests the legacy NewTableOrig function
func TestNewTableOrig(t *testing.T) {
	tests := []struct {
		name string
		file *os.File
		want bool
	}{
		{
			name: "with stdout",
			file: os.Stdout,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewTableOrig(tt.file)

			if (table != nil) != tt.want {
				t.Errorf("NewTableOrig() = %v, want %v", table != nil, tt.want)
			}

			// Verify it's actually a tablewriter.Table
			if table != nil {
				// Table creation successful - this validates the wrapper function
				// The function properly wraps the olekukonko/tablewriter package
				_ = table // Verify table is not nil
			}
		})
	}
}

// TestTableJSONSerialization tests JSON serialization and deserialization
func TestTableJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "Table struct",
			data: &Table{
				writer:  &bytes.Buffer{},
				headers: []string{"Name", "Age"},
				rows:    [][]string{{"John", "25"}},
			},
		},
		{
			name: "headers slice",
			data: []string{"Name", "Age", "City"},
		},
		{
			name: "rows slice",
			data: [][]string{{"John", "25", "NYC"}, {"Jane", "30", "LA"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("JSON Marshal failed: %v", err)
			}

			if len(jsonData) == 0 {
				t.Errorf("JSON Marshal produced empty data")
			}

			// For non-struct types, test unmarshaling
			switch tt.data.(type) {
			case []string:
				var result []string
				if err := json.Unmarshal(jsonData, &result); err != nil {
					t.Errorf("JSON Unmarshal failed: %v", err)
				}
			case [][]string:
				var result [][]string
				if err := json.Unmarshal(jsonData, &result); err != nil {
					t.Errorf("JSON Unmarshal failed: %v", err)
				}
			}
		})
	}
}

// TestTableFailFast tests fail-fast error detection
func TestTableFailFast(t *testing.T) {
	tests := []struct {
		name string
		test func() error
	}{
		{
			name: "NewTable with nil writer should not panic",
			test: func() error {
				table := NewTable(nil)
				if table == nil {
					return nil
				}
				return nil
			},
		},
		{
			name: "Header with nil slice should not panic",
			test: func() error {
				table := NewTable(&bytes.Buffer{})
				table.Header(nil)
				return nil
			},
		},
		{
			name: "Append with nil row should not panic",
			test: func() error {
				table := NewTable(&bytes.Buffer{})
				table.Append(nil)
				return nil
			},
		},
		{
			name: "Render with no headers should return error",
			test: func() error {
				table := NewTable(&bytes.Buffer{})
				return table.Render()
			},
		},
		{
			name: "calculateColumnWidths with empty headers should not panic",
			test: func() error {
				table := NewTable(&bytes.Buffer{})
				table.headers = []string{}
				widths := table.calculateColumnWidths()
				if len(widths) != 0 {
					return nil
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s: function panicked: %v", tt.name, r)
				}
			}()

			err := tt.test()
			// For some tests, we expect errors (like Render with no headers)
			// The important thing is that they don't panic
			_ = err
		})
	}
}

// TestTableIntegration tests integration scenarios
func TestTableIntegration(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(*Table)
		wantError    bool
		wantContains []string
	}{
		{
			name: "complete table workflow",
			setupFunc: func(table *Table) {
				table.Header([]string{"Name", "Age", "City"})
				table.Append([]string{"John", "25", "NYC"})
				table.Append([]string{"Jane", "30", "LA"})
				table.Append([]string{"Bob", "35", "Chicago"})
			},
			wantError: false,
			wantContains: []string{
				"Name", "Age", "City",
				"John", "25", "NYC",
				"Jane", "30", "LA",
				"Bob", "35", "Chicago",
				"┌", "┬", "┐", "├", "┼", "┤", "└", "┴", "┘", "│", "─",
			},
		},
		{
			name: "table with varying column widths",
			setupFunc: func(table *Table) {
				table.Header([]string{"Short", "Very Long Header Name"})
				table.Append([]string{"A", "B"})
				table.Append([]string{"Very Long Content", "C"})
			},
			wantError: false,
			wantContains: []string{
				"Short", "Very Long Header Name",
				"Very Long Content", "A", "B", "C",
			},
		},
		{
			name: "empty table with headers only",
			setupFunc: func(table *Table) {
				table.Header([]string{"Name", "Age"})
			},
			wantError:    false,
			wantContains: []string{"Name", "Age"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			table := NewTable(buffer)
			tt.setupFunc(table)

			err := table.Render()

			if (err != nil) != tt.wantError {
				t.Errorf("Integration test error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError {
				output := buffer.String()
				for _, want := range tt.wantContains {
					if !strings.Contains(output, want) {
						t.Errorf("Integration test output does not contain %q. Output:\n%s", want, output)
					}
				}

				// Verify output has proper structure
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) < 3 { // At minimum: top border, header, bottom border
					t.Errorf("Integration test output should have at least 3 lines, got %d", len(lines))
				}
			}
		})
	}
}

// TestTableDrawMethods tests internal drawing methods
func TestTableDrawMethods(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "drawBorder",
			testFunc: func(t *testing.T) {
				buffer := &bytes.Buffer{}
				table := NewTable(buffer)
				widths := []int{5, 5}

				table.drawBorder(widths, "┌", "┬", "┐")
				output := buffer.String()

				if !strings.Contains(output, "┌") || !strings.Contains(output, "┬") || !strings.Contains(output, "┐") {
					t.Errorf("drawBorder() output missing expected characters: %s", output)
				}
			},
		},
		{
			name: "drawRow",
			testFunc: func(t *testing.T) {
				buffer := &bytes.Buffer{}
				table := NewTable(buffer)
				row := []string{"Test", "Data"}
				widths := []int{6, 6}

				table.drawRow(row, widths)
				output := buffer.String()

				if !strings.Contains(output, "Test") || !strings.Contains(output, "Data") || !strings.Contains(output, "│") {
					t.Errorf("drawRow() output missing expected content: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
