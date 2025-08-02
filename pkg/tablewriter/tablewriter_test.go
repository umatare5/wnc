package tablewriter

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestTableAppend tests the Append method
func TestTableAppend(t *testing.T) {
	var buf bytes.Buffer
	table := NewTable(&buf)

	row := []string{"John", "25", "NYC"}
	table.Append(row)

	if len(table.rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(table.rows))
	}
}

// TestTableHeader tests the Header method
func TestTableHeader(t *testing.T) {
	var buf bytes.Buffer
	table := NewTable(&buf)

	headers := []string{"Name", "Age", "City"}
	table.Header(headers)

	if len(table.headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(table.headers))
	}
}

// TestTableRender tests the Render method
func TestTableRender(t *testing.T) {
	var buf bytes.Buffer
	table := NewTable(&buf)

	table.Header([]string{"Name", "Age"})
	table.Append([]string{"John", "25"})

	err := table.Render()
	if err != nil {
		t.Errorf("Render() returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Name") {
		t.Error("Output should contain header 'Name'")
	}
}

// TestNewTableCreation tests the NewTable function
func TestNewTableCreation(t *testing.T) {
	var buf bytes.Buffer
	table := NewTable(&buf)

	if table == nil {
		t.Error("NewTable should return non-nil table")
	}
}

// TestNewTableOrig tests the NewTableOrig function
func TestNewTableOrig(t *testing.T) {
	t.Run("test_new_table_orig_function", func(t *testing.T) {
		// Test that NewTableOrig creates a table writer
		table := NewTableOrig(os.Stdout)

		if table == nil {
			t.Error("NewTableOrig should return non-nil table")
		}
	})

	t.Run("test_new_table_orig_with_different_output", func(t *testing.T) {
		// Create a temporary file for testing
		tmpFile, err := os.CreateTemp("", "test_table_*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		// Test NewTableOrig with file output
		table := NewTableOrig(tmpFile)

		if table == nil {
			t.Error("NewTableOrig should return non-nil table")
		}
	})
}
