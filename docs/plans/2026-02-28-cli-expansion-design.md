# CLI Expansion Design

## Problem

RonDO's CLI exposes 5 commands (`add`, `done`, `list`, `journal`, `export`) out of 20+ features available in the TUI. The data layer already supports every feature — the gap is purely in CLI exposure. Power users cannot script task editing, time logging, journal browsing, Pomodoro tracking, or subtask management without opening the TUI.

## Current CLI Surface

| Command | Flags | Notes |
|---------|-------|-------|
| `rondo add "title"` | `--priority`, `--due`, `--tags` | Missing: `--description`, `--recur` |
| `rondo done <id>` | none | Only sets Done; no InProgress. Uses `List()` scan instead of `GetByID` |
| `rondo list` | `--status`, `--format` | No `--tag`, `--sort`, `--priority` filters. JSON output omits description, subtasks, time logs |
| `rondo journal "text"` | none | Append-only to today. No list/show/edit/delete |
| `rondo export` | `--format`, `--output`, `--journal` | Works well as-is |

**Architectural gap:** `cli.Run()` receives only `taskStore` and `journalStore`. The `focusStore` and `config` are never passed, making focus and config commands impossible without a signature change.

## Design Decisions

### 1. CLI Framework: Adopt Cobra

The current stdlib `flag` + `switch` dispatch works for 5 commands but cannot scale to 30+ with subcommand nesting (`subtask add`, `journal list`, `focus start`).

Cobra provides:
- Subcommand nesting with automatic help generation
- Persistent flags (`--format`, `--quiet`, `--no-color`) inherited by all subcommands
- Built-in shell completion for bash/zsh/fish/powershell
- `SilenceUsage`/`SilenceErrors` for controlled error formatting

Trade-off: adds one dependency. Worth it given the scope of expansion.

### 2. Command Style: Verb-First with Groups

Keep the current verb-first pattern (`rondo add`, `rondo list`) for backward compatibility. Use noun-verb groups only for non-task domains:

```
# Task commands (top-level, backward compatible)
rondo add, list, show, edit, delete, done, status

# Grouped commands (new domains)
rondo subtask add|list|done|delete
rondo journal  add|list|show|edit|delete|hide
rondo timelog  add|list|summary
rondo focus    start|status|stats
rondo config   get|set|list|reset
rondo recur    set|clear
rondo stats
rondo completion
```

`rondo journal "text"` remains a backward-compatible alias for `rondo journal add "text"`.

### 3. Output and Formatting

| Rule | Implementation |
|------|---------------|
| TTY detection | Auto-disable ANSI when stdout is not a terminal |
| `--format table\|json\|plain` | Persistent flag on root command |
| `--json` | Shorthand for `--format json` |
| `--quiet` / `-q` | Suppress confirmations; on `add`, emit only the created ID |
| `--no-color` | Force plain output even on TTY |
| Errors | Always to stderr |

### 4. Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General/internal error |
| 2 | Usage error (bad flags, missing args) |
| 3 | Resource not found |

### 5. Destructive Operations

All destructive commands (`delete`, `journal delete`, `subtask delete`, `config reset`) prompt for confirmation unless `--force` / `-y` is passed. When stdin is not a TTY, require `--force` explicitly.

## Command Specifications

### Task Commands

#### `rondo add "title" [flags]`
Existing command, extended with:
- `--desc TEXT` — task description
- `--recur daily|weekly|monthly|yearly` — recurrence frequency

#### `rondo show <id> [--format table|json]`
Display full task detail: title, description, status, priority, due date, tags, recurrence, subtasks with progress, time logs, blockers.

#### `rondo edit <id> [flags]`
Patch semantics — only explicitly provided flags update the task.
- `--title TEXT`
- `--desc TEXT`
- `--priority low|medium|high|urgent`
- `--due YYYY-MM-DD`
- `--tags tag1,tag2` (replaces all tags)
- `--clear-due` (removes due date)

#### `rondo delete <id> [--force]`
Delete task with confirmation prompt.

#### `rondo status <id> [pending|active|done]`
Set status explicitly. Without a value, cycle to the next status (same as `s` in TUI). Replaces the limited `done` command for new workflows.

#### `rondo done <id...>`
Mark one or more tasks as done. Fix: use `GetByID` instead of `List()` scan. When completing a recurring task, spawn the next occurrence (matches TUI behavior).

#### `rondo list [flags]`
Extend with filters:
- `--priority low|medium|high|urgent`
- `--tag TAG` (repeatable)
- `--sort created|due|priority` (default: created)
- `--due-before YYYY-MM-DD`
- `--due-after YYYY-MM-DD`
- `--overdue` (due date < today, status != done)
- `--search TEXT` (title/description substring match)
- `--limit N`

Fix JSON output to include all fields: description, subtasks, time logs, recurrence, blocked_by.

### Subtask Commands

```
rondo subtask add <task-id> "title"
rondo subtask list <task-id> [--format table|json]
rondo subtask done <task-id> <subtask-id>
rondo subtask edit <task-id> <subtask-id> "new title"
rondo subtask delete <task-id> <subtask-id> [--force]
```

All store methods exist: `AddSubtask`, `ToggleSubtask`, `UpdateSubtask`, `DeleteSubtask`.

### Time Log Commands

```
rondo timelog add <task-id> <duration> [--note TEXT]
rondo timelog list <task-id> [--format table|json]
rondo timelog summary [--days N] [--format table|json]
```

Duration format: `1h30m`, `45m`, `2h` (parsed by existing `task.ParseDuration`). Store methods: `AddTimeLog`, `ListTimeLogs`.

### Journal Commands

```
rondo journal add "text" [--date YYYY-MM-DD]
rondo journal list [--date YYYY-MM-DD] [--hidden] [--format table|json]
rondo journal show [today|yesterday|YYYY-MM-DD] [--format table|json]
rondo journal edit <entry-id> "new text"
rondo journal delete <entry-id> [--force]
rondo journal hide <date>
```

All store methods exist: `ListNotes`, `ListEntries`, `UpdateEntry`, `DeleteEntry`, `ToggleHidden`.

### Focus/Pomodoro Commands

```
rondo focus start [--task-id ID] [--duration 25m]
rondo focus status [--format table|json]
rondo focus stats [--days N] [--format table|json]
```

`focus status` is the key integration point for tmux/waybar/i3bar:
```bash
# tmux status line
rondo focus status --format json | jq -r '.remaining // empty'
```

Store methods: `Create`, `Complete`, `TodayWorkCount`, `Streak`, `WeeklySummary`, `TotalMinutesFocused`.

**Requires:** passing `focusStore` to `cli.Run()`.

### Config Commands

```
rondo config list
rondo config get <key>
rondo config set <key> <value>
rondo config reset [--force]
```

Keys: `panel_ratio`, `focus.work_duration_min`, `focus.short_break_duration_min`, `focus.long_break_duration_min`, `focus.long_break_interval`, `focus.daily_goal`, `focus.auto_start_break`, `focus.sound`.

### Stats Command

```
rondo stats [--format table|json]
```

Summary dashboard: task counts by status/priority, focus sessions today with goal progress, streak, total time logged. All data available from existing store methods.

### Shell Completion

```
rondo completion bash|zsh|fish|powershell
```

Cobra generates these automatically via `GenBashCompletion`, `GenZshCompletion`, etc.

## Architecture

### CLI Struct (dependency injection)

```go
type CLI struct {
    taskStore    *task.Store
    journalStore *journal.Store
    focusStore   *focus.Store
    cfg          config.Config
    format       string
    quiet        bool
    noColor      bool
}
```

### Entry Point

```go
func Run(args []string, ts *task.Store, js *journal.Store,
    fs *focus.Store, cfg config.Config) error {
    root := New(ts, js, fs, cfg)
    root.SetArgs(args)
    return root.Execute()
}
```

### main.go Change

```go
if err := cli.Run(os.Args[1:], taskStore, journalStore, focusStore, cfg); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    if cli.IsNotFound(err) {
        os.Exit(3)
    }
    os.Exit(1)
}
```

### Shared Printer

```go
type Printer struct {
    format  string
    quiet   bool
    noColor bool
    w       io.Writer
}

func (p *Printer) Success(format string, args ...any)  // suppressed by --quiet
func (p *Printer) Table(headers []string, rows [][]string) // tabwriter
func (p *Printer) JSON(v any) error                         // indented JSON
```

### Error Types

```go
type NotFoundError struct {
    Type string // "task", "subtask", "entry", "note"
    ID   int64
}
```

### File Structure

```
internal/cli/
  cli.go         # Cobra root command, CLI struct, Run()
  tasks.go       # add, list, show, edit, delete, done, status
  subtasks.go    # subtask add|list|done|edit|delete
  timelog.go     # timelog add|list|summary
  focus.go       # focus start|status|stats
  journal.go     # journal add|list|show|edit|delete|hide
  config.go      # config get|set|list|reset
  stats.go       # stats
  recur.go       # recur set|clear
  export.go      # export (existing, adapted to Cobra)
  output.go      # Printer
  errors.go      # NotFoundError
  confirm.go     # confirmation prompt helper
  completion.go  # completion command
```

## Implementation Phases

### Phase 1 — Structural Refactor
- `go get github.com/spf13/cobra`
- Rewrite `cli.go` with Cobra root command and `CLI` struct
- Migrate existing 5 commands to Cobra (identical behavior)
- Update `main.go` to pass `focusStore` and `cfg`
- Add `output.go` (Printer), `errors.go`, `confirm.go`
- Update tests

### Phase 2 — Core Task CRUD
- Add `show`, `edit`, `delete`, `status`
- Fix `done` to use `GetByID` and spawn recurring task on completion
- Extend `list` with `--tag`, `--sort`, `--priority`, `--due-before/after`, `--overdue`, `--search`
- Extend `add` with `--desc`, `--recur`
- Fix JSON output to include all fields

### Phase 3 — Subcommand Groups
- `subtask` group: add, list, done, edit, delete
- `journal` group: add, list, show, edit, delete, hide
- `timelog` group: add, list, summary
- `config` group: get, set, list, reset
- `recur` group: set, clear
- `stats` command

### Phase 4 — Focus + Polish
- `focus` group: start, status, stats
- `completion` command (bash/zsh/fish/powershell)
- Semantic exit codes in `main.go`
- TTY detection for auto-disable ANSI
- `--no-color` + `NO_COLOR` env var support

## Data Layer Readiness

No new store methods needed. All 17 task.Store methods, 8 journal.Store methods, and 9 focus.Store methods already exist. The CLI currently uses only 5 of 34 available methods.

## Backward Compatibility

- `rondo add/done/list/export` — all existing flags preserved
- `rondo journal "text"` — works as alias for `journal add`
- `rondo done <id>` — kept alongside new `status` command
- No existing command changes behavior
