package store

import "errors"

// Store errors
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrForbidden     = errors.New("forbidden")
	ErrIsBucket      = errors.New("is a bucket")
)
