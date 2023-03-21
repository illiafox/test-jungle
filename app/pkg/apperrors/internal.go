package apperrors

import (
	"runtime"
	"strconv"
)

const Separator = ": "

type InternalError struct {
	Err   error
	Scope string
	Line  string
}

func NewInternal(scope string, err error) error {
	if err == nil {
		return nil
	}

	_, file, line, _ := runtime.Caller(1)

	return InternalError{
		Err:   err,
		Scope: scope,
		Line:  file + Separator + strconv.Itoa(line),
	}
}

func (i InternalError) Error() string {
	return i.Scope + Separator + i.Err.Error()
}

func (i InternalError) Wrap(scope string) error {
	if i.Scope != "" {
		i.Scope = i.Scope + Separator + scope
	} else {
		i.Scope = scope
	}

	return i
}

func (i InternalError) Unwrap() error {
	return i.Err
}

func (i InternalError) Cause() error {
	return i.Err
}
