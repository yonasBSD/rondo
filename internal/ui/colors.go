package ui

import "github.com/charmbracelet/lipgloss"

// Shared color palette for the entire application.
// Colors are initialized by InitTheme() based on terminal background detection.
// Default values are for dark terminals (backward compatible).
var (
	Cyan    = lipgloss.Color("#00BCD4")
	White   = lipgloss.Color("#FAFAFA")
	Gray    = lipgloss.Color("#666666")
	DimGray = lipgloss.Color("#444444")
	Green   = lipgloss.Color("#4CAF50")
	Red     = lipgloss.Color("#F44336")
	Yellow  = lipgloss.Color("#FFC107")
	Magenta = lipgloss.Color("#E040FB")
	Orange  = lipgloss.Color("#FF9800")

	// SelectionBg is the background color for selected list items.
	SelectionBg = lipgloss.Color("#1a1a2e")

	// OverlayDim is the whitespace fill color for dialog overlays.
	OverlayDim = lipgloss.Color("#111111")
)

// isDarkTheme tracks the current theme for conditional logic (e.g. Huh form theme).
var isDarkTheme = true

// IsDark returns whether the current theme is dark.
func IsDark() bool {
	return isDarkTheme
}

// InitTheme sets all palette colors based on the terminal background.
// Call once at startup with the result of lipgloss.HasDarkBackground().
func InitTheme(dark bool) {
	isDarkTheme = dark
	if dark {
		Cyan = lipgloss.Color("#00BCD4")
		White = lipgloss.Color("#FAFAFA")
		Gray = lipgloss.Color("#666666")
		DimGray = lipgloss.Color("#444444")
		Green = lipgloss.Color("#4CAF50")
		Red = lipgloss.Color("#F44336")
		Yellow = lipgloss.Color("#FFC107")
		Magenta = lipgloss.Color("#E040FB")
		Orange = lipgloss.Color("#FF9800")
		SelectionBg = lipgloss.Color("#1a1a2e")
		OverlayDim = lipgloss.Color("#111111")
	} else {
		Cyan = lipgloss.Color("#00838F")
		White = lipgloss.Color("#1A1A2E")
		Gray = lipgloss.Color("#5C5C5C")
		DimGray = lipgloss.Color("#999999")
		Green = lipgloss.Color("#2E7D32")
		Red = lipgloss.Color("#C62828")
		Yellow = lipgloss.Color("#AB6A00")
		Magenta = lipgloss.Color("#9C27B0")
		Orange = lipgloss.Color("#E65100")
		SelectionBg = lipgloss.Color("#F0F0F0")
		OverlayDim = lipgloss.Color("#F5F5F5")
	}
}
