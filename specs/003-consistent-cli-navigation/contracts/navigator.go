// Package contracts defines the core interfaces for the navigation framework.
// This is a contract definition file - it specifies the API but is not compiled.
//
// These interfaces enable consistent CLI navigation across all cloud providers
// and commands, following industry best practices from research on 11 major CLI tools.
package contracts

import (
	"context"
)

// Navigator orchestrates multi-step interactive CLI flows with support for
// back navigation, cancellation, and error recovery.
//
// Design based on research findings:
// - gcloud-style back navigation (only modern CLI with this feature)
// - Ctrl+C/Esc cancellation (universal pattern, 11/11 tools)
// - Empty state handling (10/11 tools return exit 0)
// - Validate-on-Enter (9/10 interactive tools)
type Navigator interface {
	// Run executes a Flow and returns the final result.
	// Returns error only for fatal errors (not user cancellation).
	//
	// Cancellation behaviors:
	// - Ctrl+C: Returns context.Canceled error
	// - 'q' or Esc: Returns ErrFlowCanceled
	// - Empty state: Returns ErrEmptyState (not fatal)
	Run(ctx context.Context, flow Flow) (Result, error)

	// RunStep executes a single Step and returns the result.
	// Used for testing individual steps or simple single-step operations.
	RunStep(ctx context.Context, step Step) (Result, error)
}

// Flow represents a multi-step interactive sequence with history and state management.
//
// Example: Creating a droplet = Flow with steps [SelectRegion, SelectSize, SelectImage, Confirm]
type Flow interface {
	// Name returns the flow's display name (e.g., "Create Droplet").
	Name() string

	// Steps returns all steps in the flow, in order.
	// Steps are executed sequentially unless user navigates back.
	Steps() []Step

	// State returns the current flow state (history, current step, selections).
	State() State
}

// Step represents a single interactive prompt or action in a Flow.
//
// Design principles from research:
// - Show defaults (gh, npm, gcloud pattern)
// - Validate on Enter only (9/10 tools)
// - Support navigation keys (back, quit, cancel)
type Step interface {
	// Name returns the step's identifier (e.g., "select_region").
	Name() string

	// Prompt returns the text shown to the user.
	// Should be short and clear (e.g., "Select region:").
	Prompt() string

	// Execute runs the step's interaction (prompt, validate, return result).
	// Returns:
	// - Result: User's selection or input
	// - error: Fatal error only (validation errors re-prompt)
	//
	// Special errors:
	// - ErrGoBack: User pressed 'b' or ←
	// - ErrCancel: User pressed 'q' or Esc
	// - context.Canceled: User pressed Ctrl+C
	Execute(ctx context.Context, state State) (Result, error)

	// Validate checks if a result is valid for this step.
	// Called AFTER user presses Enter (not per-keystroke).
	// Returns error message if invalid, empty string if valid.
	Validate(result Result) error

	// Default returns the default value for this step (if any).
	// Shown in prompt: "? Prompt (default):"
	Default() interface{}
}

// State manages the flow's execution state, including history for back navigation.
//
// Inspired by git's rebase navigation and gcloud init's back functionality.
type State interface {
	// History returns all completed steps and their results.
	// Ordered from oldest to newest: [step1, step2, step3]
	History() []StepResult

	// CurrentStep returns the index of the current step (0-based).
	CurrentStep() int

	// SetCurrentStep moves to a specific step index.
	// Used for back navigation: SetCurrentStep(CurrentStep() - 1)
	SetCurrentStep(index int) error

	// AddResult records a step's result in history.
	AddResult(step Step, result Result)

	// GetResult retrieves the result for a specific step by name.
	// Returns (result, found bool).
	GetResult(stepName string) (Result, bool)

	// Clear resets all history and state (used when starting new flow).
	Clear()

	// CanGoBack returns true if there are previous steps to return to.
	CanGoBack() bool

	// Back moves to the previous step and returns its index.
	// Returns error if already at first step.
	Back() (int, error)
}

// Result represents the outcome of a Step's execution.
//
// Generic container for:
// - Selected item from list (droplet, region, size)
// - User input (name, quantity)
// - Confirmation (yes/no)
type Result struct {
	// Value is the actual result data (string, int, struct, etc.).
	Value interface{}

	// Metadata stores additional context (e.g., selected index, display name).
	Metadata map[string]interface{}
}

// StepResult pairs a Step with its Result for history tracking.
type StepResult struct {
	Step   Step
	Result Result
}

// Sentinel errors for navigation control flow.
// These are NOT fatal errors - they indicate user actions.
var (
	// ErrGoBack indicates user wants to return to previous step.
	// Triggered by: 'b' key, ← arrow, or custom back option.
	ErrGoBack error

	// ErrCancel indicates user wants to quit the entire flow.
	// Triggered by: 'q' key or Esc.
	ErrCancel error

	// ErrEmptyState indicates an operation has no resources to work with.
	// Example: "No droplets found" when trying to destroy.
	// NOT a fatal error - should display helpful message and exit 0.
	ErrEmptyState error

	// ErrValidationFailed indicates user input failed validation.
	// Should re-prompt, not exit. Contains the validation message.
	ErrValidationFailed error
)

