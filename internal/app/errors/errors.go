package errors

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrRowDeleted = errors.New("row deleted")
)
