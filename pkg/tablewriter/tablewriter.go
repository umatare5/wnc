// Package tablewriter provides a simple wrapper around the tablewriter library
package tablewriter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
	writer  io.Writer
	headers []string
	rows    [][]string
}

func NewTable(writer io.Writer) *Table {
	if writer == nil {
		writer = os.Stdout
	}
	return &Table{writer: writer}
}

func (t *Table) Header(headers []string) {
	t.headers = headers
}

func (t *Table) Append(row []string) {
	t.rows = append(t.rows, row)
}

func (t *Table) Render() error {
	if len(t.headers) == 0 {
		return fmt.Errorf("no headers set")
	}

	// Calculate column widths
	widths := t.calculateColumnWidths()

	// Render table
	t.drawBorder(widths, "┌", "┬", "┐")
	t.drawRow(t.headers, widths)
	t.drawBorder(widths, "├", "┼", "┤")

	for _, row := range t.rows {
		t.drawRow(row, widths)
	}

	t.drawBorder(widths, "└", "┴", "┘")
	return nil
}

func (t *Table) calculateColumnWidths() []int {
	widths := make([]int, len(t.headers))

	// Get header widths
	for i, header := range t.headers {
		widths[i] = utf8.RuneCountInString(header)
	}

	// Get row widths
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(widths) {
				if cellWidth := utf8.RuneCountInString(cell); cellWidth > widths[i] {
					widths[i] = cellWidth
				}
			}
		}
	}

	// Add padding
	for i := range widths {
		widths[i] += 2 // 1 space on each side
	}

	return widths
}

func (t *Table) drawBorder(widths []int, left, middle, right string) {
	var border strings.Builder
	border.WriteString(left)
	for i, width := range widths {
		border.WriteString(strings.Repeat("─", width))
		if i < len(widths)-1 {
			border.WriteString(middle)
		}
	}
	border.WriteString(right)
	_, _ = fmt.Fprintln(t.writer, border.String())
}

func (t *Table) drawRow(row []string, widths []int) {
	var rowBuilder strings.Builder
	rowBuilder.WriteString("│")

	for i, width := range widths {
		var cell string
		if i < len(row) {
			cell = row[i]
		}

		// Calculate padding
		cellWidth := utf8.RuneCountInString(cell)
		padding := width - cellWidth
		leftPad := 1
		rightPad := padding - leftPad

		rowBuilder.WriteString(strings.Repeat(" ", leftPad))
		rowBuilder.WriteString(cell)
		rowBuilder.WriteString(strings.Repeat(" ", rightPad))
		rowBuilder.WriteString("│")
	}

	_, _ = fmt.Fprintln(t.writer, rowBuilder.String())
}

// Legacy function for compatibility - use olekukonko/tablewriter for fallback
func NewTableOrig(stdout *os.File) *tablewriter.Table {
	return tablewriter.NewWriter(stdout)
}
