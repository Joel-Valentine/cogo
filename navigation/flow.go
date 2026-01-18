package navigation

import (
	"context"
	"fmt"
)

// Flow represents a multi-step interactive sequence with history and state management.
//
// Example: Creating a droplet = Flow with steps [SelectRegion, SelectSize, SelectImage, Confirm]
//
// Research-based design:
// - Inspired by gcloud init (only modern CLI with full back navigation)
// - State management from git rebase (history + undo)
// - Step-by-step pattern from terraform (plan → confirm → apply)
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
	// Used for state tracking and history.
	Name() string

	// Prompt returns the text shown to the user.
	// Should be short and clear (e.g., "Select region:").
	Prompt() string

	// Execute runs the step's interaction (prompt, validate, return result).
	// Returns:
	// - Result: User's selection or input
	// - error: Fatal error only (validation errors should re-prompt internally)
	//
	// Special errors:
	// - ErrGoBack: User pressed 'b' or ←
	// - ErrCancel: User pressed 'q' or Esc
	// - context.Canceled: User pressed Ctrl+C
	// - ErrEmptyState: No resources available (not fatal)
	Execute(ctx context.Context, state State) (Result, error)

	// Validate checks if a result is valid for this step.
	// Called AFTER user presses Enter (not per-keystroke).
	// Returns error if invalid, nil if valid.
	//
	// Note: Most steps should validate within Execute() and re-prompt.
	// This is a secondary validation for flow-level consistency checks.
	Validate(result Result) error

	// Default returns the default value for this step (if any).
	// Shown in prompt: "? Prompt (default):"
	// Return nil if no default.
	Default() interface{}
}

// flow is a concrete implementation of Flow interface.
type flow struct {
	name  string
	steps []Step
	state State
}

// NewFlow creates a new Flow with the given name and steps.
func NewFlow(name string, steps []Step) Flow {
	return &flow{
		name:  name,
		steps: steps,
		state: NewState(),
	}
}

// NewFlowWithState creates a new Flow with existing state (for resuming/testing).
func NewFlowWithState(name string, steps []Step, state State) Flow {
	return &flow{
		name:  name,
		steps: steps,
		state: state,
	}
}

func (f *flow) Name() string {
	return f.name
}

func (f *flow) Steps() []Step {
	return f.steps
}

func (f *flow) State() State {
	return f.state
}

// SimpleStep is a basic Step implementation for common use cases.
// For more complex steps, implement the Step interface directly.
type SimpleStep struct {
	StepName     string
	PromptText   string
	DefaultValue interface{}
	ExecuteFunc  func(ctx context.Context, state State) (Result, error)
	ValidateFunc func(result Result) error
}

func (s *SimpleStep) Name() string {
	return s.StepName
}

func (s *SimpleStep) Prompt() string {
	return s.PromptText
}

func (s *SimpleStep) Execute(ctx context.Context, state State) (Result, error) {
	if s.ExecuteFunc == nil {
		return Result{}, fmt.Errorf("SimpleStep %q has no ExecuteFunc", s.StepName)
	}
	return s.ExecuteFunc(ctx, state)
}

func (s *SimpleStep) Validate(result Result) error {
	if s.ValidateFunc == nil {
		return nil // No validation
	}
	return s.ValidateFunc(result)
}

func (s *SimpleStep) Default() interface{} {
	return s.DefaultValue
}

