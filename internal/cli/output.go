package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/term"
	"github.com/mattn/go-isatty"
	"github.com/roniel/todo-app/internal/ui"
)

// Printer handles formatted output to a writer.
type Printer struct {
	format  string
	quiet   bool
	noColor bool
	w       io.Writer
}

// newPrinter creates a Printer writing to w.
func newPrinter(w io.Writer, format string, quiet, noColor bool) *Printer {
	return &Printer{format: format, quiet: quiet, noColor: noColor, w: w}
}

// isTTY reports whether the given file is connected to a terminal.
func isTTY(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

// Success prints a success message unless --quiet is set.
// When color is enabled, it prepends a green checkmark.
func (p *Printer) Success(format string, args ...any) {
	if p.quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if p.noColor {
		fmt.Fprintln(p.w, msg)
	} else {
		prefix := lipgloss.NewStyle().Foreground(ui.Green).Render("✓")
		fmt.Fprintf(p.w, "%s %s\n", prefix, msg)
	}
}

// Bold returns s in bold when color is enabled.
func (p *Printer) Bold(s string) string {
	if p.noColor {
		return s
	}
	return lipgloss.NewStyle().Bold(true).Render(s)
}

// Dim returns s in gray when color is enabled.
func (p *Printer) Dim(s string) string {
	if p.noColor {
		return s
	}
	return lipgloss.NewStyle().Foreground(ui.Gray).Render(s)
}

// Table prints tabular output. When color is enabled it renders a styled
// table with rounded Unicode borders; otherwise it falls back to plain
// tab-aligned text for piped output.
func (p *Printer) Table(headers []string, rows [][]string) {
	if p.noColor {
		tw := tabwriter.NewWriter(p.w, 0, 4, 2, ' ', 0)
		for i, h := range headers {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			fmt.Fprint(tw, h)
		}
		fmt.Fprintln(tw)
		for i, h := range headers {
			if i > 0 {
				fmt.Fprint(tw, "\t")
			}
			for j := 0; j < len(h); j++ {
				fmt.Fprint(tw, "-")
			}
		}
		fmt.Fprintln(tw)
		for _, row := range rows {
			for i, cell := range row {
				if i > 0 {
					fmt.Fprint(tw, "\t")
				}
				fmt.Fprint(tw, cell)
			}
			fmt.Fprintln(tw)
		}
		tw.Flush()
		return
	}

	t := table.New().
		Headers(headers...).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(ui.Cyan)).
		Width(writerWidth(p.w)).
		Wrap(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Padding(0, 1)
			}
			return lipgloss.NewStyle().Padding(0, 1)
		})

	for _, row := range rows {
		t.Row(row...)
	}

	fmt.Fprintln(p.w, t.Render())
}

// writerWidth returns the terminal width if w is a terminal file,
// or 80 as a fallback.
func writerWidth(w io.Writer) int {
	if f, ok := w.(*os.File); ok {
		width, _, err := term.GetSize(f.Fd())
		if err == nil && width > 0 {
			return width
		}
	}
	return 80
}

// Colored applies the given color to s when color output is enabled.
func (p *Printer) Colored(s string, color lipgloss.TerminalColor) string {
	if p.noColor {
		return s
	}
	return lipgloss.NewStyle().Foreground(color).Render(s)
}

// JSON prints v as indented JSON to the printer's writer.
func (p *Printer) JSON(v any) error {
	enc := json.NewEncoder(p.w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
