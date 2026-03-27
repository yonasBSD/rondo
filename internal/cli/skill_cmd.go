package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func (c *CLI) skillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage Claude Code skill integration",
	}

	cmd.AddCommand(c.skillInstallCmd())
	cmd.AddCommand(c.skillUninstallCmd())
	return cmd
}

func (c *CLI) skillInstallCmd() *cobra.Command {
	var project bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install rondo skill for Claude Code",
		Long: `Install the rondo skill so Claude Code can manage tasks, journal entries,
subtasks, time logs, and focus sessions.

By default installs to ~/.claude/skills/rondo/ (available in all projects).
Use --project to install to ./.claude/skills/rondo/ (current project only).`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := skillDir(project)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("create directory %s: %w", dir, err)
			}

			path := filepath.Join(dir, "SKILL.md")
			if err := os.WriteFile(path, []byte(skillContent), 0o644); err != nil {
				return fmt.Errorf("write skill file: %w", err)
			}

			p := c.printer(os.Stdout)
			p.Success("Skill installed at %s", path)
			return nil
		},
	}

	cmd.Flags().BoolVar(&project, "project", false, "Install to current project (.claude/skills/) instead of global (~/.claude/skills/)")
	return cmd
}

func (c *CLI) skillUninstallCmd() *cobra.Command {
	var project bool

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove rondo skill from Claude Code",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := skillDir(project)
			if err != nil {
				return err
			}

			if _, err := os.Stat(dir); err != nil {
				if os.IsNotExist(err) {
					p := c.printer(os.Stdout)
					p.Success("Skill not installed, nothing to remove")
					return nil
				}
				return fmt.Errorf("check skill directory: %w", err)
			}

			if err := os.RemoveAll(dir); err != nil {
				return fmt.Errorf("remove skill directory: %w", err)
			}

			p := c.printer(os.Stdout)
			p.Success("Skill removed from %s", dir)
			return nil
		},
	}

	cmd.Flags().BoolVar(&project, "project", false, "Remove from current project (.claude/skills/) instead of global (~/.claude/skills/)")
	return cmd
}

// skillDir returns the target directory for the rondo skill.
func skillDir(project bool) (string, error) {
	if project {
		return filepath.Join(".claude", "skills", "rondo"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".claude", "skills", "rondo"), nil
}
