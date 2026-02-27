package app

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/roniel/todo-app/internal/ui"
)

// dialogStyle returns the dialog overlay style with the current theme colors.
func dialogStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ui.Cyan).
		Padding(1, 2).
		Width(60)
}
