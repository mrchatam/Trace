package store

import (
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound is the sentinel every not-found store lookup wraps.
var ErrNotFound = errors.New("store: entity not found")

// NotFoundError reports a missing row and names the entity type that was
// actually looked up, so a caller can say "no task" instead of guessing from
// the function it happened to call through (#155).
type NotFoundError struct {
	Entity string
	ID     string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("store: %s %q: not found", e.Entity, e.ID)
}

// Unwrap exposes the sentinel for errors.Is(err, ErrNotFound).
func (e *NotFoundError) Unwrap() error { return ErrNotFound }

// Is keeps errors.Is(err, sql.ErrNoRows) true for every existing caller that
// predates ErrNotFound (domain/scope.go, retrieval/exact.go, httpapi).
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound || target == sql.ErrNoRows
}

// notFound builds the canonical missing-row error for an entity lookup.
func notFound(entity, id string) error {
	return &NotFoundError{Entity: entity, ID: id}
}

// IsNotFound reports whether err is or wraps a missing-row error.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	var nf *NotFoundError
	if errors.As(err, &nf) {
		return true
	}
	return errors.Is(err, ErrNotFound) || errors.Is(err, sql.ErrNoRows)
}
