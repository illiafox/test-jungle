package apperrors

import "errors"

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrNotFound      = errors.New("not found")
)
