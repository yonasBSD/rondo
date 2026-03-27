# Changelog

All notable changes to RonDO are documented here.

---

## [v0.7.0] — 2026-03-27

### ✨ New Features

- **Claude Code Integration**: Install a skill directly from the CLI with `rondo skill install`. Once installed, you can manage your tasks and journal without leaving your editor — Claude will pick up `rondo` commands automatically based on context.
  - `rondo skill install` — global install (`~/.claude/skills/rondo/`)
  - `rondo skill install --project` — project-scoped install (`.claude/skills/rondo/`)
  - `rondo skill uninstall` — removes the skill

---

## [v0.6.1] — 2026-03-27

### 🐛 Bug Fixes

- **Journal scroll**: The detail panel no longer jumps back to the top when navigating between journal entries with `j`/`k`. The view now follows the selected entry smoothly, keeping it visible within the panel.

---

## [v0.6.0] — 2026-03-04

### ✨ New Features

- **Task Metadata**: Attach structured key-value data to any task with `--meta key=value`. Filter tasks by metadata with `rondo list --meta key=value` (supports multiple filters with AND logic).
- **Task Notes**: Add timestamped comments to tasks (`task_notes`) — useful for progress updates or context that doesn't belong in the description.
- **Batch Mode**: `rondo batch` reads newline-delimited JSON commands from stdin and returns results as JSON. Ideal for scripting and automation pipelines.
- **Task Dependencies (TUI)**: Block relationships between tasks are now fully visible and manageable from the TUI, not just the CLI.

### 🔒 Reliability

- **Delete Guard**: RonDO now refuses to delete tasks that block other tasks. The CLI returns exit code 1 with a clear error unless `--cascade` is passed. The TUI shows a yellow warning dialog requiring double confirmation.
- **Self-Block Prevention**: Tasks can no longer be set to block themselves.
- **UTC Timestamps**: All stored timestamps now use UTC consistently to avoid timezone-related bugs.
- **JSON Compatibility**: `blocked_by` and `blocks` fields in JSON output remain backward-compatible `[]int64` arrays, with new `_detail` fields added alongside them.

---

## [v0.5.0] — 2026-03-03

### ✨ New Features

- **Configurable Date & Time Formats**: Choose how dates and times are displayed across the app. Set your preferred format via `rondo config set` or pick from built-in presets (ISO, US, European, and more).

### 🔧 Improvements

- Consistent date rendering throughout the UI — one unified render path instead of scattered format strings.
- Added preset layout options to quickly switch between common date/time display styles.

---

## [v0.4.1] — 2026-02-28

### 🔧 Improvements

- Added a commit-message hook to enforce clean, attribution-free commit messages in the repository.

---

## [v0.4.0] — 2026-02-28

### ✨ New Features

- **Full CLI**: RonDO now ships a complete Cobra-based CLI alongside the TUI. Manage tasks, journal entries, and more from any script or terminal without launching the interactive interface.
  - Commands: `add`, `done`, `list`, `show`, `edit`, `delete`, `status`, `subtask`, `timelog`, `note`, `journal`, `focus`, `export`, `stats`, `config`, `completion`
  - **TTY-aware output**: Styled Unicode tables with color when connected to a terminal; clean plain-text when piped.
  - **Shell completions**: `rondo completion bash|zsh|fish|powershell`
  - **Global flags**: `--format table|json`, `--json`, `--quiet`, `--no-color`

---

## [v0.3.0] — 2026-02-27

### ✨ New Features

- **Pomodoro Timer**: Full focus timer built into the app. Press `P` to start a session linked to your selected task.
  - Work → Short Break → Work → ... → Long Break (4-session sets)
  - Visual cycle indicator (●●●○) and distinct colors per session type
  - Configurable durations, daily goal, and auto-start breaks
  - Terminal bell notification when a phase completes
  - Stats overlay (`G`) with daily goal progress, weekly summary, and streak tracking

---

## [v0.2.3] — 2026-02-27

### 🎨 Improvements

- **Adaptive Color Theme**: The app now automatically detects whether your terminal is using a light or dark background and adjusts its color palette accordingly.

---

## [v0.2.2] — 2026-02-25

### 🐛 Bug Fixes

- Fixed a layout overflow bug where long task titles would push UI elements out of bounds.

---

## [v0.2.1] — 2026-02-23

### 🔧 Improvements

- Added AUR (Arch User Repository) publishing to the release workflow — RonDO is now installable via `yay -S rondo` on Arch Linux.

---

## [v0.2.0] — 2026-02-20

### 🐛 Bug Fixes

- Fixed search/filter not working correctly in the task list.
- Fixed the help dialog overflowing in smaller terminals.
- Suppressed a spurious backup warning that appeared on first launch.

---

## [v0.1.0] — 2026-02-20

### 🚀 Initial Release

RonDO is a keyboard-driven terminal productivity app combining task management and a daily journal in a single TUI.

**Task Management**
- Create, edit, and delete tasks with priorities, due dates, tags, and subtasks
- Cycle status between Pending, In Progress, and Done
- Sort by creation date, due date, or priority (F1/F2/F3)
- Fuzzy search and filter
- Time logging and recurring tasks
- Export to Markdown or JSON

**Journal**
- One note per calendar day, auto-created
- Multiple timestamped entries per note
- Hide/restore old notes

**General**
- Vim-style navigation (`j`/`k`) with two-panel resizable layout
- SQLite persistence at `~/.todo-app/todo.db`
- Auto backups and JSON config
- Undo last destructive action (`Ctrl+Z`)
- Homebrew distribution via automated release workflow
