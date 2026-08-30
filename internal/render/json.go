package render

import (
	"encoding/json/v2"
	"fmt"
	"io"
)

// JSON writes the rows as a flat array with the field names --sort-by accepts. json/v2 marshals a
// nil slice to [] and escapes no HTML, and MarshalWrite adds no trailing newline, so one is
// appended here.
func JSON[T any](w io.Writer, rows []T) error {
	if err := json.MarshalWrite(w, rows); err != nil {
		return fmt.Errorf("encoding rows as JSON: %w", err)
	}

	if _, err := io.WriteString(w, "\n"); err != nil {
		return fmt.Errorf("writing the JSON terminator: %w", err)
	}

	return nil
}
