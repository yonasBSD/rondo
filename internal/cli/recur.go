package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/roniel/todo-app/internal/task"
	"github.com/spf13/cobra"
)

func (c *CLI) recurCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recur",
		Short: "Manage task recurrence",
	}
	cmd.AddCommand(
		c.recurSetCmd(),
		c.recurClearCmd(),
	)
	return cmd
}

func (c *CLI) recurSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <task-id> <daily|weekly|monthly|yearly>",
		Short: "Set recurrence for a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task ID %q: %w", args[0], err)
			}
			t, err := c.getTaskOrNotFound(taskID)
			if err != nil {
				return err
			}
			freq := task.ParseRecurFreq(strings.ToLower(args[1]))
			if freq == task.RecurNone {
				return fmt.Errorf("invalid frequency %q: must be daily, weekly, monthly, or yearly", args[1])
			}
			interval := t.RecurInterval
			if interval <= 0 {
				interval = 1
			}
			if err := c.taskStore.UpdateRecurrence(taskID, freq, interval); err != nil {
				return fmt.Errorf("set recurrence: %w", err)
			}
			c.printer(os.Stdout).Success("Set task #%d to recur %s", taskID, freq.String())
			return nil
		},
	}
}

func (c *CLI) recurClearCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clear <task-id>",
		Short: "Clear recurrence for a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task ID %q: %w", args[0], err)
			}
			if _, err := c.getTaskOrNotFound(taskID); err != nil {
				return err
			}
			if err := c.taskStore.UpdateRecurrence(taskID, task.RecurNone, 0); err != nil {
				return fmt.Errorf("clear recurrence: %w", err)
			}
			c.printer(os.Stdout).Success("Cleared recurrence for task #%d", taskID)
			return nil
		},
	}
}
