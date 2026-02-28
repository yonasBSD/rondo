package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/roniel/todo-app/internal/task"
	"github.com/roniel/todo-app/internal/ui"
	"github.com/spf13/cobra"
)

func (c *CLI) subtaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subtask",
		Short: "Manage subtasks",
	}
	cmd.AddCommand(
		c.subtaskAddCmd(),
		c.subtaskListCmd(),
		c.subtaskDoneCmd(),
		c.subtaskEditCmd(),
		c.subtaskDeleteCmd(),
	)
	return cmd
}

func (c *CLI) subtaskAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <task-id> \"title\"",
		Short: "Add a subtask to a task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task ID %q: %w", args[0], err)
			}
			if _, err := c.getTaskOrNotFound(taskID); err != nil {
				return err
			}
			title := args[1]
			if err := c.taskStore.AddSubtask(taskID, title); err != nil {
				return fmt.Errorf("add subtask: %w", err)
			}
			c.printer(os.Stdout).Success("Added subtask to task #%d: %s", taskID, title)
			return nil
		},
	}
}

func (c *CLI) subtaskListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <task-id>",
		Short: "List subtasks for a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task ID %q: %w", args[0], err)
			}
			t, err := c.getTaskOrNotFound(taskID)
			if err != nil {
				return err
			}
			p := c.printer(os.Stdout)
			if strings.ToLower(c.format) == "json" {
				type jsonSub struct {
					ID       int64  `json:"id"`
					Done     bool   `json:"done"`
					Title    string `json:"title"`
					Position int    `json:"position"`
				}
				out := make([]jsonSub, 0, len(t.Subtasks))
				for _, s := range t.Subtasks {
					out = append(out, jsonSub{
						ID:       s.ID,
						Done:     s.Completed,
						Title:    s.Title,
						Position: s.Position,
					})
				}
				return p.JSON(out)
			}
			rows := make([][]string, 0, len(t.Subtasks))
			for _, s := range t.Subtasks {
				done := p.Colored("○", ui.Gray)
				if s.Completed {
					done = p.Colored("✓", ui.Green)
				}
				rows = append(rows, []string{
					strconv.FormatInt(s.ID, 10),
					done,
					s.Title,
				})
			}
			p.Table([]string{"ID", "DONE", "TITLE"}, rows)
			return nil
		},
	}
}

func (c *CLI) subtaskDoneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "done <task-id> <subtask-id>",
		Short: "Toggle subtask completion",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task ID %q: %w", args[0], err)
			}
			subtaskID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid subtask ID %q: %w", args[1], err)
			}
			t, err := c.getTaskOrNotFound(taskID)
			if err != nil {
				return err
			}
			if !subtaskBelongsToTask(t, subtaskID) {
				return &NotFoundError{Type: "subtask", ID: subtaskID}
			}
			if err := c.taskStore.ToggleSubtask(subtaskID); err != nil {
				return fmt.Errorf("toggle subtask: %w", err)
			}
			c.printer(os.Stdout).Success("Toggled subtask #%d", subtaskID)
			return nil
		},
	}
}

func (c *CLI) subtaskEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit <task-id> <subtask-id> \"new title\"",
		Short: "Edit a subtask title",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task ID %q: %w", args[0], err)
			}
			subtaskID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid subtask ID %q: %w", args[1], err)
			}
			newTitle := args[2]
			t, err := c.getTaskOrNotFound(taskID)
			if err != nil {
				return err
			}
			if !subtaskBelongsToTask(t, subtaskID) {
				return &NotFoundError{Type: "subtask", ID: subtaskID}
			}
			if err := c.taskStore.UpdateSubtask(subtaskID, newTitle); err != nil {
				return fmt.Errorf("update subtask: %w", err)
			}
			c.printer(os.Stdout).Success("Updated subtask #%d: %s", subtaskID, newTitle)
			return nil
		},
	}
}

func (c *CLI) subtaskDeleteCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "delete <task-id> <subtask-id>",
		Short: "Delete a subtask",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid task ID %q: %w", args[0], err)
			}
			subtaskID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid subtask ID %q: %w", args[1], err)
			}
			t, err := c.getTaskOrNotFound(taskID)
			if err != nil {
				return err
			}
			var subtaskTitle string
			found := false
			for _, s := range t.Subtasks {
				if s.ID == subtaskID {
					subtaskTitle = s.Title
					found = true
					break
				}
			}
			if !found {
				return &NotFoundError{Type: "subtask", ID: subtaskID}
			}
			ok, err := Confirm(fmt.Sprintf("Delete subtask #%d %q?", subtaskID, subtaskTitle), force)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Fprintln(os.Stderr, "Cancelled.")
				return nil
			}
			if err := c.taskStore.DeleteSubtask(subtaskID); err != nil {
				return fmt.Errorf("delete subtask: %w", err)
			}
			c.printer(os.Stdout).Success("Deleted subtask #%d", subtaskID)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "y", false, "Skip confirmation prompt")
	return cmd
}

// subtaskBelongsToTask reports whether subtaskID is among the task's subtasks.
func subtaskBelongsToTask(t *task.Task, subtaskID int64) bool {
	for _, s := range t.Subtasks {
		if s.ID == subtaskID {
			return true
		}
	}
	return false
}
