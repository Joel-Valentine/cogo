package navigation

import (
	"context"
	"fmt"
)

// Navigator orchestrates multi-step interactive CLI flows with support for
// back navigation, cancellation, and error recovery.
//
// Research-based design:
// - gcloud-style back navigation (2/11 tools)
// - Universal Ctrl+C cancellation (11/11 tools)
// - Empty state handling (10/11 tools exit 0)
// - Validate-on-Enter (9/10 interactive tools)
type Navigator interface {
	// Run executes a Flow and returns the final result.
	// Returns error only for fatal errors (not user cancellation).
	Run(ctx context.Context, flow Flow) (Result, error)

	// RunStep executes a single Step and returns the result.
	// Used for testing individual steps or simple single-step operations.
	RunStep(ctx context.Context, step Step) (Result, error)
}

// navigator is the concrete implementation of Navigator interface.
type navigator struct {
	// Options for customizing behavior (future extensibility)
}

// NewNavigator creates a new Navigator instance.
func NewNavigator() Navigator {
	return &navigator{}
}

// Run executes a multi-step Flow with support for back navigation.
//
// Flow execution:
//  1. Execute steps sequentially from State.CurrentStep()
//  2. On success: Add result to state, move to next step
//  3. On ErrGoBack: Move to previous step, re-execute
//  4. On ErrCancel or context.Canceled: Exit cleanly
//  5. On ErrEmptyState: Return immediately (not fatal)
//  6. On other error: Return error
//
// Returns:
// - Final result from last step if completed
// - Empty result + ErrCancel if user quit
// - Empty result + ErrEmptyState if empty state encountered
// - Empty result + context.Canceled if Ctrl+C pressed
// - Error for fatal errors
func (n *navigator) Run(ctx context.Context, flow Flow) (Result, error) {
	// Validate flow
	if flow == nil {
		return Result{}, fmt.Errorf("flow cannot be nil")
	}

	steps := flow.Steps()
	if len(steps) == 0 {
		return Result{}, ErrNoSteps
	}

	state := flow.State()

	// Execute steps
	for state.CurrentStep() < len(steps) {
		// Check context cancellation before each step
		select {
		case <-ctx.Done():
			return Result{}, ctx.Err()
		default:
		}

		step := steps[state.CurrentStep()]

		// Execute step
		result, err := step.Execute(ctx, state)

		// Handle errors
		if err != nil {
			// User wants to go back
			if err == ErrGoBack {
				if !state.CanGoBack() {
					// Already at first step - treat as cancel
					return Result{}, ErrCancel
				}
				_, _ = state.Back()
				continue
			}

			// User canceled
			if err == ErrCancel {
				return Result{}, ErrCancel
			}

			// Ctrl+C
			if err == context.Canceled {
				return Result{}, context.Canceled
			}

			// Empty state (not fatal)
			if err == ErrEmptyState {
				return Result{}, ErrEmptyState
			}

			// Fatal error
			return Result{}, fmt.Errorf("step %q failed: %w", step.Name(), err)
		}

		// Validate result (optional - steps can validate in Execute)
		if err := step.Validate(result); err != nil {
			// Validation failed - should have been caught in Execute,
			// but handle gracefully if step didn't validate
			fmt.Printf("✗ Validation error: %v\n\n", err)
			continue // Re-run same step
		}

		// Add result to state and move to next step
		state.AddResult(step.Name(), result)
	}

	// All steps completed - return final result
	history := state.History()
	if len(history) == 0 {
		return Result{}, fmt.Errorf("flow completed but no results recorded")
	}

	return history[len(history)-1].Result, nil
}

// RunStep executes a single Step independently.
// Useful for:
// - Testing individual steps
// - Simple single-step commands (list, destroy one droplet, etc.)
// - Steps that don't need flow context
func (n *navigator) RunStep(ctx context.Context, step Step) (Result, error) {
	if step == nil {
		return Result{}, fmt.Errorf("step cannot be nil")
	}

	// Create temporary state for step
	state := NewState()

	// Execute step
	result, err := step.Execute(ctx, state)
	if err != nil {
		return Result{}, err
	}

	// Validate result
	if err := step.Validate(result); err != nil {
		return Result{}, err
	}

	return result, nil
}

