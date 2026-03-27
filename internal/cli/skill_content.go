package cli

// skillContent is the embedded SKILL.md content installed by `rondo skill install`.
const skillContent = `---
name: rondo
description: Use when managing tasks, journal entries, subtasks, time logs, focus sessions, or tracking work progress. Invoke for any request involving todos, task lists, daily notes, pomodoro, or productivity tracking.
---

# RonDO — Terminal Productivity CLI

Task management, journaling, time tracking, and focus sessions from the command line.
Always use ` + "`" + `--json` + "`" + ` when you need to parse output programmatically.

## Global Flags

All commands support: ` + "`" + `--json` + "`" + `, ` + "`" + `--format table|json` + "`" + `, ` + "`" + `--quiet` + "`" + ` / ` + "`" + `-q` + "`" + `, ` + "`" + `--no-color` + "`" + `

## Tasks

` + "```" + `bash
# Add a task
rondo add "title" [--priority low|medium|high|urgent] [--due YYYY-MM-DD] \
  [--tags t1,t2] [--desc "..."] [--meta key=value] [--blocks 2,3] [--recur daily|weekly|monthly|yearly]

# List tasks (supports rich filtering)
rondo list [--status pending|active|done|all] [--priority high] [--tag work] \
  [--meta key=value] [--sort created|due|priority] [--due-before YYYY-MM-DD] \
  [--due-after YYYY-MM-DD] [--overdue] [--search text] [--limit N] [--json]

# Show task details
rondo show <id> [--json]

# Edit task (only specified flags are updated)
rondo edit <id> [--title "..."] [--desc "..."] [--priority ...] [--due ...] \
  [--tags ...] [--meta key=value] [--blocks 1,2] [--clear-blocks] [--clear-due] [--recur ...]

# Mark done (supports multiple IDs; spawns next for recurring tasks)
rondo done <id> [<id2> ...]

# Delete task (--cascade if it blocks others, --force/-y to skip confirm)
rondo delete <id> [--force] [--cascade]

# Set or cycle status
rondo status <id> [pending|active|done]
` + "```" + `

## Subtasks

` + "```" + `bash
rondo subtask add <task-id> "title"
rondo subtask list <task-id> [--json]
rondo subtask done <task-id> <subtask-id>      # toggles completion
rondo subtask edit <task-id> <subtask-id> "new title"
rondo subtask delete <task-id> <subtask-id> [--force]
` + "```" + `

## Task Notes

` + "```" + `bash
rondo note add <task-id> "note body"
rondo note list <task-id> [--json]
rondo note edit <task-id> <note-id> "new body"
rondo note delete <task-id> <note-id> [--force]
` + "```" + `

## Time Logging

` + "```" + `bash
# Duration format: 1h30m, 45m, 2h
rondo timelog add <task-id> <duration> [--note "what I did"]
rondo timelog list <task-id> [--json]
rondo timelog summary [--days 7] [--json]
` + "```" + `

## Recurrence

` + "```" + `bash
rondo recur set <id> daily|weekly|monthly|yearly
rondo recur clear <id>
` + "```" + `

## Journal

` + "```" + `bash
# Quick add to today (shorthand)
rondo journal "entry text"

# Add with date control
rondo journal add "entry text" [--date today|yesterday|YYYY-MM-DD]

# List notes
rondo journal list [--date YYYY-MM-DD] [--hidden] [--json]

# Show entries for a date (default: today)
rondo journal show [today|yesterday|YYYY-MM-DD] [--json]

# Edit / delete entries
rondo journal edit <entry-id> "new text"
rondo journal delete <entry-id> [--force]

# Toggle note visibility
rondo journal hide <date>
` + "```" + `

## Focus / Pomodoro

` + "```" + `bash
# Record a completed focus session
rondo focus start [--task-id <id>] [--duration 25m]

# Today's progress
rondo focus status [--json]

# Historical stats
rondo focus stats [--days 7] [--json]
` + "```" + `

## Stats & Export

` + "```" + `bash
rondo stats [--json]
rondo export [--format md|json] [--output file.md] [--journal]
` + "```" + `

## Batch Mode

Send multiple commands via stdin as newline-delimited JSON:

` + "```" + `bash
echo '{"cmd":"add","args":["Deploy fix","--priority","urgent"]}
{"cmd":"list","args":["--status","active","--json"]}' | rondo batch
` + "```" + `

Returns JSON array: ` + "`" + `[{"cmd":"add","ok":true}, ...]` + "`" + `

## Config

` + "```" + `bash
rondo config list [--json]
rondo config get <key>
rondo config set <key> <value>
rondo config reset [--force]
` + "```" + `

## Shell Completions

` + "```" + `bash
rondo completion bash|zsh|fish|powershell
` + "```" + `

## Tips

- Use ` + "`" + `--json` + "`" + ` to get structured output for parsing
- Use ` + "`" + `--quiet` + "`" + ` to suppress success messages
- Date fields accept: YYYY-MM-DD, "today", "yesterday"
- Metadata filters use AND logic: ` + "`" + `--meta a=1 --meta b=2` + "`" + ` matches both
- Tag filters use OR logic: ` + "`" + `--tag a --tag b` + "`" + ` matches either
- Delete guard: tasks blocking others need ` + "`" + `--cascade` + "`" + ` to delete
- Recurring tasks auto-spawn next occurrence on ` + "`" + `done` + "`" + `
`
