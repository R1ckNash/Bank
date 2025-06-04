package models

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrInternal      = errors.New("internal server error")
)
