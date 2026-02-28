package focus

import (
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	// Enable foreign keys for consistency with production.
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestNewStore(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestCreateAndListByTask(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	now := time.Now().Truncate(time.Second)
	sess := &Session{
		TaskID:    42,
		Duration:  DefaultDuration,
		StartedAt: now,
	}

	if err := store.Create(sess); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if sess.ID == 0 {
		t.Error("expected session ID to be set after Create")
	}

	sessions, err := store.ListByTask(42)
	if err != nil {
		t.Fatalf("ListByTask: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	got := sessions[0]
	if got.ID != sess.ID {
		t.Errorf("ID = %d, want %d", got.ID, sess.ID)
	}
	if got.TaskID != 42 {
		t.Errorf("TaskID = %d, want 42", got.TaskID)
	}
	if got.Duration != DefaultDuration {
		t.Errorf("Duration = %v, want %v", got.Duration, DefaultDuration)
	}
	if got.IsCompleted() {
		t.Error("expected incomplete session")
	}
}

func TestComplete(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	sess := &Session{
		TaskID:    1,
		Duration:  DefaultDuration,
		StartedAt: time.Now().Truncate(time.Second),
	}
	if err := store.Create(sess); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := store.Complete(sess.ID); err != nil {
		t.Fatalf("Complete: %v", err)
	}

	sessions, err := store.ListByTask(1)
	if err != nil {
		t.Fatalf("ListByTask: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if !sessions[0].IsCompleted() {
		t.Error("expected session to be completed")
	}
}

func TestCompleteNotFound(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	if err := store.Complete(999); err == nil {
		t.Error("expected error completing non-existent session")
	}
}

func TestListByTaskEmpty(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	sessions, err := store.ListByTask(999)
	if err != nil {
		t.Fatalf("ListByTask: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions, got %d", len(sessions))
	}
}

func TestListByTaskOrdering(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	base := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		sess := &Session{
			TaskID:    5,
			Duration:  DefaultDuration,
			StartedAt: base.Add(time.Duration(i) * time.Hour),
		}
		if err := store.Create(sess); err != nil {
			t.Fatalf("Create #%d: %v", i, err)
		}
	}

	sessions, err := store.ListByTask(5)
	if err != nil {
		t.Fatalf("ListByTask: %v", err)
	}
	if len(sessions) != 3 {
		t.Fatalf("expected 3 sessions, got %d", len(sessions))
	}
	// Should be ordered DESC by started_at.
	if !sessions[0].StartedAt.After(sessions[1].StartedAt) {
		t.Error("sessions not ordered by started_at DESC")
	}
	if !sessions[1].StartedAt.After(sessions[2].StartedAt) {
		t.Error("sessions not ordered by started_at DESC")
	}
}

func TestTodayCount(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	count, err := store.TodayCount()
	if err != nil {
		t.Fatalf("TodayCount: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}

	// Create and complete a session today.
	now := time.Now().Truncate(time.Second)
	sess := &Session{
		TaskID:    1,
		Duration:  DefaultDuration,
		StartedAt: now.Add(-30 * time.Minute),
	}
	if err := store.Create(sess); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := store.Complete(sess.ID); err != nil {
		t.Fatalf("Complete: %v", err)
	}

	count, err = store.TodayCount()
	if err != nil {
		t.Fatalf("TodayCount: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestCompletionsByDay(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	now := time.Now().Truncate(time.Second)

	// Create 2 completed sessions today.
	for i := 0; i < 2; i++ {
		sess := &Session{
			TaskID:    1,
			Duration:  DefaultDuration,
			StartedAt: now.Add(time.Duration(-i) * time.Hour),
		}
		if err := store.Create(sess); err != nil {
			t.Fatalf("Create: %v", err)
		}
		if err := store.Complete(sess.ID); err != nil {
			t.Fatalf("Complete: %v", err)
		}
	}

	// Create 1 incomplete session (should not be counted).
	incomplete := &Session{
		TaskID:    1,
		Duration:  DefaultDuration,
		StartedAt: now,
	}
	if err := store.Create(incomplete); err != nil {
		t.Fatalf("Create: %v", err)
	}

	result, err := store.CompletionsByDay(7)
	if err != nil {
		t.Fatalf("CompletionsByDay: %v", err)
	}

	today := now.UTC().Format(time.DateOnly)
	if result[today] != 2 {
		t.Errorf("expected 2 completions today, got %d", result[today])
	}
}

func TestTodayWorkCount(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	count, err := store.TodayWorkCount()
	if err != nil {
		t.Fatalf("TodayWorkCount: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}

	now := time.Now().Truncate(time.Second)

	// Create and complete a work session today.
	work := &Session{Duration: DefaultDuration, StartedAt: now.Add(-30 * time.Minute), Kind: KindWork}
	if err := store.Create(work); err != nil {
		t.Fatalf("Create work: %v", err)
	}
	if err := store.Complete(work.ID); err != nil {
		t.Fatalf("Complete work: %v", err)
	}

	// Create and complete a short break session today (should not count).
	brk := &Session{Duration: 5 * time.Minute, StartedAt: now.Add(-25 * time.Minute), Kind: KindShortBreak}
	if err := store.Create(brk); err != nil {
		t.Fatalf("Create break: %v", err)
	}
	if err := store.Complete(brk.ID); err != nil {
		t.Fatalf("Complete break: %v", err)
	}

	// Create an incomplete work session (should not count).
	incomplete := &Session{Duration: DefaultDuration, StartedAt: now, Kind: KindWork}
	if err := store.Create(incomplete); err != nil {
		t.Fatalf("Create incomplete: %v", err)
	}

	count, err = store.TodayWorkCount()
	if err != nil {
		t.Fatalf("TodayWorkCount: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
}

func TestWeeklySummary(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	now := time.Now().Truncate(time.Second)

	// Create 2 completed work sessions today.
	for i := 0; i < 2; i++ {
		sess := &Session{
			Duration:  DefaultDuration,
			StartedAt: now.Add(time.Duration(-i) * time.Hour),
			Kind:      KindWork,
		}
		if err := store.Create(sess); err != nil {
			t.Fatalf("Create #%d: %v", i, err)
		}
		if err := store.Complete(sess.ID); err != nil {
			t.Fatalf("Complete #%d: %v", i, err)
		}
	}

	// Create a completed break session (should not appear in weekly summary).
	brk := &Session{Duration: 5 * time.Minute, StartedAt: now, Kind: KindShortBreak}
	if err := store.Create(brk); err != nil {
		t.Fatalf("Create break: %v", err)
	}
	if err := store.Complete(brk.ID); err != nil {
		t.Fatalf("Complete break: %v", err)
	}

	result, err := store.WeeklySummary()
	if err != nil {
		t.Fatalf("WeeklySummary: %v", err)
	}

	today := now.UTC().Format(time.DateOnly)
	if result[today] != 2 {
		t.Errorf("expected 2 work sessions today, got %d", result[today])
	}
	// Break sessions must not appear.
	total := 0
	for _, v := range result {
		total += v
	}
	if total != 2 {
		t.Errorf("expected total 2 across all days, got %d", total)
	}
}

func TestStreak(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	// No sessions → streak = 0.
	streak, err := store.Streak()
	if err != nil {
		t.Fatalf("Streak: %v", err)
	}
	if streak != 0 {
		t.Errorf("expected streak 0, got %d", streak)
	}

	now := time.Now().Truncate(time.Second)

	// Add a completed work session today.
	today := &Session{Duration: DefaultDuration, StartedAt: now.Add(-30 * time.Minute), Kind: KindWork}
	if err := store.Create(today); err != nil {
		t.Fatalf("Create today: %v", err)
	}
	if err := store.Complete(today.ID); err != nil {
		t.Fatalf("Complete today: %v", err)
	}

	streak, err = store.Streak()
	if err != nil {
		t.Fatalf("Streak: %v", err)
	}
	if streak != 1 {
		t.Errorf("expected streak 1, got %d", streak)
	}
}

func TestTotalMinutesFocused(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	mins, err := store.TotalMinutesFocused(7)
	if err != nil {
		t.Fatalf("TotalMinutesFocused: %v", err)
	}
	if mins != 0 {
		t.Errorf("expected 0, got %d", mins)
	}

	now := time.Now().Truncate(time.Second)

	// Create and complete two 25-min work sessions.
	for i := 0; i < 2; i++ {
		sess := &Session{
			Duration:  25 * time.Minute,
			StartedAt: now.Add(time.Duration(-i) * time.Hour),
			Kind:      KindWork,
		}
		if err := store.Create(sess); err != nil {
			t.Fatalf("Create #%d: %v", i, err)
		}
		if err := store.Complete(sess.ID); err != nil {
			t.Fatalf("Complete #%d: %v", i, err)
		}
	}

	// A completed break session should not count.
	brk := &Session{Duration: 5 * time.Minute, StartedAt: now, Kind: KindShortBreak}
	if err := store.Create(brk); err != nil {
		t.Fatalf("Create break: %v", err)
	}
	if err := store.Complete(brk.ID); err != nil {
		t.Fatalf("Complete break: %v", err)
	}

	mins, err = store.TotalMinutesFocused(7)
	if err != nil {
		t.Fatalf("TotalMinutesFocused: %v", err)
	}
	if mins != 50 {
		t.Errorf("expected 50 minutes, got %d", mins)
	}
}

func TestCreateWithZeroTaskID(t *testing.T) {
	db := openTestDB(t)
	store, err := NewStore(db)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	sess := &Session{
		TaskID:    0,
		Duration:  DefaultDuration,
		StartedAt: time.Now().Truncate(time.Second),
	}
	if err := store.Create(sess); err != nil {
		t.Fatalf("Create with zero TaskID: %v", err)
	}

	sessions, err := store.ListByTask(0)
	if err != nil {
		t.Fatalf("ListByTask(0): %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].TaskID != 0 {
		t.Errorf("TaskID = %d, want 0", sessions[0].TaskID)
	}
}
