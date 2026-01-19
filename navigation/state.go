package navigation

import "fmt"

// State manages the flow's execution state, including history for back navigation.
//
// Inspired by git's rebase navigation and gcloud init's back functionality.
// Research finding: Only 2/11 tools (gcloud, git) support back navigation, but
// this is a key UX improvement for interactive CLIs.
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
	AddResult(stepName string, result Result)

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

type state struct {
	history     []StepResult
	currentStep int
	maxSteps    int
}

func NewState() State {
	return &state{
		history:     make([]StepResult, 0),
		currentStep: 0,
		maxSteps:    -1,
	}
}

func NewStateWithMax(maxSteps int) State {
	return &state{
		history:     make([]StepResult, 0, maxSteps),
		currentStep: 0,
		maxSteps:    maxSteps,
	}
}

func (s *state) History() []StepResult {
	historyCopy := make([]StepResult, len(s.history))
	copy(historyCopy, s.history)
	return historyCopy
}

func (s *state) CurrentStep() int {
	return s.currentStep
}

func (s *state) SetCurrentStep(index int) error {
	if index < 0 {
		return fmt.Errorf("%w: %d (must be >= 0)", ErrInvalidStepIndex, index)
	}
	if s.maxSteps > 0 && index >= s.maxSteps {
		return fmt.Errorf("%w: %d (max is %d)", ErrInvalidStepIndex, index, s.maxSteps-1)
	}
	s.currentStep = index
	return nil
}

func (s *state) AddResult(stepName string, result Result) {
	stepResult := StepResult{
		StepName: stepName,
		Result:   result,
	}

	if s.currentStep < len(s.history) {
		s.history = s.history[:s.currentStep]
	}

	s.history = append(s.history, stepResult)
	s.currentStep++
}

func (s *state) GetResult(stepName string) (Result, bool) {
	for _, sr := range s.history {
		if sr.StepName == stepName {
			return sr.Result, true
		}
	}
	return Result{}, false
}

func (s *state) Clear() {
	s.history = make([]StepResult, 0)
	s.currentStep = 0
}

func (s *state) CanGoBack() bool {
	return s.currentStep > 0
}

func (s *state) Back() (int, error) {
	if !s.CanGoBack() {
		return s.currentStep, ErrCannotGoBack
	}
	s.currentStep--
	return s.currentStep, nil
}

