package navigation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNavigationErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrGoBack", ErrGoBack},
		{"ErrCancel", ErrCancel},
		{"ErrEmptyState", ErrEmptyState},
		{"ErrValidationFailed", ErrValidationFailed},
		{"ErrNoSteps", ErrNoSteps},
		{"ErrInvalidStepIndex", ErrInvalidStepIndex},
		{"ErrCannotGoBack", ErrCannotGoBack},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
			assert.NotEmpty(t, tt.err.Error())
		})
	}
}

func TestValidationError(t *testing.T) {
	msg := "field is required"
	err := NewValidationError(msg)

	assert.Error(t, err)
	assert.Equal(t, msg, err.Error())

	// Should unwrap to ErrValidationFailed
	assert.ErrorIs(t, err, ErrValidationFailed)
}

func TestValidationError_Unwrap(t *testing.T) {
	err := NewValidationError("test message")

	// errors.Is should work
	assert.True(t, errors.Is(err, ErrValidationFailed))

	// errors.As should work
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Equal(t, "test message", valErr.Message)
}

