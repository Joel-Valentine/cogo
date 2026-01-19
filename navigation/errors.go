package navigation

import "errors"

var (
	ErrGoBack           = errors.New("user requested to go back")
	ErrCancel           = errors.New("user canceled operation")
	ErrEmptyState       = errors.New("no resources available")
	ErrValidationFailed = errors.New("validation failed")
	ErrNoSteps          = errors.New("flow has no steps")
	ErrInvalidStepIndex = errors.New("invalid step index")
	ErrCannotGoBack     = errors.New("cannot go back from first step")
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return ErrValidationFailed
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

