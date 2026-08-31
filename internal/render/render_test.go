package render

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/olekukonko/tablewriter/pkg/twwidth"
)

type row struct {
	Name  string  `json:"name"`
	Count *int    `json:"count,omitzero"`
	Note  *string `json:"note,omitzero"`
}

func cols() []Column[row] {
	return []Column[row]{
		{Key: "name", Header: "Name", Cell: func(r row) string { return Str(r.Name) }},
		{
			Key: "count", Header: "Count",
			Cell: func(r row) string { return IntPtr(r.Count) },
			Sort: func(r row) any { return SortValue(r.Count) },
		},
		{Key: "note", Header: "Note", Cell: func(r row) string { return StrPtr(r.Note) }},
	}
}

func TestTableIsBorderless(t *testing.T) {
	t.Parallel()

	zero, note := 0, "ok"
	rows := []row{
		{Name: "TEST-AP01", Count: &zero, Note: &note},
		{Name: "TEST-AP02"},
	}

	var buf bytes.Buffer
	if err := Table(&buf, cols(), rows); err != nil {
		t.Fatalf("Table: %v", err)
	}

	got := buf.String()

	for _, ch := range []string{"|", "+", "-+-", "─", "│"} {
		if strings.Contains(got, ch) && ch != "-+-" {
			t.Errorf("output holds the rule character %q:\n%s", ch, got)
		}
	}

	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected a heading and two rows, got %d lines:\n%s", len(lines), got)
	}

	if !strings.HasPrefix(strings.TrimSpace(lines[0]), "Name") {
		t.Errorf("heading line = %q", lines[0])
	}

	// A reported zero must survive as 0, and an unreported count must not.
	if !strings.Contains(lines[1], "0") {
		t.Errorf("a reported zero was lost: %q", lines[1])
	}

	if !strings.Contains(lines[2], Absent) {
		t.Errorf("an unreported value did not render %q: %q", Absent, lines[2])
	}
}

// An empty fleet is a normal answer, so the heading still prints: that is what
// separates it from a read that failed.
func TestTableEmptyKeepsTheHeading(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := Table(&buf, cols(), []row{}); err != nil {
		t.Fatalf("Table: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(got, "Name") {
		t.Errorf("output = %q, want the heading alone", got)
	}

	if strings.Count(got, "\n") != 0 {
		t.Errorf("output has more than the heading line: %q", got)
	}
}

func TestJSONShape(t *testing.T) {
	t.Parallel()

	zero, note := 0, "ok"

	var buf bytes.Buffer
	if err := JSON(&buf, []row{{Name: "a", Count: &zero, Note: &note}, {Name: "b"}}); err != nil {
		t.Fatalf("JSON: %v", err)
	}

	want := `[{"name":"a","count":0,"note":"ok"},{"name":"b"}]` + "\n"
	if got := buf.String(); got != want {
		t.Errorf("JSON =\n%s\nwant\n%s", got, want)
	}
}

// json/v2 renders a nil slice as [] where v1 wrote null, so an empty result needs
// no normalizing step of its own.
func TestJSONEmptyIsAnArray(t *testing.T) {
	t.Parallel()

	for name, rows := range map[string][]row{"nil": nil, "empty": {}} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := JSON(&buf, rows); err != nil {
				t.Fatalf("JSON: %v", err)
			}

			if got := buf.String(); got != "[]\n" {
				t.Errorf("JSON = %q, want %q", got, "[]\n")
			}
		})
	}
}

// A controller may name an SSID with a character v1 would have escaped. v2 does
// not, so the value appears as the controller sent it.
func TestJSONDoesNotEscapeHTML(t *testing.T) {
	t.Parallel()

	name := "a&b<c>"

	var buf bytes.Buffer
	if err := JSON(&buf, []row{{Name: name}}); err != nil {
		t.Fatalf("JSON: %v", err)
	}

	if !strings.Contains(buf.String(), name) {
		t.Errorf("JSON = %s, want the raw %q", buf.String(), name)
	}
}

func TestIEC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *uint64
		want string
	}{
		{name: "absent", in: nil, want: Absent},
		{name: "zero", in: u64(0), want: "0B"},
		{name: "below a kibibyte", in: u64(1023), want: "1023B"},
		{name: "one kibibyte", in: u64(1024), want: "1.0KiB"},
		{name: "mebibytes", in: u64(5 * 1024 * 1024), want: "5.0MiB"},
		{name: "gibibytes", in: u64(3*1024*1024*1024 + 512*1024*1024), want: "3.5GiB"},
		{name: "exbibytes", in: u64(2 << 60), want: "2.0EiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := IEC(tt.in); got != tt.want {
				t.Errorf("IEC(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   *int64
		want string
	}{
		{name: "absent", in: nil, want: Absent},
		{name: "negative", in: i64(-1), want: Absent},
		{name: "zero", in: i64(0), want: "0s"},
		{name: "seconds", in: i64(45), want: "45s"},
		{name: "minutes", in: i64(200), want: "3m20s"},
		{name: "whole minutes", in: i64(180), want: "3m"},
		{name: "hours", in: i64(3725), want: "1h2m"},
		{name: "days", in: i64(97200), want: "1d3h"},
		{name: "weeks", in: i64(1209600), want: "14d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Duration(tt.in); got != tt.want {
				t.Errorf("Duration(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	t.Parallel()

	yes, no := true, false

	if got := Bool(nil); got != Absent {
		t.Errorf("Bool(nil) = %q, want %q", got, Absent)
	}

	if got := Bool(&yes); got != "Yes" {
		t.Errorf("Bool(true) = %q", got)
	}

	if got := Bool(&no); got != "No" {
		t.Errorf("Bool(false) = %q", got)
	}
}

// The controller returns the Unix epoch on several timestamp leaves to mean "no
// instant", so an age must not be computed from one.
func TestSecondsSinceRejectsSentinels(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)

	if got := SecondsSince(now, time.Time{}); got != nil {
		t.Errorf("zero time gave %v, want nil", *got)
	}

	if got := SecondsSince(now, time.Unix(0, 0)); got != nil {
		t.Errorf("epoch gave %v, want nil", *got)
	}

	if got := SecondsSince(now, now.Add(time.Hour)); got != nil {
		t.Errorf("a future instant gave %v, want nil", *got)
	}

	got := SecondsSince(now, now.Add(-90*time.Second))
	if got == nil || *got != 90 {
		t.Errorf("SecondsSince = %v, want 90", got)
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	if got := Join(nil, "/"); got != Absent {
		t.Errorf("Join(nil) = %q, want %q", got, Absent)
	}

	if got := Join([]string{"5", "6"}, "/"); got != "5/6" {
		t.Errorf("Join = %q", got)
	}
}

func TestSortByString(t *testing.T) {
	t.Parallel()

	rows := []row{{Name: "c"}, {Name: "a"}, {Name: "b"}}
	if err := Sort(rows, cols(), "name", false); err != nil {
		t.Fatalf("Sort: %v", err)
	}

	if got := names(rows); got != "a,b,c" {
		t.Errorf("ascending = %s", got)
	}

	if err := Sort(rows, cols(), "name", true); err != nil {
		t.Fatalf("Sort: %v", err)
	}

	if got := names(rows); got != "c,b,a" {
		t.Errorf("descending = %s", got)
	}
}

// A numeric column must not be ordered by its rendered text, or 9 would follow 10.
func TestSortByNumberNotText(t *testing.T) {
	t.Parallel()

	n2, n9, n10 := 2, 9, 10
	rows := []row{{Name: "ten", Count: &n10}, {Name: "nine", Count: &n9}, {Name: "two", Count: &n2}}

	if err := Sort(rows, cols(), "count", false); err != nil {
		t.Fatalf("Sort: %v", err)
	}

	if got := names(rows); got != "two,nine,ten" {
		t.Errorf("ascending = %s, want two,nine,ten", got)
	}
}

// Reversing the order is meant to bring the other end of the values forward, not to
// promote the rows that carry no value.
func TestSortKeepsAbsentLastInBothDirections(t *testing.T) {
	t.Parallel()

	n1, n2 := 1, 2

	for _, desc := range []bool{false, true} {
		t.Run(map[bool]string{false: "asc", true: "desc"}[desc], func(t *testing.T) {
			t.Parallel()

			rows := []row{{Name: "none"}, {Name: "one", Count: &n1}, {Name: "two", Count: &n2}}
			if err := Sort(rows, cols(), "count", desc); err != nil {
				t.Fatalf("Sort: %v", err)
			}

			if rows[len(rows)-1].Name != "none" {
				t.Errorf("order = %s, want the unreported row last", names(rows))
			}
		})
	}
}

// A string column with no Sort function falls back to the cell, where Absent must
// still be treated as absence rather than as the literal "-".
func TestSortTreatsAbsentCellAsAbsence(t *testing.T) {
	t.Parallel()

	zzz := "zzz"
	rows := []row{{Name: "x"}, {Name: "y", Note: &zzz}}

	if err := Sort(rows, cols(), "note", false); err != nil {
		t.Fatalf("Sort: %v", err)
	}

	if rows[0].Name != "y" {
		t.Errorf("order = %s, want the reported note first", names(rows))
	}
}

func TestSortRejectsUnknownKey(t *testing.T) {
	t.Parallel()

	if err := Sort([]row{}, cols(), "bogus", false); err == nil {
		t.Error("Sort accepted an unknown column")
	}
}

func TestKeysMatchColumnOrder(t *testing.T) {
	t.Parallel()

	if got := strings.Join(Keys(cols()), ","); got != "name,count,note" {
		t.Errorf("Keys = %s", got)
	}
}

func names(rows []row) string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.Name)
	}

	return strings.Join(out, ",")
}

func u64(v uint64) *uint64 { return &v }
func i64(v int64) *int64   { return &v }

// The bordered table is a style, not a second data path: the same rows must reach it,
// the rules must line up whatever a cell holds, and the default must not move.
func TestPrettyTable(t *testing.T) {
	t.Parallel()

	type row struct {
		Name  string
		Wide  string
		State *string
	}

	cols := []Column[row]{
		{Key: "name", Header: "Name", Cell: func(r row) string { return Str(r.Name) }},
		{Key: "wide", Header: "Site", Cell: func(r row) string { return Str(r.Wide) }},
		{
			Key: "state", Header: "State",
			Cell:   func(r row) string { return StrPtr(r.State) },
			Pretty: func(r row) string { return glyphFor(r.State) },
		},
	}

	rows := []row{
		{Name: "TEST-AP01", Wide: "Tokyo", State: strPtr("Up")},
		{Name: "TEST-AP02", Wide: "Osaka", State: strPtr("Down")},
		{Name: "TEST-AP03", Wide: "", State: nil},
	}

	var pretty, plain bytes.Buffer
	if err := PrettyTable(&pretty, cols, rows); err != nil {
		t.Fatalf("PrettyTable: %v", err)
	}

	if err := Table(&plain, cols, rows); err != nil {
		t.Fatalf("Table: %v", err)
	}

	// Every line of the bordered table is the same display width, or a rule falls in
	// the wrong column on the reader's terminal.
	lines := strings.Split(strings.TrimRight(pretty.String(), "\n"), "\n")
	if len(lines) != 7 {
		t.Fatalf("got %d lines, want 7 (three rules, a heading and three rows)", len(lines))
	}

	widths := map[int]bool{}
	for _, l := range lines {
		widths[twwidth.Width(l)] = true
	}

	if len(widths) != 1 {
		t.Errorf("the rules do not line up: widths %v\n%s", widths, pretty.String())
	}

	if !strings.HasPrefix(lines[0], "┌") || !strings.HasPrefix(lines[len(lines)-1], "└") {
		t.Errorf("the table is not bordered:\n%s", pretty.String())
	}

	// The glyph replaces the cell in the bordered table and nowhere else.
	if !strings.Contains(pretty.String(), "✅") {
		t.Error("the bordered table carries no glyph")
	}

	if strings.ContainsAny(plain.String(), "✅✕│") {
		t.Errorf("the plain table gained a glyph or a rule:\n%s", plain.String())
	}

	// An absent cell keeps the marker under both renderings.
	for name, out := range map[string]string{"pretty": pretty.String(), "plain": plain.String()} {
		if !strings.Contains(out, Absent) {
			t.Errorf("the %s table lost the absence marker:\n%s", name, out)
		}
	}
}

// A column with no Pretty function renders identically in both tables.
func TestPrettyTableFallsBackToCell(t *testing.T) {
	t.Parallel()

	type row struct{ Name string }

	cols := []Column[row]{{Key: "name", Header: "Name", Cell: func(r row) string { return Str(r.Name) }}}
	rows := []row{{Name: "TEST-AP01"}}

	var buf bytes.Buffer
	if err := PrettyTable(&buf, cols, rows); err != nil {
		t.Fatalf("PrettyTable: %v", err)
	}

	if !strings.Contains(buf.String(), "TEST-AP01") {
		t.Errorf("a column with no Pretty rendering lost its cell:\n%s", buf.String())
	}
}

// glyphFor deliberately mixes widths. show ap's Admin column holds a two-column check
// mark and a one-column cross in different rows, so the padding has to be measured per
// cell rather than assumed from the first one.
func glyphFor(p *string) string {
	switch {
	case p == nil:
		return Absent
	case *p == "Up":
		return "✅"
	default:
		return "✕"
	}
}

func strPtr(s string) *string { return &s }
