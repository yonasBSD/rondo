package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestInitTheme_Dark(t *testing.T) {
	InitTheme(true)
	if Cyan != lipgloss.Color("#00BCD4") {
		t.Errorf("dark Cyan = %v, want #00BCD4", Cyan)
	}
	if White != lipgloss.Color("#FAFAFA") {
		t.Errorf("dark White = %v, want #FAFAFA", White)
	}
	if SelectionBg != lipgloss.Color("#1a1a2e") {
		t.Errorf("dark SelectionBg = %v, want #1a1a2e", SelectionBg)
	}
	if OverlayDim != lipgloss.Color("#111111") {
		t.Errorf("dark OverlayDim = %v, want #111111", OverlayDim)
	}
}

func TestInitTheme_Light(t *testing.T) {
	InitTheme(false)
	if Cyan != lipgloss.Color("#00838F") {
		t.Errorf("light Cyan = %v, want #00838F", Cyan)
	}
	if White != lipgloss.Color("#1A1A2E") {
		t.Errorf("light White = %v, want #1A1A2E", White)
	}
	if SelectionBg != lipgloss.Color("#F0F0F0") {
		t.Errorf("light SelectionBg = %v, want #F0F0F0", SelectionBg)
	}
	if OverlayDim != lipgloss.Color("#F5F5F5") {
		t.Errorf("light OverlayDim = %v, want #F5F5F5", OverlayDim)
	}
	// Restore dark for other tests.
	InitTheme(true)
}

func TestIsDark(t *testing.T) {
	InitTheme(true)
	if !IsDark() {
		t.Error("IsDark() should be true after InitTheme(true)")
	}
	InitTheme(false)
	if IsDark() {
		t.Error("IsDark() should be false after InitTheme(false)")
	}
	InitTheme(true)
}
