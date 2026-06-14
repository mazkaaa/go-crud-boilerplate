package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource already exists")
	ErrInternal     = errors.New("internal server error")
	ErrInvalidInput = errors.New("invalid input")
)
