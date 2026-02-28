package cli

import (
	"errors"
	"fmt"
)

// NotFoundError indicates that a requested resource was not found.
type NotFoundError struct {
	Type string // "task", "subtask", "entry", "note"
	ID   int64
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s #%d not found", e.Type, e.ID)
}

// IsNotFound reports whether err is a *NotFoundError.
func IsNotFound(err error) bool {
	var nfe *NotFoundError
	return errors.As(err, &nfe)
}
