package navigation

import "errors"

// Navigation control flow errors.
// These are NOT fatal errors - they indicate user actions or expected states.

var (
	// ErrGoBack indicates user wants to return to previous step.
	// Triggered by: 'b' key, ← arrow, or custom back option.
	//
	// Example: User is at "Select Size" step and presses 'b' to go back to "Select Region".
	ErrGoBack = errors.New("user requested to go back")

	// ErrCancel indicates user wants to quit the entire flow.
	// Triggered by: 'q' key or Esc.
	//
	// Example: User is in middle of create flow and presses 'q' to quit entirely.
	ErrCancel = errors.New("user canceled operation")

	// ErrEmptyState indicates an operation has no resources to work with.
	// Example: "No droplets found" when trying to destroy.
	//
	// Research finding: NOT a fatal error - should display helpful message and exit 0.
	// Used by kubectl, gh, gcloud (10/11 tools return exit 0 for empty state).
	ErrEmptyState = errors.New("no resources available")

	// ErrValidationFailed indicates user input failed validation.
	// Should re-prompt, not exit. The error message contains the validation details.
	//
	// Research finding: Validate on Enter only (9/10 interactive tools).
	ErrValidationFailed = errors.New("validation failed")

	// ErrNoSteps indicates a Flow was created with no steps.
	ErrNoSteps = errors.New("flow has no steps")

	// ErrInvalidStepIndex indicates an attempt to navigate to an invalid step.
	ErrInvalidStepIndex = errors.New("invalid step index")

	// ErrCannotGoBack indicates Back() was called when already at first step.
	ErrCannotGoBack = errors.New("cannot go back from first step")
)

// ValidationError wraps a validation failure message.
// Implements error interface and can be unwrapped to ErrValidationFailed.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return ErrValidationFailed
}

// NewValidationError creates a validation error with a message.
func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

