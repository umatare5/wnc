// Package render turns a slice of rows into a borderless table, a bordered one, or a flat
// JSON array. All three are built from the same column list, so a column cannot exist in
// one and be missing from the others.
package render

// Absent is what a table cell shows where the controller reported nothing. The
// JSON form omits the key instead of writing this, so a consumer parsing the JSON
// never has to recognize a placeholder.
const Absent = "-"

// Column is one output column: its JSON field name, its table heading and the cell
// it derives from a row.
type Column[T any] struct {
	// Key is the JSON field name and the value --sort-by accepts. An invariant test
	// asserts it matches the row struct's json tag for the same position.
	Key string

	// Header is the table heading.
	Header string

	// Cell renders the table cell. It is also what Sort falls back to, so this and
	// not Pretty is the text an unsorted column is ordered by.
	Cell func(T) string

	// Pretty renders the cell for the bordered table only, and is nil on every column
	// whose two renderings are the same. It never feeds Sort or the JSON, so a glyph
	// here cannot reorder rows or reach a consumer parsing the output.
	Pretty func(T) string

	// Sort yields the value to order by, and is nil where the rendered text already orders
	// correctly: an octet count printed as "1.0KiB" and an age printed as "3d4h" do not.
	Sort func(T) any
}

func Keys[T any](cols []Column[T]) []string {
	out := make([]string, 0, len(cols))
	for _, c := range cols {
		out = append(out, c.Key)
	}

	return out
}

// SortValue converts an optional integer into the float64 a comparator takes, or nil for an
// unreported cell. A column's Sort yields one of nil, a string, a float64 or a bool.
func SortValue[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](p *T) any {
	if p == nil {
		return nil
	}

	return float64(*p)
}

func Headers[T any](cols []Column[T]) []string {
	out := make([]string, 0, len(cols))
	for _, c := range cols {
		out = append(out, c.Header)
	}

	return out
}
