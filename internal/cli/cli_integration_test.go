package cli

import (
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/roniel/todo-app/internal/config"
	"github.com/roniel/todo-app/internal/journal"
	"github.com/roniel/todo-app/internal/task"

	_ "modernc.org/sqlite"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func newTestStores(t *testing.T) (*task.Store, *journal.Store) {
	t.Helper()
	db := openTestDB(t)
	ts, err := task.NewStore(db)
	if err != nil {
		t.Fatalf("task.NewStore: %v", err)
	}
	js, err := journal.NewStore(db)
	if err != nil {
		t.Fatalf("journal.NewStore: %v", err)
	}
	return ts, js
}

// run is a convenience wrapper that calls Run with nil focusStore and default config.
func run(t *testing.T, args []string, ts *task.Store, js *journal.Store) error {
	t.Helper()
	return Run(args, ts, js, nil, config.Config{})
}

// captureStdout runs fn while capturing everything written to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = orig

	buf, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return string(buf)
}

// ---------------------------------------------------------------------------
// add command
// ---------------------------------------------------------------------------

func TestIntegration_Add_Basic(t *testing.T) {
	ts, js := newTestStores(t)

	if err := run(t, []string{"add", "Buy milk"}, ts, js); err != nil {
		t.Fatalf("add: %v", err)
	}

	tasks, err := ts.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	got := tasks[0]
	if got.Title != "Buy milk" {
		t.Errorf("Title = %q, want %q", got.Title, "Buy milk")
	}
	if got.Priority != task.Low {
		t.Errorf("Priority = %v, want Low", got.Priority)
	}
	if got.Status != task.Pending {
		t.Errorf("Status = %v, want Pending", got.Status)
	}
}

func TestIntegration_Add_AllFlags(t *testing.T) {
	ts, js := newTestStores(t)

	args := []string{"add", "--priority", "high", "--due", "2026-03-15", "--tags", "home,shopping", "Big task"}
	if err := run(t, args, ts, js); err != nil {
		t.Fatalf("add: %v", err)
	}

	tasks, err := ts.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	got := tasks[0]
	if got.Title != "Big task" {
		t.Errorf("Title = %q, want %q", got.Title, "Big task")
	}
	if got.Priority != task.High {
		t.Errorf("Priority = %v, want High", got.Priority)
	}
	if got.DueDate == nil {
		t.Fatal("expected DueDate to be set")
	}
	if got.DueDate.Format("2006-01-02") != "2026-03-15" {
		t.Errorf("DueDate = %s, want 2026-03-15", got.DueDate.Format("2006-01-02"))
	}
	if len(got.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(got.Tags))
	}
	if got.Tags[0] != "home" || got.Tags[1] != "shopping" {
		t.Errorf("Tags = %v, want [home shopping]", got.Tags)
	}
}

func TestIntegration_Add_Multiple(t *testing.T) {
	ts, js := newTestStores(t)

	for _, title := range []string{"Task 1", "Task 2", "Task 3"} {
		if err := run(t, []string{"add", title}, ts, js); err != nil {
			t.Fatalf("add(%q): %v", title, err)
		}
	}

	tasks, err := ts.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}
}

// ---------------------------------------------------------------------------
// done command
// ---------------------------------------------------------------------------

func TestIntegration_Done(t *testing.T) {
	ts, js := newTestStores(t)

	if err := run(t, []string{"add", "Finish report"}, ts, js); err != nil {
		t.Fatalf("add: %v", err)
	}

	tasks, err := ts.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	id := tasks[0].ID

	if err := run(t, []string{"done", strconv.FormatInt(id, 10)}, ts, js); err != nil {
		t.Fatalf("done: %v", err)
	}

	got, err := ts.GetByID(id)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Status != task.Done {
		t.Errorf("Status = %v, want Done", got.Status)
	}
}

func TestIntegration_Done_NotFound(t *testing.T) {
	ts, js := newTestStores(t)

	err := run(t, []string{"done", "999"}, ts, js)
	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "not found")
	}
}

// ---------------------------------------------------------------------------
// list command
// ---------------------------------------------------------------------------

func TestIntegration_List_Table(t *testing.T) {
	ts, js := newTestStores(t)

	run(t, []string{"add", "Alpha"}, ts, js)
	run(t, []string{"add", "Beta"}, ts, js)

	out := captureStdout(t, func() {
		if err := run(t, []string{"list"}, ts, js); err != nil {
			t.Fatalf("list: %v", err)
		}
	})

	if !strings.Contains(out, "Alpha") {
		t.Errorf("output missing %q:\n%s", "Alpha", out)
	}
	if !strings.Contains(out, "Beta") {
		t.Errorf("output missing %q:\n%s", "Beta", out)
	}
	if !strings.Contains(out, "TITLE") {
		t.Errorf("output missing header %q:\n%s", "TITLE", out)
	}
}

func TestIntegration_List_JSON(t *testing.T) {
	ts, js := newTestStores(t)

	run(t, []string{"add", "JSON task"}, ts, js)

	out := captureStdout(t, func() {
		if err := run(t, []string{"list", "--format", "json"}, ts, js); err != nil {
			t.Fatalf("list json: %v", err)
		}
	})

	var result []json.RawMessage
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "JSON task") {
		t.Errorf("output missing %q:\n%s", "JSON task", out)
	}
}

func TestIntegration_List_FilterDone(t *testing.T) {
	ts, js := newTestStores(t)

	run(t, []string{"add", "Stay pending"}, ts, js)
	run(t, []string{"add", "Mark done"}, ts, js)
	run(t, []string{"done", "2"}, ts, js)

	out := captureStdout(t, func() {
		if err := run(t, []string{"list", "--status", "done"}, ts, js); err != nil {
			t.Fatalf("list --status done: %v", err)
		}
	})

	if !strings.Contains(out, "Mark done") {
		t.Errorf("output missing done task:\n%s", out)
	}
	if strings.Contains(out, "Stay pending") {
		t.Errorf("output should not contain pending task:\n%s", out)
	}
}

func TestIntegration_List_FilterPending(t *testing.T) {
	ts, js := newTestStores(t)

	run(t, []string{"add", "Pending one"}, ts, js)
	run(t, []string{"add", "Done one"}, ts, js)
	run(t, []string{"done", "2"}, ts, js)

	out := captureStdout(t, func() {
		if err := run(t, []string{"list", "--status", "pending"}, ts, js); err != nil {
			t.Fatalf("list --status pending: %v", err)
		}
	})

	if !strings.Contains(out, "Pending one") {
		t.Errorf("output missing pending task:\n%s", out)
	}
	if strings.Contains(out, "Done one") {
		t.Errorf("output should not contain done task:\n%s", out)
	}
}

// ---------------------------------------------------------------------------
// journal command
// ---------------------------------------------------------------------------

func TestIntegration_JournalAdd(t *testing.T) {
	ts, js := newTestStores(t)

	if err := run(t, []string{"journal", "Great day"}, ts, js); err != nil {
		t.Fatalf("journal: %v", err)
	}

	notes, err := js.ListNotes(false)
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}

	entries := notes[0].Entries
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Body != "Great day" {
		t.Errorf("Body = %q, want %q", entries[0].Body, "Great day")
	}
}

func TestIntegration_JournalAdd_MultipleWords(t *testing.T) {
	ts, js := newTestStores(t)

	if err := run(t, []string{"journal", "Hello", "world"}, ts, js); err != nil {
		t.Fatalf("journal: %v", err)
	}

	notes, err := js.ListNotes(false)
	if err != nil {
		t.Fatalf("ListNotes: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Entries[0].Body != "Hello world" {
		t.Errorf("Body = %q, want %q", notes[0].Entries[0].Body, "Hello world")
	}
}

// ---------------------------------------------------------------------------
// export command
// ---------------------------------------------------------------------------

func TestIntegration_Export_Markdown(t *testing.T) {
	ts, js := newTestStores(t)

	run(t, []string{"add", "Export me"}, ts, js)

	out := captureStdout(t, func() {
		if err := run(t, []string{"export", "--format", "md"}, ts, js); err != nil {
			t.Fatalf("export md: %v", err)
		}
	})

	if !strings.Contains(out, "# Tasks") {
		t.Errorf("output missing %q:\n%s", "# Tasks", out)
	}
	if !strings.Contains(out, "Export me") {
		t.Errorf("output missing task title:\n%s", out)
	}
}

func TestIntegration_Export_JSON(t *testing.T) {
	ts, js := newTestStores(t)

	run(t, []string{"add", "JSON export"}, ts, js)

	out := captureStdout(t, func() {
		if err := run(t, []string{"export", "--format", "json"}, ts, js); err != nil {
			t.Fatalf("export json: %v", err)
		}
	})

	var data map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "JSON export") {
		t.Errorf("output missing task title:\n%s", out)
	}
}

func TestIntegration_Export_ToFile(t *testing.T) {
	ts, js := newTestStores(t)

	run(t, []string{"add", "File export"}, ts, js)

	tmpFile := filepath.Join(t.TempDir(), "export.md")
	if err := run(t, []string{"export", "--format", "md", "--output", tmpFile}, ts, js); err != nil {
		t.Fatalf("export to file: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("read export file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "# Tasks") {
		t.Errorf("file missing %q:\n%s", "# Tasks", content)
	}
	if !strings.Contains(content, "File export") {
		t.Errorf("file missing task title:\n%s", content)
	}
}

func TestIntegration_Export_WithJournal(t *testing.T) {
	ts, js := newTestStores(t)

	run(t, []string{"add", "My task"}, ts, js)
	run(t, []string{"journal", "My entry"}, ts, js)

	out := captureStdout(t, func() {
		if err := run(t, []string{"export", "--format", "md", "--journal"}, ts, js); err != nil {
			t.Fatalf("export with journal: %v", err)
		}
	})

	if !strings.Contains(out, "# Tasks") {
		t.Errorf("output missing %q:\n%s", "# Tasks", out)
	}
	if !strings.Contains(out, "# Journal") {
		t.Errorf("output missing %q:\n%s", "# Journal", out)
	}
}

// ---------------------------------------------------------------------------
// show command
// ---------------------------------------------------------------------------

func TestIntegration_Show_Table(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "--desc", "My description", "Show me"}, ts, js)

	tasks, _ := ts.List()
	id := tasks[0].ID

	out := captureStdout(t, func() {
		if err := run(t, []string{"show", strconv.FormatInt(id, 10)}, ts, js); err != nil {
			t.Fatalf("show: %v", err)
		}
	})

	if !strings.Contains(out, "Show me") {
		t.Errorf("output missing title:\n%s", out)
	}
	if !strings.Contains(out, "My description") {
		t.Errorf("output missing description:\n%s", out)
	}
}

func TestIntegration_Show_JSON(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "--desc", "desc here", "JSON show"}, ts, js)
	tasks, _ := ts.List()
	id := tasks[0].ID

	out := captureStdout(t, func() {
		if err := run(t, []string{"show", "--format", "json", strconv.FormatInt(id, 10)}, ts, js); err != nil {
			t.Fatalf("show json: %v", err)
		}
	})

	var result map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "JSON show") {
		t.Errorf("output missing title:\n%s", out)
	}
	if !strings.Contains(out, "desc here") {
		t.Errorf("output missing description:\n%s", out)
	}
}

func TestIntegration_Show_NotFound(t *testing.T) {
	ts, js := newTestStores(t)
	err := run(t, []string{"show", "999"}, ts, js)
	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "not found")
	}
}

// ---------------------------------------------------------------------------
// edit command
// ---------------------------------------------------------------------------

func TestIntegration_Edit_Title(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "Original title"}, ts, js)
	tasks, _ := ts.List()
	id := strconv.FormatInt(tasks[0].ID, 10)

	if err := run(t, []string{"edit", id, "--title", "Updated title"}, ts, js); err != nil {
		t.Fatalf("edit: %v", err)
	}

	got, err := ts.GetByID(tasks[0].ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Title != "Updated title" {
		t.Errorf("Title = %q, want %q", got.Title, "Updated title")
	}
}

func TestIntegration_Edit_Priority(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "My task"}, ts, js)
	tasks, _ := ts.List()
	id := strconv.FormatInt(tasks[0].ID, 10)

	if err := run(t, []string{"edit", id, "--priority", "urgent"}, ts, js); err != nil {
		t.Fatalf("edit: %v", err)
	}

	got, _ := ts.GetByID(tasks[0].ID)
	if got.Priority != task.Urgent {
		t.Errorf("Priority = %v, want Urgent", got.Priority)
	}
}

func TestIntegration_Edit_NoFlags(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "My task"}, ts, js)
	tasks, _ := ts.List()

	err := run(t, []string{"edit", strconv.FormatInt(tasks[0].ID, 10)}, ts, js)
	if err == nil {
		t.Fatal("expected error for edit with no flags, got nil")
	}
}

func TestIntegration_Edit_NotFound(t *testing.T) {
	ts, js := newTestStores(t)
	err := run(t, []string{"edit", "999", "--title", "x"}, ts, js)
	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "not found")
	}
}

// ---------------------------------------------------------------------------
// delete command
// ---------------------------------------------------------------------------

func TestIntegration_Delete_Force(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "Bye"}, ts, js)
	tasks, _ := ts.List()
	id := strconv.FormatInt(tasks[0].ID, 10)

	if err := run(t, []string{"delete", "--force", id}, ts, js); err != nil {
		t.Fatalf("delete: %v", err)
	}

	remaining, err := ts.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(remaining) != 0 {
		t.Errorf("expected 0 tasks after delete, got %d", len(remaining))
	}
}

func TestIntegration_Delete_NotFound(t *testing.T) {
	ts, js := newTestStores(t)
	err := run(t, []string{"delete", "--force", "999"}, ts, js)
	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "not found")
	}
}

// ---------------------------------------------------------------------------
// status command
// ---------------------------------------------------------------------------

func TestIntegration_Status_SetExplicit(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "Work item"}, ts, js)
	tasks, _ := ts.List()
	id := strconv.FormatInt(tasks[0].ID, 10)

	if err := run(t, []string{"status", id, "active"}, ts, js); err != nil {
		t.Fatalf("status: %v", err)
	}

	got, _ := ts.GetByID(tasks[0].ID)
	if got.Status != task.InProgress {
		t.Errorf("Status = %v, want InProgress", got.Status)
	}
}

func TestIntegration_Status_Cycle(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "Cycle me"}, ts, js)
	tasks, _ := ts.List()
	id := strconv.FormatInt(tasks[0].ID, 10)

	// Pending → InProgress
	run(t, []string{"status", id}, ts, js)
	got, _ := ts.GetByID(tasks[0].ID)
	if got.Status != task.InProgress {
		t.Errorf("after 1st cycle: Status = %v, want InProgress", got.Status)
	}

	// InProgress → Done
	run(t, []string{"status", id}, ts, js)
	got, _ = ts.GetByID(tasks[0].ID)
	if got.Status != task.Done {
		t.Errorf("after 2nd cycle: Status = %v, want Done", got.Status)
	}
}

func TestIntegration_Status_InvalidValue(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "Task"}, ts, js)
	tasks, _ := ts.List()

	err := run(t, []string{"status", strconv.FormatInt(tasks[0].ID, 10), "flying"}, ts, js)
	if err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}
}

// ---------------------------------------------------------------------------
// done command (extended)
// ---------------------------------------------------------------------------

func TestIntegration_Done_MultipleIDs(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "Task A"}, ts, js)
	run(t, []string{"add", "Task B"}, ts, js)
	tasks, _ := ts.List()

	ids := make([]string, len(tasks))
	for i, tk := range tasks {
		ids[i] = strconv.FormatInt(tk.ID, 10)
	}

	args := append([]string{"done"}, ids...)
	if err := run(t, args, ts, js); err != nil {
		t.Fatalf("done multiple: %v", err)
	}

	remaining, _ := ts.List()
	for _, tk := range remaining {
		if tk.Status != task.Done {
			t.Errorf("task #%d Status = %v, want Done", tk.ID, tk.Status)
		}
	}
}

func TestIntegration_Done_RecurringSpawnsNext(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "--recur", "daily", "Daily standup"}, ts, js)
	tasks, _ := ts.List()
	id := strconv.FormatInt(tasks[0].ID, 10)

	if err := run(t, []string{"done", id}, ts, js); err != nil {
		t.Fatalf("done recurring: %v", err)
	}

	all, _ := ts.List()
	if len(all) != 2 {
		t.Fatalf("expected 2 tasks after completing recurring (original + next), got %d", len(all))
	}
}

// ---------------------------------------------------------------------------
// add command (extended)
// ---------------------------------------------------------------------------

func TestIntegration_Add_WithDesc(t *testing.T) {
	ts, js := newTestStores(t)
	if err := run(t, []string{"add", "--desc", "Some description", "Task with desc"}, ts, js); err != nil {
		t.Fatalf("add: %v", err)
	}
	tasks, _ := ts.List()
	if tasks[0].Description != "Some description" {
		t.Errorf("Description = %q, want %q", tasks[0].Description, "Some description")
	}
}

func TestIntegration_Add_WithRecur(t *testing.T) {
	ts, js := newTestStores(t)
	if err := run(t, []string{"add", "--recur", "weekly", "Weekly review"}, ts, js); err != nil {
		t.Fatalf("add: %v", err)
	}
	tasks, _ := ts.List()
	if tasks[0].RecurFreq != task.RecurWeekly {
		t.Errorf("RecurFreq = %v, want RecurWeekly", tasks[0].RecurFreq)
	}
}

func TestIntegration_Add_InvalidRecur(t *testing.T) {
	ts, js := newTestStores(t)
	err := run(t, []string{"add", "--recur", "hourly", "Bad recur"}, ts, js)
	if err == nil {
		t.Fatal("expected error for invalid recurrence, got nil")
	}
}

// ---------------------------------------------------------------------------
// list command (extended filters)
// ---------------------------------------------------------------------------

func TestIntegration_List_FilterByPriority(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "--priority", "urgent", "Urgent task"}, ts, js)
	run(t, []string{"add", "--priority", "low", "Low task"}, ts, js)

	out := captureStdout(t, func() {
		run(t, []string{"list", "--priority", "urgent"}, ts, js)
	})

	if !strings.Contains(out, "Urgent task") {
		t.Errorf("output missing urgent task:\n%s", out)
	}
	if strings.Contains(out, "Low task") {
		t.Errorf("output should not contain low priority task:\n%s", out)
	}
}

func TestIntegration_List_FilterByTag(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "--tags", "work,golang", "Work task"}, ts, js)
	run(t, []string{"add", "--tags", "personal", "Personal task"}, ts, js)

	out := captureStdout(t, func() {
		run(t, []string{"list", "--tag", "work"}, ts, js)
	})

	if !strings.Contains(out, "Work task") {
		t.Errorf("output missing tagged task:\n%s", out)
	}
	if strings.Contains(out, "Personal task") {
		t.Errorf("output should not contain untagged task:\n%s", out)
	}
}

func TestIntegration_List_Search(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "Buy groceries"}, ts, js)
	run(t, []string{"add", "Write tests"}, ts, js)

	out := captureStdout(t, func() {
		run(t, []string{"list", "--search", "groceries"}, ts, js)
	})

	if !strings.Contains(out, "Buy groceries") {
		t.Errorf("output missing matched task:\n%s", out)
	}
	if strings.Contains(out, "Write tests") {
		t.Errorf("output should not contain unmatched task:\n%s", out)
	}
}

func TestIntegration_List_Limit(t *testing.T) {
	ts, js := newTestStores(t)
	for i := 0; i < 5; i++ {
		run(t, []string{"add", "Task"}, ts, js)
	}

	out := captureStdout(t, func() {
		run(t, []string{"list", "--limit", "2"}, ts, js)
	})

	// Count occurrences of "Task" in output (minus header line)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	dataLines := 0
	for _, line := range lines {
		if strings.Contains(line, "Task") && !strings.Contains(line, "TITLE") {
			dataLines++
		}
	}
	if dataLines > 2 {
		t.Errorf("expected at most 2 data rows, got %d:\n%s", dataLines, out)
	}
}

func TestIntegration_List_SortByPriority(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "--priority", "low", "Low"}, ts, js)
	run(t, []string{"add", "--priority", "urgent", "Urgent"}, ts, js)

	out := captureStdout(t, func() {
		run(t, []string{"list", "--sort", "priority"}, ts, js)
	})

	urgentIdx := strings.Index(out, "Urgent")
	lowIdx := strings.Index(out, "Low")
	if urgentIdx < 0 || lowIdx < 0 {
		t.Fatalf("output missing expected tasks:\n%s", out)
	}
	if urgentIdx > lowIdx {
		t.Errorf("expected Urgent before Low when sorted by priority:\n%s", out)
	}
}

func TestIntegration_List_JSON_AllFields(t *testing.T) {
	ts, js := newTestStores(t)
	run(t, []string{"add", "--desc", "desc", "--tags", "a,b", "Full task"}, ts, js)

	out := captureStdout(t, func() {
		run(t, []string{"list", "--format", "json"}, ts, js)
	})

	var result []map[string]json.RawMessage
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON: %v\noutput: %s", err, out)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 task in JSON output, got %d", len(result))
	}
	for _, field := range []string{"description", "subtasks", "time_logs", "blocked_by", "created_at", "updated_at"} {
		if _, ok := result[0][field]; !ok {
			t.Errorf("JSON output missing field %q", field)
		}
	}
}

// ---------------------------------------------------------------------------
// dispatch (Run)
// ---------------------------------------------------------------------------

func TestIntegration_Run_Dispatch(t *testing.T) {
	ts, js := newTestStores(t)

	if err := Run([]string{"add", "My task"}, ts, js, nil, config.Config{}); err != nil {
		t.Fatalf("Run add: %v", err)
	}

	tasks, err := ts.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Title != "My task" {
		t.Errorf("Title = %q, want %q", tasks[0].Title, "My task")
	}
}
