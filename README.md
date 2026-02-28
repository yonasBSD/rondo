<h1 align="center">
  RonDO
</h1>

<p align="center">
  <strong>A modern terminal productivity app that combines task management with a daily journal.</strong>
</p>

<p align="center">
  <a href="https://github.com/roniel-rhack/rondo/releases/latest"><img src="https://img.shields.io/github/v/tag/roniel-rhack/rondo?style=flat-square&amp;label=release&amp;color=00bcd4" alt="Release"></a>
  <a href="https://github.com/roniel-rhack/rondo/blob/main/LICENSE"><img src="https://img.shields.io/github/license/roniel-rhack/rondo?style=flat-square&amp;color=00bcd4" alt="License"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&amp;logo=go&amp;logoColor=white" alt="Go"></a>
</p>

<p align="center">
  Fast, keyboard-driven, and built with Go and the <a href="https://charm.sh">Charm</a> ecosystem.
</p>

---

<p align="center">
  <img src="assets/tasks.png" width="720" alt="Task management view">
</p>

<p align="center">
  <img src="assets/journal.png" width="720" alt="Journal view">
</p>

---

## Install

### Homebrew

```bash
brew tap roniel-rhack/tap
brew install rondo
```

### Arch Linux (AUR)

```bash
yay -S rondo
```

### From source

```bash
git clone https://github.com/roniel-rhack/rondo.git
cd rondo
go build -o rondo ./cmd/todo
mv rondo /usr/local/bin/
```

## Features

### Task Management

- **Full CRUD** with validated forms
- **Subtasks** with completion tracking and progress bar
- **Status workflow** — Pending, In Progress, Done
- **Priority levels** — Low, Medium, High, Urgent (color-coded)
- **Due dates** with sort support
- **Tags** — comma-separated, filterable with tag bar (`F4`)
- **Recurring tasks** — daily, weekly, monthly, or yearly recurrence
- **Task dependencies** — mark tasks as blocked by others
- **Time logging** — log time spent on tasks from the detail panel
- **Sorting** by creation date, due date, or priority (`F1`/`F2`/`F3`)
- **Fuzzy search** — filter tasks with `/`

### Productivity Tools

- **Pomodoro timer** — full work/break cycles with configurable durations (`p`)
  - Work → Short Break → Work → ... → Long Break (4-session sets)
  - Phase-aware timer: 🍅 work, ☕ short break, 🌿 long break
  - Cycle indicator (●●●○), terminal bell on completion
  - Configurable via settings form (`P`) or `config.json`
- **Statistics overlay** — task counts, priority breakdown, tag cloud, focus sessions, daily goal, streaks (`G`)
- **Export** — Markdown or JSON, with optional journal inclusion (`X`)
- **Undo** — revert the last destructive action (`Ctrl+Z`)

### Daily Journal

- **One note per day** — auto-created on first entry
- **Timestamped entries** — each entry records the time
- **Edit & delete entries** — cursor-based selection
- **Hide/restore notes** — archive without deleting
- **Smart date labels** — "Today", "Yesterday", weekday names

### Interface

- **Two-panel layout** — list + detail with `1`/`2` focus switching
- **Resizable panels** — adjust split ratio with `<`/`>`
- **Four tabs** — All, Active, Done, Journal (live counts)
- **Vim-style navigation** — `j`/`k` everywhere
- **Context-sensitive status bar** — keybinding hints update per panel
- **Modal forms** — Huh-powered with Dracula theme
- **Confirmation dialogs** for all destructive actions
- **Help overlay** — `?` for full keybinding reference
- **Config file** — persistent settings at `~/.todo-app/config.json`
- **Auto backups** — daily SQLite backups at `~/.todo-app/backups/`

### CLI Mode

Full-featured CLI with styled terminal output (auto-detected), JSON support, and shell completions.

```bash
# Tasks
rondo add "Buy groceries" --priority high --due 2026-03-15 --tags "home,shopping"
rondo list --status pending --sort priority --limit 10
rondo list --priority urgent --overdue --format json
rondo show 3
rondo edit 3 --title "Buy organic groceries" --due 2026-03-20
rondo done 3 4 5
rondo delete 3 --force
rondo status 3 active

# Subtasks
rondo subtask add 3 "Pick up milk"
rondo subtask list 3
rondo subtask done 3 1

# Time tracking
rondo timelog add 3 1h30m --note "Deep work session"
rondo timelog list 3
rondo timelog summary --days 30

# Recurrence
rondo recur set 3 weekly
rondo recur clear 3

# Journal
rondo journal "Productive day"
rondo journal add "Wrapped up the feature" --date yesterday
rondo journal list
rondo journal show today

# Focus (Pomodoro)
rondo focus start --task-id 3 --duration 25m
rondo focus status
rondo focus stats --days 14

# Utilities
rondo stats                        # Task + focus summary
rondo export --format json --journal --output backup.json
rondo config list                  # View all settings
rondo config set focus.work_duration_min 30
rondo completion zsh               # Shell completions
```

#### Global Flags

| Flag | Description |
|------|-------------|
| `--format table\|json` | Output format (default: table) |
| `--json` | Shorthand for `--format json` |
| `-q, --quiet` | Suppress non-essential output |
| `--no-color` | Disable ANSI colors (auto-detected when piped) |

## Keyboard Shortcuts

### Global

| Key | Action |
|-----|--------|
| `Tab` | Switch tabs |
| `1` / `2` | Focus left / right panel |
| `<` / `>` | Resize panels |
| `Esc` | Return to list / clear filter |
| `?` | Help overlay |
| `p` | Pomodoro timer (start / stop / break) |
| `P` | Pomodoro settings |
| `X` | Export |
| `G` | Statistics |
| `Ctrl+Z` | Undo last action |
| `q` | Quit |

### Tasks — Panel 1

| Key | Action |
|-----|--------|
| `j`/`k` | Navigate |
| `a` | Add task |
| `e` | Edit task |
| `d` | Delete task |
| `s` | Cycle status |
| `t` | Add subtask |
| `/` | Search |
| `F1`/`F2`/`F3` | Sort by created / due / priority |
| `F4` | Toggle tag filter bar |

### Tasks — Panel 2 (Details)

| Key | Action |
|-----|--------|
| `j`/`k` | Navigate subtasks |
| `a` | Add subtask |
| `e` | Edit subtask |
| `d` | Delete subtask |
| `s` | Toggle subtask |
| `l` | Log time |
| `b` | View blockers |

### Journal — Panel 1 (Notes)

| Key | Action |
|-----|--------|
| `j`/`k` | Navigate notes |
| `a` | Add entry (today) |
| `h` | Hide / restore note |
| `H` | Toggle show hidden |
| `/` | Search notes |

### Journal — Panel 2 (Entries)

| Key | Action |
|-----|--------|
| `j`/`k` | Navigate entries |
| `a` | Add entry (today) |
| `e` | Edit entry |
| `d` | Delete entry |

## Architecture

```
cmd/todo/main.go                # Entry point (TUI + CLI dispatch)
internal/
  app/
    model.go                    # Main Bubbletea Model + Update + View
    model_journal.go            # Journal tab handlers
    model_forms.go              # Form submission + confirmation dialogs
    model_overlays.go           # Help, stats, blocker overlays + panel renderer
    model_tasks.go              # Task list helpers (filter, sort, reload, export)
    model_features.go           # Feature handlers (focus, tags, undo, blockers)
    keys.go                     # Keybinding definitions
    styles.go                   # Lip Gloss styles
    delegate.go                 # Task list item delegate
    delegate_journal.go         # Journal note list item delegate
  cli/
    cli.go                      # Cobra root command + global flags
    output.go                   # Styled output (TTY-aware tables, colors)
    errors.go                   # NotFoundError type
    confirm.go                  # Confirmation prompts
    tasks.go                    # add, done, list, show, edit, delete, status
    journal.go                  # journal (add, list, show, edit, delete, hide)
    export.go                   # export (md, json, file output)
    subtasks.go                 # subtask (add, list, done, edit, delete)
    timelog.go                  # timelog (add, list, summary)
    recur.go                    # recur (set, clear)
    focus.go                    # focus (start, status, stats)
    stats.go                    # stats (task + focus summary)
    config_cmd.go               # config (list, get, set, reset)
    completion.go               # Shell completion (bash, zsh, fish, powershell)
  config/
    config.go                   # JSON config (~/.todo-app/config.json)
  database/
    db.go                       # SQLite connection + backup
  export/
    export.go                   # Markdown + JSON export writers
  focus/
    focus.go                    # Focus/Pomodoro session model
    store.go                    # Focus session SQLite repository
  journal/
    journal.go                  # Note & Entry domain types
    store.go                    # Journal SQLite repository
  task/
    task.go                     # Task & Subtask domain types
    store.go                    # Task SQLite repository
    recur.go                    # Recurring task logic
    timelog.go                  # Time log model
  ui/
    colors.go                   # Shared color palette
    views.go                    # Rendering (tabs, detail, status bar, dialogs)
    form.go                     # Huh form builders
```

Follows the **Bubbletea MVU** (Model-Update-View) pattern. All data persists in a single SQLite database at `~/.todo-app/todo.db` (WAL mode, single connection, `ON DELETE CASCADE`).

## Development

```bash
go build ./cmd/todo     # Build
go run ./cmd/todo       # Run
go test ./...           # Test
go vet ./...            # Vet
go mod tidy             # Tidy deps
```

## Built With

[Bubbletea](https://github.com/charmbracelet/bubbletea) ·
[Bubbles](https://github.com/charmbracelet/bubbles) ·
[Lip Gloss](https://github.com/charmbracelet/lipgloss) ·
[Huh](https://github.com/charmbracelet/huh) ·
[Cobra](https://github.com/spf13/cobra) ·
[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)

## License

[MIT](LICENSE)
