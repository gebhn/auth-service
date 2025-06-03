package store

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrNotImplemented = errors.New("not implemented")
)
