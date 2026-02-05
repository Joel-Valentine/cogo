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

type navigator struct{}

func NewNavigator() Navigator {
	return &navigator{}
}

func (n *navigator) Run(ctx context.Context, flow Flow) (Result, error) {
	if flow == nil {
		return Result{}, fmt.Errorf("flow cannot be nil")
	}

	steps := flow.Steps()
	if len(steps) == 0 {
		return Result{}, ErrNoSteps
	}

	state := flow.State()

	for state.CurrentStep() < len(steps) {
		select {
		case <-ctx.Done():
			return Result{}, ctx.Err()
		default:
		}

		step := steps[state.CurrentStep()]
		result, err := step.Execute(ctx, state)

		if err != nil {
			if err == ErrGoBack {
				if !state.CanGoBack() {
					return Result{}, ErrCancel
				}
				_, _ = state.Back()
				continue
			}

			if err == ErrCancel {
				return Result{}, ErrCancel
			}

			if err == context.Canceled {
				return Result{}, context.Canceled
			}

			if err == ErrEmptyState {
				return Result{}, ErrEmptyState
			}

			return Result{}, fmt.Errorf("step %q failed: %w", step.Name(), err)
		}

		if err := step.Validate(result); err != nil {
			fmt.Printf("✗ Validation error: %v\n\n", err)
			continue
		}

		state.AddResult(step.Name(), result)
	}

	history := state.History()
	if len(history) == 0 {
		return Result{}, fmt.Errorf("flow completed but no results recorded")
	}

	return history[len(history)-1].Result, nil
}

func (n *navigator) RunStep(ctx context.Context, step Step) (Result, error) {
	if step == nil {
		return Result{}, fmt.Errorf("step cannot be nil")
	}

	state := NewState()

	result, err := step.Execute(ctx, state)
	if err != nil {
		return Result{}, err
	}

	if err := step.Validate(result); err != nil {
		return Result{}, err
	}

	return result, nil
}

