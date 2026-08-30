package render

import (
	"bytes"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// Table writes the rows as a borderless, space-aligned table for piping: a full border spends
// three columns per field on views that run to twenty fields. An empty row set still prints the
// heading line, which distinguishes an empty fleet from a read that produced nothing.
func Table[T any](w io.Writer, cols []Column[T], rows []T) error {
	return write(w, cols, rows, plainCell[T], borderless(),
		tw.Padding{Left: "", Right: "  ", Overwrite: true}, tw.AlignLeft)
}

// PrettyTable writes the rows as a light-ruled, bordered table with a glyph in the columns that
// declare one. It is what --pretty selects, for reading on a terminal rather than for piping: the
// rules cost three columns per field, and three of the five glyphs are two columns wide.
func PrettyTable[T any](w io.Writer, cols []Column[T], rows []T) error {
	return write(w, cols, rows, prettyCell[T], bordered(), tw.PaddingDefault, tw.AlignCenter)
}

func plainCell[T any](c Column[T]) func(T) string { return c.Cell }

func prettyCell[T any](c Column[T]) func(T) string {
	if c.Pretty != nil {
		return c.Pretty
	}

	return c.Cell
}

func write[T any](
	w io.Writer, cols []Column[T], rows []T, cell func(Column[T]) func(T) string,
	rendition tw.Rendition, padding tw.Padding, header tw.Align,
) error {
	// The table is built into a buffer so each line can be stripped of the padding
	// tablewriter leaves past the last column. Trailing blanks are invisible on a
	// terminal and turn up in every checksum, diff and copied cell downstream.
	var buf bytes.Buffer

	t := tablewriter.NewTable(&buf,
		tablewriter.WithRenderer(renderer.NewBlueprint(rendition)),
		tablewriter.WithHeaderAlignment(header),
		tablewriter.WithRowAlignment(tw.AlignLeft),
		tablewriter.WithHeaderAutoFormat(tw.Off),
		tablewriter.WithTrimSpace(tw.On),
		tablewriter.WithPadding(padding),
	)

	t.Header(Headers(cols))

	for _, row := range rows {
		cells := make([]string, 0, len(cols))
		for _, c := range cols {
			cells = append(cells, cell(c)(row))
		}

		if err := t.Append(cells); err != nil {
			return fmt.Errorf("appending a table row: %w", err)
		}
	}

	if err := t.Render(); err != nil {
		return fmt.Errorf("rendering the table: %w", err)
	}

	return writeTrimmed(w, buf.Bytes())
}

// writeTrimmed copies the rendered table out, one line at a time, without the
// trailing padding.
func writeTrimmed(w io.Writer, b []byte) error {
	for line := range bytes.SplitSeq(bytes.TrimRight(b, "\n"), []byte("\n")) {
		if _, err := w.Write(append(bytes.TrimRight(line, " \t"), '\n')); err != nil {
			return fmt.Errorf("writing a table line: %w", err)
		}
	}

	return nil
}

// bordered draws tw.StyleLight, so the box-drawing set is the library's and the column widths are
// measured by it rather than by counting runes.
func bordered() tw.Rendition {
	return tw.Rendition{
		Borders: tw.Border{Left: tw.On, Right: tw.On, Top: tw.On, Bottom: tw.On},
		Symbols: tw.NewSymbols(tw.StyleLight),
		Settings: tw.Settings{
			Separators: tw.Separators{
				ShowHeader:     tw.On,
				ShowFooter:     tw.Off,
				BetweenRows:    tw.Off,
				BetweenColumns: tw.On,
			},
			Lines: tw.Lines{
				ShowTop:        tw.On,
				ShowBottom:     tw.On,
				ShowHeaderLine: tw.On,
				ShowFooterLine: tw.Off,
			},
		},
	}
}

// borderless turns off every rule and separator, leaving two spaces of padding as
// the only column gap.
func borderless() tw.Rendition {
	return tw.Rendition{
		Borders: tw.Border{Left: tw.Off, Right: tw.Off, Top: tw.Off, Bottom: tw.Off},
		Settings: tw.Settings{
			Separators: tw.Separators{
				ShowHeader:     tw.Off,
				ShowFooter:     tw.Off,
				BetweenRows:    tw.Off,
				BetweenColumns: tw.Off,
			},
			Lines: tw.Lines{
				ShowTop:        tw.Off,
				ShowBottom:     tw.Off,
				ShowHeaderLine: tw.Off,
				ShowFooterLine: tw.Off,
			},
		},
	}
}
