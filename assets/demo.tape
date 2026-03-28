# RonDO Demo — Terminal Productivity App
# Usage: vhs demo.tape

Output demo.gif

Require rondo

Set Shell zsh
Set Width 1200
Set Height 900
Set FontSize 15
Set Padding 15
Set TypingSpeed 50ms
Set Theme "Catppuccin Mocha"

# ─── Launch RonDO ───
Type "rondo"
Enter
Sleep 2.5s

# ─── Add first task ───
Type "a"
Sleep 1.5s

# Title
Type "Build landing page"
Tab
Sleep 500ms

# Description (Text field — use Alt+Enter or Tab to exit)
Type "Design and implement the marketing site"
Tab
Sleep 500ms

# Priority — Select: use Down to pick "High"
Down
Down
Sleep 300ms
Tab
Sleep 500ms

# Due Date
Type "2026-04-15"
Tab
Sleep 500ms

# Tags
Type "frontend,design"
Tab
Sleep 500ms

# Recurrence — Select: keep "None"
Tab
Sleep 500ms

# Metadata (Text field) — skip empty
Tab
Sleep 500ms

# Blocks — skip empty, submit form
Enter
Sleep 2s

# ─── Add second task ───
Type "a"
Sleep 1.5s

# Title
Type "Set up CI/CD pipeline"
Tab
Sleep 500ms

# Description
Type "Configure GitHub Actions for testing"
Tab
Sleep 500ms

# Priority — keep Low (default)
Tab
Sleep 500ms

# Due Date
Type "2026-04-10"
Tab
Sleep 500ms

# Tags
Type "devops"
Tab
Sleep 500ms

# Recurrence — keep None
Tab
Sleep 500ms

# Metadata — skip
Tab
Sleep 500ms

# Blocks — submit
Enter
Sleep 2s

# ─── Navigate tasks ───
Type "j"
Sleep 600ms
Type "k"
Sleep 1s

# ─── Add subtasks to first task ───
Type "t"
Sleep 1s
Type "Create wireframes"
Tab
Sleep 300ms
Enter
Sleep 1s

Type "t"
Sleep 1s
Type "Implement responsive grid"
Tab
Sleep 300ms
Enter
Sleep 1s

Type "t"
Sleep 1s
Type "Add dark mode support"
Tab
Sleep 300ms
Enter
Sleep 1.5s

# ─── Switch to detail panel and toggle subtasks ───
Type "2"
Sleep 800ms
Type "s"
Sleep 600ms
Type "j"
Sleep 400ms
Type "s"
Sleep 1.5s

# ─── Back to list, cycle task status → In Progress ───
Type "1"
Sleep 500ms
Type "s"
Sleep 1s

# ─── Switch tabs: Active ───
Tab
Sleep 1.5s

# ─── Switch tabs: Done ───
Tab
Sleep 1.5s

# ─── Switch tabs: Journal ───
Tab
Sleep 1.5s

# ─── Add journal entry ───
Type "a"
Sleep 1s
Type "Started working on the landing page. Good progress on wireframes today."
Sleep 500ms
Enter
Sleep 1.5s

# ─── Add another journal entry ───
Type "a"
Sleep 1s
Type "CI/CD pipeline research done. GitHub Actions looks like the best option."
Sleep 500ms
Enter
Sleep 1.5s

# ─── Show journal entries in detail panel ───
Type "2"
Sleep 800ms
Type "j"
Sleep 600ms
Type "k"
Sleep 1.5s

# ─── Back to All tab ───
Type "1"
Sleep 500ms
Tab
Sleep 1.5s

# ─── Help overlay ───
Type "?"
Sleep 3s
Escape
Sleep 1.5s

# ─── Search/filter ───
Type "/"
Sleep 600ms
Type "landing"
Sleep 1.5s
Escape
Sleep 2s

# ─── Quit ───
Type "q"
Sleep 500ms
