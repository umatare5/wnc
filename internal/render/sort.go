package render

import (
	"cmp"
	"fmt"
	"slices"
)

// Sort orders the rows by one column, named by its JSON field name, on the JSON value rather than
// the rendered text. An unreported cell sorts after every reported one in both directions, because
// reversing the order must not promote the rows that have no value.
func Sort[T any](rows []T, cols []Column[T], key string, desc bool) error {
	col, ok := find(cols, key)
	if !ok {
		return fmt.Errorf("no such column %q", key)
	}

	value := col.Sort
	if value == nil {
		// The JSON value of a string or enum column is the cell text, so the rendered
		// cell orders it correctly. Absent stays absent.
		value = func(row T) any {
			if s := col.Cell(row); s != Absent {
				return s
			}

			return nil
		}
	}

	slices.SortStableFunc(rows, func(a, b T) int {
		return compare(value(a), value(b), desc)
	})

	return nil
}

func find[T any](cols []Column[T], key string) (Column[T], bool) {
	for _, c := range cols {
		if c.Key == key {
			return c, true
		}
	}

	return Column[T]{}, false
}

// compare orders two cell values, holding absence at the end regardless of
// direction. A pair of differing dynamic types cannot arise from one column, so it
// is reported as equal rather than guessed at.
func compare(a, b any, desc bool) int {
	switch {
	case a == nil && b == nil:
		return 0
	case a == nil:
		return 1
	case b == nil:
		return -1
	}

	n := compareSame(a, b)
	if desc {
		return -n
	}

	return n
}

// compareSame orders two present values of the same shape.
func compareSame(a, b any) int {
	switch x := a.(type) {
	case string:
		y, ok := b.(string)
		if !ok {
			return 0
		}

		return cmp.Compare(x, y)
	case float64:
		y, ok := b.(float64)
		if !ok {
			return 0
		}

		return cmp.Compare(x, y)
	case bool:
		y, ok := b.(bool)
		if !ok {
			return 0
		}

		return compareBool(x, y)
	default:
		return 0
	}
}

// compareBool puts false before true, matching how the rendered No and Yes order.
func compareBool(a, b bool) int {
	switch {
	case a == b:
		return 0
	case !a:
		return -1
	default:
		return 1
	}
}
