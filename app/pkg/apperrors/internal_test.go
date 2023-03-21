package apperrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const scope = "test"

func TestWrap(t *testing.T) {
	err := errors.New("error")

	internal := NewInternal(scope, err)
	internal = internal.(InternalError).Wrap(scope)

	require.ErrorIs(t, err, errors.Unwrap(internal))
}

func TestUnwrap(t *testing.T) {
	err := errors.New("error")

	internal := NewInternal(scope, err)

	require.ErrorIs(t, err, errors.Unwrap(internal))
}

func TestCause(t *testing.T) {
	err := errors.New("error")

	internal := NewInternal(scope, err)

	require.ErrorIs(t, err, internal.(InternalError).Cause())
}
