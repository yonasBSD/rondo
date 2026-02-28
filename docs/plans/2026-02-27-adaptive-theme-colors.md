# Adaptive Theme Colors Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix issue #2 — Make all colors adapt to light/dark terminal backgrounds so the app is usable on both.

**Architecture:** Detect terminal background once at startup via `lipgloss.HasDarkBackground()`. Initialize the existing exported color variables in `ui/colors.go` from either a dark or light palette. All downstream code continues using `ui.Cyan`, `ui.White`, etc. unchanged — zero API surface changes. Add two new palette colors (`SelectionBg`, `OverlayDim`) to replace inline hardcoded hex values. Conditionally select Huh form theme. Remove accessibility-problematic `Blink(true)`.

**Tech Stack:** Go, Lip Gloss (lipgloss.HasDarkBackground, lipgloss.Color), Huh (ThemeDracula/ThemeBase)

---

### Task 1: Add theme detection and dual palette to `colors.go`

**Files:**
- Modify: `internal/ui/colors.go`
- Create: `internal/ui/colors_test.go`

**Step 1: Write the failing test**

```go
// internal/ui/colors_test.go
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
	if SelectionBg != lipgloss.Color("#E8F0FE") {
		t.Errorf("light SelectionBg = %v, want #E8F0FE", SelectionBg)
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/ui/ -run TestInitTheme -v`
Expected: FAIL — `InitTheme` and `SelectionBg` not defined.

**Step 3: Write the implementation**

Replace `internal/ui/colors.go` with:

```go
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
		SelectionBg = lipgloss.Color("#E8F0FE")
		OverlayDim = lipgloss.Color("#F5F5F5")
	}
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/ui/ -run TestInitTheme -v && go test ./internal/ui/ -run TestIsDark -v`
Expected: PASS

**Step 5: Run full test suite to verify no regressions**

Run: `go vet ./... && go test ./...`
Expected: All pass. Existing code still compiles because exported variable names are unchanged.

**Step 6: Commit**

```bash
git add internal/ui/colors.go internal/ui/colors_test.go
git commit -m "feat: add dual theme palette with dark/light support

Add InitTheme() to switch all palette colors based on terminal
background detection. Add SelectionBg and OverlayDim to the
palette for inline hardcoded values. Default values remain
dark-theme for backward compatibility."
```

---

### Task 2: Call `InitTheme` at startup in `main.go`

**Files:**
- Modify: `cmd/todo/main.go`

**Step 1: Add theme detection before TUI launch**

In `cmd/todo/main.go`, add the import and call `ui.InitTheme()` before creating the model. Insert after the config load block (after line 68) and before `app.New()` (line 70):

```go
import "github.com/charmbracelet/lipgloss"
import "github.com/roniel/todo-app/internal/ui"
```

```go
	// Detect terminal background and initialize color theme.
	ui.InitTheme(lipgloss.HasDarkBackground())
```

**Step 2: Build to verify compilation**

Run: `go build ./cmd/todo`
Expected: Builds successfully.

**Step 3: Commit**

```bash
git add cmd/todo/main.go
git commit -m "feat: detect terminal background at startup

Call lipgloss.HasDarkBackground() and ui.InitTheme() before
launching the TUI to select appropriate color palette."
```

---

### Task 3: Replace inline `#1a1a2e` with `ui.SelectionBg`

**Files:**
- Modify: `internal/app/delegate.go:123`
- Modify: `internal/app/delegate_journal.go:53,77`

**Step 1: Replace in `delegate.go`**

Line 123, change:
```go
selStyle := lipgloss.NewStyle().Background(lipgloss.Color("#1a1a2e")).Width(availWidth)
```
to:
```go
selStyle := lipgloss.NewStyle().Background(ui.SelectionBg).Width(availWidth)
```

**Step 2: Replace in `delegate_journal.go`**

Line 53, change:
```go
line = lipgloss.NewStyle().Background(lipgloss.Color("#1a1a2e")).Render(line)
```
to:
```go
line = lipgloss.NewStyle().Background(ui.SelectionBg).Render(line)
```

Line 77, change:
```go
line = lipgloss.NewStyle().Background(lipgloss.Color("#1a1a2e")).Render(line)
```
to:
```go
line = lipgloss.NewStyle().Background(ui.SelectionBg).Render(line)
```

**Step 3: Build to verify**

Run: `go build ./cmd/todo`
Expected: Builds successfully.

**Step 4: Commit**

```bash
git add internal/app/delegate.go internal/app/delegate_journal.go
git commit -m "refactor: use ui.SelectionBg instead of hardcoded #1a1a2e"
```

---

### Task 4: Replace inline `#111111` with `ui.OverlayDim`

**Files:**
- Modify: `internal/app/model.go` (10 occurrences: lines 729, 743, 753, 763, 775, 786, 796, 805, 811, 817)
- Modify: `internal/app/model_journal.go` (6 occurrences: lines 490, 505, 522, 529, 535, 545)

**Step 1: Replace all occurrences in `model.go`**

Find and replace all instances of:
```go
lipgloss.WithWhitespaceForeground(lipgloss.Color("#111111"))
```
with:
```go
lipgloss.WithWhitespaceForeground(ui.OverlayDim)
```

**Step 2: Replace all occurrences in `model_journal.go`**

Same replacement — all 6 instances.

**Step 3: Build to verify**

Run: `go build ./cmd/todo`
Expected: Builds successfully.

**Step 4: Commit**

```bash
git add internal/app/model.go internal/app/model_journal.go
git commit -m "refactor: use ui.OverlayDim instead of hardcoded #111111"
```

---

### Task 5: Make Huh form theme adaptive

**Files:**
- Modify: `internal/ui/form.go`

**Step 1: Add a FormTheme helper and replace all 6 ThemeDracula calls**

Add at the top of `form.go` (after imports):

```go
// FormTheme returns the appropriate Huh form theme for the current terminal background.
func FormTheme() *huh.Theme {
	if IsDark() {
		return huh.ThemeDracula()
	}
	return huh.ThemeBase()
}
```

Then replace all 6 occurrences of `.WithTheme(huh.ThemeDracula())` with `.WithTheme(FormTheme())` in:
- `NewTaskForm` (line 86)
- `EditTaskForm` (line 135)
- `SubtaskForm` (line 147)
- `JournalEntryForm` (line 161)
- `ExportForm` (line 180)
- `TimeLogForm` (line 197)

**Step 2: Build to verify**

Run: `go build ./cmd/todo`
Expected: Builds successfully.

**Step 3: Commit**

```bash
git add internal/ui/form.go
git commit -m "feat: adaptive Huh form theme based on terminal background

Use ThemeDracula for dark terminals, ThemeBase for light terminals."
```

---

### Task 6: Remove `Blink(true)` from overdue style

**Files:**
- Modify: `internal/ui/overdue.go:46`

**Step 1: Remove Blink**

Line 46, change:
```go
return lipgloss.NewStyle().Foreground(Red).Bold(true).Blink(true)
```
to:
```go
return lipgloss.NewStyle().Foreground(Red).Bold(true)
```

The `OVERDUE` text badge already provides sufficient emphasis. Blink is a WCAG 2.2.2 concern and unreliable across terminals.

**Step 2: Run existing tests**

Run: `go test ./internal/ui/ -run TestDueStyle -v`
Expected: PASS

**Step 3: Build to verify**

Run: `go build ./cmd/todo`
Expected: Builds successfully.

**Step 4: Commit**

```bash
git add internal/ui/overdue.go
git commit -m "fix: remove Blink from overdue style for accessibility

Blink is a WCAG 2.2.2 concern and unreliable across terminals.
The OVERDUE text badge provides sufficient emphasis."
```

---

### Task 7: Final verification

**Step 1: Run full test suite**

Run: `go vet ./... && go test ./...`
Expected: All pass.

**Step 2: Build the binary**

Run: `go build -o rondo ./cmd/todo`
Expected: Builds successfully.

**Step 3: Verify the color values are correct**

Run: `go test ./internal/ui/ -v`
Expected: All tests pass including the new `TestInitTheme_Dark`, `TestInitTheme_Light`, and `TestIsDark`.

---

## Summary

| Task | Files Changed | Nature |
|------|--------------|--------|
| 1 | `ui/colors.go`, `ui/colors_test.go` | Core: dual palette + `InitTheme()` |
| 2 | `cmd/todo/main.go` | Wiring: detect + init at startup |
| 3 | `app/delegate.go`, `app/delegate_journal.go` | Mechanical: 3 inline replacements |
| 4 | `app/model.go`, `app/model_journal.go` | Mechanical: 16 inline replacements |
| 5 | `ui/form.go` | Feature: adaptive Huh theme |
| 6 | `ui/overdue.go` | A11y fix: remove Blink |
| 7 | (none) | Verification |

**Total: 8 files modified, 1 file created. Zero API signature changes.**

### Light Theme Palette Reference

| Variable | Dark | Light | Light contrast vs white |
|----------|------|-------|------------------------|
| `Cyan` | `#00BCD4` | `#00838F` | 5.4:1 |
| `White` | `#FAFAFA` | `#1A1A2E` | >15:1 |
| `Gray` | `#666666` | `#5C5C5C` | 5.9:1 |
| `DimGray` | `#444444` | `#999999` | 2.8:1 (intentional: dim role) |
| `Green` | `#4CAF50` | `#2E7D32` | 5.8:1 |
| `Red` | `#F44336` | `#C62828` | 7.1:1 |
| `Yellow` | `#FFC107` | `#AB6A00` | 4.8:1 |
| `Magenta` | `#E040FB` | `#9C27B0` | 6.6:1 |
| `Orange` | `#FF9800` | `#E65100` | 5.1:1 |
| `SelectionBg` | `#1a1a2e` | `#E8F0FE` | — |
| `OverlayDim` | `#111111` | `#F5F5F5` | — |
