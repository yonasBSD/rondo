# Pomodoro Enhancement Design

> Research by: Software Architect, UX Designer, UI Designer

## Problem

The current focus feature is a single 25-minute countdown timer. It lacks breaks, cycles, configurable durations, notifications, daily goals, and meaningful stats — all core to the Pomodoro Technique.

## The Pomodoro Technique

```
Work 25m → Short Break 5m → Work 25m → Short Break 5m →
Work 25m → Short Break 5m → Work 25m → Long Break 15m → repeat
```

One "set" = 4 work sessions + 3 short breaks + 1 long break.

---

## Architecture

### State Machine

```
phaseIdle
  --[p key]--> phaseWork (start tick)

phaseWork
  --[timer done]--> phaseWorkDone (show overlay, ring bell)
  --[p key]--> modeFocusConfirmCancel (existing)

phaseWorkDone (overlay: Enter=break, s=skip, Esc=dismiss)
  --[Enter]--> phaseBreak (short or long, based on cycle)
  --[s/Esc]--> phaseIdle

phaseBreak
  --[timer done]--> phaseBreakDone (show overlay, ring bell)
  --[p key]--> modeFocusConfirmCancel

phaseBreakDone (overlay: Enter=work, s=skip, Esc=dismiss)
  --[Enter]--> phaseWork
  --[s/Esc]--> phaseIdle
```

When `autoStartBreak` is ON, `phaseWorkDone` skips the overlay and transitions directly to `phaseBreak`.

### Phase Types

```go
type focusPhase int

const (
    phaseIdle      focusPhase = iota
    phaseWork
    phaseBreak
    phaseWorkDone
    phaseBreakDone
)
```

### New Modes

```go
modeFocusSessionEnd  // work done overlay
modeFocusBreakEnd    // break done overlay
modeFocusSettings    // P key: Huh form for settings
```

### Model Field Changes

Replace:
```go
focusActive  bool
```

With:
```go
focusPhase    focusPhase
focusCyclePos int          // 0-3, completed work sessions in current cycle
```

Add `isFocusActive()` method derived from `focusPhase != phaseIdle`.

### Cycle Logic

```
cyclePos  work completions  next break
    0           1           short
    1           2           short
    2           3           short
    3           4           long → reset to 0
```

`focusCyclePos` is in-memory, recovered on startup via `TodayWorkCount() % LongBreakInterval`.

---

## Data Model

### Domain Type Changes (`internal/focus/focus.go`)

```go
type SessionKind int

const (
    KindWork       SessionKind = 0  // zero value = backward compatible
    KindShortBreak SessionKind = 1
    KindLongBreak  SessionKind = 2
)

type Session struct {
    ID          int64
    TaskID      int64
    Kind        SessionKind   // NEW
    CyclePos    int           // NEW: 1-4 for work, 0 for breaks
    Duration    time.Duration
    StartedAt   time.Time
    CompletedAt *time.Time
}
```

### Schema Migration (`internal/focus/store.go`)

Using `addColumnIfNotExists` pattern (already in `task/store.go`):

```sql
ALTER TABLE focus_sessions ADD COLUMN kind      INTEGER NOT NULL DEFAULT 0;
ALTER TABLE focus_sessions ADD COLUMN cycle_pos INTEGER NOT NULL DEFAULT 0;
```

Existing rows get `kind=0` (Work) — correct since only work sessions existed.

### New Store Methods

```go
func (s *Store) TodayWorkCount() (int, error)
func (s *Store) WeeklySummary() (map[string]int, error)
func (s *Store) Streak() (int, error)
func (s *Store) TotalMinutesFocused(days int) (int, error)
```

---

## Configuration

Extend existing `config.Config` in `internal/config/config.go`:

```go
type FocusConfig struct {
    WorkDuration       int  `json:"work_duration_min"`        // default: 25
    ShortBreakDuration int  `json:"short_break_duration_min"` // default: 5
    LongBreakDuration  int  `json:"long_break_duration_min"`  // default: 15
    LongBreakInterval  int  `json:"long_break_interval"`      // default: 4
    DailyGoal          int  `json:"daily_goal"`               // default: 8 (0=disabled)
    AutoStartBreak     bool `json:"auto_start_break"`         // default: false
    Sound              bool `json:"sound"`                    // default: true (BEL)
}
```

Zero values from old config files → `validate()` applies defaults. No config migration needed.

---

## Keyboard Shortcuts

| Key | Action | Notes |
|-----|--------|-------|
| `p` | Start/cancel focus (unchanged) | Context-aware: starts work, break, or shows cancel dialog |
| `P` | Open focus settings form | New. Follows `h`/`H` pattern |
| `Enter` | Start break/work (in overlays) | Only active in completion overlays |
| `s` | Skip break/work (in overlays) | Only active in completion overlays |

No other new bindings needed. Stats integrate into existing `G` overlay.

---

## UX Flows

### Progressive Disclosure

```
Layer 1: Press p         → Just works (25/5/15 defaults, zero config)
Layer 2: Press P         → Settings form (customize durations, daily goal)
Layer 3: Edit config.json → Power control (all parameters)
Layer 4: Press G          → Stats (streaks, weekly chart, daily progress)
```

### Session Completion

Default (auto_break OFF):
1. Terminal bell (`\a`)
2. Overlay dialog appears with cycle info and choices
3. User picks: Enter (start break), s (skip), Esc (dismiss)

Optional (auto_break ON):
1. Terminal bell
2. Status message: "Break started (5 min)"
3. Break timer starts automatically

### Task-Level Focus Data

Detail panel shows focus info automatically when sessions exist:
```
Focus        3 sessions (1h 15m total)
             Last: Today at 2:30 PM
```

---

## UI Design

### Timer Display

Status bar progress bar (always visible during session):
```
🍅 ████████████░░░░░░░░ 18:30 ●●●○  │  p:stop
```

### Visual States

| State | Emoji | Color | Border |
|-------|-------|-------|--------|
| Work | 🍅 | Orange (`#FF9800`) | Orange |
| Short Break | ☕ | Green (`#4CAF50`) | Green |
| Long Break | 🌿 | Cyan (`#00BCD4`) | Cyan |

All colors already exist in the palette — zero new definitions needed.

### Cycle Indicator

`●●●○` — Green filled circles for completed, DimGray empty for remaining. Terminal-safe 1-width characters.

### Session Complete Overlay

```
╔═══════════════════════════════════════╗
║                                       ║
║     🍅 Session Complete!              ║
║                                       ║
║     3 of 4 sessions done  ●●●○       ║
║     Time for a 5 min break            ║
║                                       ║
║     [Enter] Start break               ║
║     [s] Skip  [Esc] Dismiss           ║
║                                       ║
╚═══════════════════════════════════════╝
```

Green border (completion = success).

### Full Mockup: During Work Session

```
┌──────────────────────────────────────────────────────────────────┐
│  RonDO  │  All (7)  │  Active (4)  │  Done (3)  │  Journal (5) │
├─────────────────────────┬────────────────────────────────────────┤
│ ▸ ◐ ●  Fix login bug   │  Title       Fix login bug             │
│   ○ ●  Update docs     │  Status      ◐ In Progress             │
│   ○ ●  Write tests     │  Priority    High                      │
│   ○ ●  Deploy v2       │  Due         Mar 01, 2026 SOON         │
│   ✓ ●  Setup CI        │                                        │
│                         │  Subtasks    1/3                       │
│                         │  ██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   │
│                         │    [x] Reproduce bug                   │
│                         │    [ ] Write test                      │
│                         │    [ ] Fix handler                     │
├─────────────────────────┴────────────────────────────────────────┤
│ 🍅 ████████████████░░░░░░░░ 18:30 ●●●○ │ [1:Tasks] p:stop ?:help│
└──────────────────────────────────────────────────────────────────┘
```

### Full Mockup: During Short Break

```
┌──────────────────────────────────────────────────────────────────┐
│  RonDO  │  All (7)  │  Active (4)  │  Done (3)  │  Journal (5) │
├─────────────────────────┬────────────────────────────────────────┤
│  (normal task list)     │  (normal detail view)                  │
├─────────────────────────┴────────────────────────────────────────┤
│ ☕ ████████░░░░░░░░░░░░░░ 03:15 ●●●○  │ [1:Tasks] p:skip ?:help│
└──────────────────────────────────────────────────────────────────┘
```

### Enhanced Stats Overlay (Focus Section)

```
Focus
  Today: 5/8 sessions (2h 5m)
  This week:
  Mon ████████  4
  Tue ██████████████  7
  Wed ████████████  6
  Thu ██████  3
  Fri ██████████  5
  Streak: 12 days  ▁▂▃▅▇▅▃▁▁▃▅▇▅▃
```

Uses existing sparkline and progress bar infrastructure.

### Notifications

1. **BEL character** (`\a` via stderr) — universal terminal notification
2. **Overlay dialog** — ensures user sees completion
3. **Window title** (`tea.SetWindowTitle`) — visible in taskbar/dock

---

## Files to Modify

| File | Change |
|------|--------|
| `internal/focus/focus.go` | Add `SessionKind`, `CyclePos` to Session |
| `internal/focus/store.go` | Migration + new query methods |
| `internal/config/config.go` | Add `FocusConfig` struct |
| `internal/app/model.go` | `focusPhase` type, new modes, cycle fields |
| `internal/app/model_features.go` | Rewrite toggle, add phase helpers, `focusBell()` |
| `internal/app/model_overlays.go` | 2 new overlay renderers, enhanced stats |
| `internal/ui/views.go` | Status bar timer enhancement |
| `internal/ui/form.go` | `FocusSettingsForm()` for P key |
| `internal/app/keys.go` | Add `P` binding |

**0 new files. 9 files modified.**

---

## Implementation Phases

### Phase 1: Core Cycle (~200 LOC)

- `SessionKind` + `CyclePos` in domain model
- Schema migration
- `focusPhase` state machine in Model
- Work/break transitions with overlay prompts
- Enhanced `focusTimerStr()` with phase + cycle
- Terminal bell on completion
- `TodayWorkCount()` for cycle recovery on startup

### Phase 2: Configuration (~150 LOC)

- `FocusConfig` in config.go with defaults + validation
- `P` keybinding + `modeFocusSettings`
- `FocusSettingsForm()` Huh form
- Wire config values into session creation
- Auto-break behavior path

### Phase 3: Enhanced Stats (~200 LOC)

- Daily goal progress bar in stats overlay
- 7-day bar chart
- Streak tracking
- Per-task focus data in detail panel
- `WeeklySummary()`, `Streak()`, `TotalMinutesFocused()` store methods

### Phase 4: Polish (~150 LOC)

- Progress bar in status bar
- Phase-colored borders during sessions
- Daily goal indicator in status bar
- Color-coded timer text by phase
- Optional journal entry on session completion

---

## Competitive Research Sources

- [Pomofocus.io](https://pomofocus.io/)
- [Zapier: 6 Best Pomodoro Apps](https://zapier.com/blog/best-pomodoro-apps/)
- [Focus Keeper Blog](https://focuskeeper.co/blog/pomofocus-alternative)
- [Reclaim: Top 11 Pomodoro Timer Apps 2026](https://reclaim.ai/blog/best-pomodoro-timer-apps)
- [Terminal Trove: pomo](https://terminaltrove.com/pomo/)
- [h-sifat/productivity-timer](https://github.com/h-sifat/productivity-timer)
- [Todoist: Pomodoro Technique](https://www.todoist.com/productivity-methods/pomodoro-technique)
