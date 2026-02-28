package focus

import (
	"database/sql"
	"fmt"
	"time"
)

// Store handles focus session persistence in SQLite.
type Store struct {
	db *sql.DB
}

// NewStore creates a focus store using the provided database connection.
// The caller is responsible for opening and closing the DB.
func NewStore(db *sql.DB) (*Store, error) {
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS focus_sessions (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			task_id      INTEGER NOT NULL DEFAULT 0,
			duration     INTEGER NOT NULL,
			started_at   TEXT NOT NULL,
			completed_at TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_focus_sessions_task ON focus_sessions(task_id)`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("focus migrate: %w", err)
		}
	}
	if err := addColumnIfNotExists(db, "focus_sessions", "kind", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return fmt.Errorf("focus migrate kind: %w", err)
	}
	if err := addColumnIfNotExists(db, "focus_sessions", "cycle_pos", "INTEGER NOT NULL DEFAULT 0"); err != nil {
		return fmt.Errorf("focus migrate cycle_pos: %w", err)
	}
	return nil
}

func addColumnIfNotExists(db *sql.DB, table, column, colDef string) error {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull int
		var dflt sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return err
		}
		if name == column {
			return nil // column already exists
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, colDef))
	return err
}

// Create inserts a new session and sets its ID.
func (s *Store) Create(session *Session) error {
	var completedAt *string
	if session.CompletedAt != nil {
		v := session.CompletedAt.UTC().Format(time.RFC3339)
		completedAt = &v
	}
	res, err := s.db.Exec(
		`INSERT INTO focus_sessions (task_id, duration, started_at, completed_at, kind, cycle_pos) VALUES (?,?,?,?,?,?)`,
		session.TaskID,
		int64(session.Duration),
		session.StartedAt.UTC().Format(time.RFC3339),
		completedAt,
		int(session.Kind),
		session.CyclePos,
	)
	if err != nil {
		return fmt.Errorf("create focus session: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}
	session.ID = id
	return nil
}

// Complete marks a session as completed by setting completed_at to now.
func (s *Store) Complete(id int64) error {
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := s.db.Exec(`UPDATE focus_sessions SET completed_at = ? WHERE id = ?`, now, id)
	if err != nil {
		return fmt.Errorf("complete focus session: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("focus session %d not found", id)
	}
	return nil
}

// ListByTask returns sessions for a given task, ordered by started_at DESC.
func (s *Store) ListByTask(taskID int64) ([]Session, error) {
	rows, err := s.db.Query(
		`SELECT id, task_id, duration, started_at, completed_at, kind, cycle_pos FROM focus_sessions WHERE task_id = ? ORDER BY started_at DESC`,
		taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("list focus sessions: %w", err)
	}
	defer rows.Close()

	var sessions []Session
	for rows.Next() {
		sess, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}
	return sessions, rows.Err()
}

// CompletionsByDay returns the count of completed sessions per day for the
// last N days, keyed by "YYYY-MM-DD".
func (s *Store) CompletionsByDay(days int) (map[string]int, error) {
	cutoff := time.Now().UTC().AddDate(0, 0, -days).Format(time.RFC3339)
	rows, err := s.db.Query(
		`SELECT DATE(completed_at) AS day, COUNT(*) FROM focus_sessions
		 WHERE completed_at IS NOT NULL AND completed_at >= ?
		 GROUP BY day`,
		cutoff,
	)
	if err != nil {
		return nil, fmt.Errorf("completions by day: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var day string
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		result[day] = count
	}
	return result, rows.Err()
}

// TodayCount returns the number of sessions completed today.
func (s *Store) TodayCount() (int, error) {
	today := time.Now().UTC().Format(time.DateOnly)
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM focus_sessions WHERE completed_at IS NOT NULL AND DATE(completed_at) = ?`,
		today,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("today count: %w", err)
	}
	return count, nil
}

// TodayWorkCount returns the number of completed work sessions today.
func (s *Store) TodayWorkCount() (int, error) {
	today := time.Now().UTC().Format(time.DateOnly)
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM focus_sessions WHERE completed_at IS NOT NULL AND kind = 0 AND DATE(completed_at) = ?`,
		today,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("today work count: %w", err)
	}
	return count, nil
}

// WeeklySummary returns the count of completed work sessions per day for the
// last 7 days, keyed by "YYYY-MM-DD".
func (s *Store) WeeklySummary() (map[string]int, error) {
	cutoff := time.Now().UTC().AddDate(0, 0, -7).Format(time.RFC3339)
	rows, err := s.db.Query(
		`SELECT DATE(completed_at) AS day, COUNT(*) FROM focus_sessions
		 WHERE completed_at IS NOT NULL AND kind = 0 AND completed_at >= ?
		 GROUP BY day`,
		cutoff,
	)
	if err != nil {
		return nil, fmt.Errorf("weekly summary: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var day string
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		result[day] = count
	}
	return result, rows.Err()
}

// Streak returns the number of consecutive days (walking back from today)
// that have at least one completed work session.
func (s *Store) Streak() (int, error) {
	rows, err := s.db.Query(
		`SELECT DISTINCT DATE(completed_at) AS day FROM focus_sessions
		 WHERE completed_at IS NOT NULL AND kind = 0
		 ORDER BY day DESC`,
	)
	if err != nil {
		return 0, fmt.Errorf("streak: %w", err)
	}
	defer rows.Close()

	var days []string
	for rows.Next() {
		var day string
		if err := rows.Scan(&day); err != nil {
			return 0, err
		}
		days = append(days, day)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	if len(days) == 0 {
		return 0, nil
	}

	streak := 0
	expected := time.Now().UTC().Format(time.DateOnly)
	for _, day := range days {
		if day == expected {
			streak++
			t, _ := time.Parse(time.DateOnly, expected)
			expected = t.AddDate(0, 0, -1).Format(time.DateOnly)
		} else {
			break
		}
	}
	return streak, nil
}

// TotalMinutesFocused returns the total minutes spent in completed work
// sessions over the last N days.
func (s *Store) TotalMinutesFocused(days int) (int, error) {
	cutoff := time.Now().UTC().AddDate(0, 0, -days).Format(time.RFC3339)
	var totalNs int64
	err := s.db.QueryRow(
		`SELECT COALESCE(SUM(duration), 0) FROM focus_sessions
		 WHERE completed_at IS NOT NULL AND kind = 0 AND completed_at >= ?`,
		cutoff,
	).Scan(&totalNs)
	if err != nil {
		return 0, fmt.Errorf("total minutes focused: %w", err)
	}
	return int(time.Duration(totalNs).Minutes()), nil
}

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanSession(s scanner) (Session, error) {
	var sess Session
	var durationNs int64
	var startedAt string
	var completedAt sql.NullString
	var kind int

	if err := s.Scan(&sess.ID, &sess.TaskID, &durationNs, &startedAt, &completedAt, &kind, &sess.CyclePos); err != nil {
		return Session{}, err
	}

	sess.Duration = time.Duration(durationNs)
	sess.Kind = SessionKind(kind)

	t, err := time.Parse(time.RFC3339, startedAt)
	if err != nil {
		return Session{}, fmt.Errorf("parse started_at %q: %w", startedAt, err)
	}
	sess.StartedAt = t

	if completedAt.Valid {
		t, err := time.Parse(time.RFC3339, completedAt.String)
		if err != nil {
			return Session{}, fmt.Errorf("parse completed_at %q: %w", completedAt.String, err)
		}
		sess.CompletedAt = &t
	}
	return sess, nil
}
