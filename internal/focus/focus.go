package focus

import (
	"fmt"
	"time"
)

// DefaultDuration is the standard focus session length (Pomodoro).
const DefaultDuration = 25 * time.Minute

// SessionKind distinguishes work sessions from break sessions.
type SessionKind int

const (
	KindWork       SessionKind = 0
	KindShortBreak SessionKind = 1
	KindLongBreak  SessionKind = 2
)

// String returns the full name of the session kind.
func (k SessionKind) String() string {
	switch k {
	case KindWork:
		return "Work"
	case KindShortBreak:
		return "Short Break"
	case KindLongBreak:
		return "Long Break"
	default:
		return "Unknown"
	}
}

// Label returns a shorter display label for the session kind.
func (k SessionKind) Label() string {
	switch k {
	case KindWork:
		return "Focus"
	case KindShortBreak:
		return "Break"
	case KindLongBreak:
		return "Long Break"
	default:
		return "Unknown"
	}
}

// Session represents a single focus/pomodoro session.
type Session struct {
	ID          int64
	TaskID      int64         // 0 if no associated task
	Duration    time.Duration // planned duration
	StartedAt   time.Time
	CompletedAt *time.Time  // nil if abandoned or in-progress
	Kind        SessionKind
	CyclePos    int // 1-4 for work sessions, 0 for breaks
}

// IsCompleted reports whether the session was completed.
func (s Session) IsCompleted() bool {
	return s.CompletedAt != nil
}

// Elapsed returns the time elapsed since StartedAt, capped at Duration.
func (s Session) Elapsed(now time.Time) time.Duration {
	elapsed := now.Sub(s.StartedAt)
	if elapsed < 0 {
		return 0
	}
	if elapsed > s.Duration {
		return s.Duration
	}
	return elapsed
}

// Remaining returns Duration minus Elapsed, with a minimum of 0.
func (s Session) Remaining(now time.Time) time.Duration {
	r := s.Duration - s.Elapsed(now)
	if r < 0 {
		return 0
	}
	return r
}

// FormatTimer formats a duration as "MM:SS" (e.g., "25:00", "04:30").
func FormatTimer(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d.Seconds())
	minutes := total / 60
	seconds := total % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
