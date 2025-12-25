package storage

import "errors"

var (
	// ErrNotFound indicates that a requested resource was not found.
	ErrNotFound = errors.New("resource not found")

	// ErrDuplicateEntry indicates that a duplicate entry was attempted to be inserted.
	ErrDuplicateEntry = errors.New("duplicate entry")

	// ErrInvalidInput indicates that the input provided is invalid.
	ErrInvalidInput = errors.New("invalid input")
)
