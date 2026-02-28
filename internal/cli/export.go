package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/roniel/todo-app/internal/export"
	"github.com/roniel/todo-app/internal/journal"
	"github.com/spf13/cobra"
)

func (c *CLI) exportCmd() *cobra.Command {
	var format, output string
	var includeJournal bool

	cmd := &cobra.Command{
		Use:   "export [flags]",
		Short: "Export tasks and optionally journal to a file or stdout",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tasks, err := c.taskStore.List()
			if err != nil {
				return fmt.Errorf("list tasks: %w", err)
			}

			var notes []journal.Note
			if includeJournal {
				notes, err = c.journalStore.ListNotes(false)
				if err != nil {
					return fmt.Errorf("list journal notes: %w", err)
				}
			}

			var w io.Writer = os.Stdout
			var bw *bufio.Writer
			if output != "" {
				f, err := os.Create(output)
				if err != nil {
					return fmt.Errorf("create output file: %w", err)
				}
				defer f.Close()
				bw = bufio.NewWriter(f)
				w = bw
			}

			switch format {
			case "md", "markdown":
				if err := export.WriteTasks(w, tasks); err != nil {
					return fmt.Errorf("write tasks: %w", err)
				}
				if includeJournal {
					if _, err := fmt.Fprintln(w); err != nil {
						return err
					}
					if err := export.WriteNotes(w, notes); err != nil {
						return fmt.Errorf("write journal: %w", err)
					}
				}
			case "json":
				if !includeJournal {
					notes = nil
				}
				if err := export.WriteJSON(w, tasks, notes); err != nil {
					return fmt.Errorf("write json: %w", err)
				}
			default:
				return fmt.Errorf("invalid format %q: must be md or json", format)
			}

			if bw != nil {
				if err := bw.Flush(); err != nil {
					return fmt.Errorf("flush output: %w", err)
				}
			}

			if output != "" {
				fmt.Fprintf(os.Stderr, "Exported to %s\n", output)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "md", "Export format: md, json")
	cmd.Flags().StringVar(&output, "output", "", "Output file path (default: stdout)")
	cmd.Flags().BoolVar(&includeJournal, "journal", false, "Include journal entries")

	return cmd
}
