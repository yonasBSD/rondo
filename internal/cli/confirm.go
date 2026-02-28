package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/roniel/todo-app/internal/ui"
)

// Confirm prompts the user for confirmation and returns their answer.
// If force is true, it returns true without prompting.
// If stdin is not a TTY and force is false, it returns an error.
func Confirm(prompt string, force bool) (bool, error) {
	if force {
		return true, nil
	}
	if !isTTY(os.Stdin) {
		return false, fmt.Errorf("stdin is not a TTY: use --force to skip confirmation")
	}
	styledPrompt := prompt
	if isTTY(os.Stderr) {
		warn := lipgloss.NewStyle().Foreground(ui.Yellow).Render("?")
		hint := lipgloss.NewStyle().Foreground(ui.Gray).Render("[y/N]")
		styledPrompt = fmt.Sprintf("%s %s %s", warn, prompt, hint)
	} else {
		styledPrompt = fmt.Sprintf("%s [y/N]", prompt)
	}
	fmt.Fprintf(os.Stderr, "%s: ", styledPrompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("read confirmation: %w", err)
	}
	answer := strings.TrimSpace(strings.ToLower(line))
	return answer == "y" || answer == "yes", nil
}
