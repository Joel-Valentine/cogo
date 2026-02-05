package navigation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState_NewState(t *testing.T) {
	state := NewState()

	assert.NotNil(t, state)
	assert.Equal(t, 0, state.CurrentStep())
	assert.Empty(t, state.History())
	assert.False(t, state.CanGoBack())
}

func TestState_AddResult(t *testing.T) {
	state := NewState()

	// Add first result
	result1 := NewResult("nyc3")
	state.AddResult("select_region", result1)

	assert.Equal(t, 1, state.CurrentStep())
	assert.Len(t, state.History(), 1)

	// Add second result
	result2 := NewResult("s-1vcpu-1gb")
	state.AddResult("select_size", result2)

	assert.Equal(t, 2, state.CurrentStep())
	assert.Len(t, state.History(), 2)
}

func TestState_GetResult(t *testing.T) {
	state := NewState()

	// Add results
	result1 := NewResult("nyc3")
	state.AddResult("select_region", result1)

	result2 := NewResult("s-1vcpu-1gb")
	state.AddResult("select_size", result2)

	// Get existing result
	got, found := state.GetResult("select_region")
	assert.True(t, found)
	assert.Equal(t, "nyc3", got.Value)

	// Get non-existent result
	_, found = state.GetResult("nonexistent")
	assert.False(t, found)
}

func TestState_Back(t *testing.T) {
	state := NewState()

	// Cannot go back from initial state
	_, err := state.Back()
	assert.ErrorIs(t, err, ErrCannotGoBack)
	assert.False(t, state.CanGoBack())

	// Add some results
	state.AddResult("step1", NewResult("value1"))
	state.AddResult("step2", NewResult("value2"))
	assert.Equal(t, 2, state.CurrentStep())

	// Go back
	assert.True(t, state.CanGoBack())
	index, err := state.Back()
	require.NoError(t, err)
	assert.Equal(t, 1, index)
	assert.Equal(t, 1, state.CurrentStep())

	// Go back again
	assert.True(t, state.CanGoBack())
	index, err = state.Back()
	require.NoError(t, err)
	assert.Equal(t, 0, index)
	assert.Equal(t, 0, state.CurrentStep())

	// Cannot go back further
	assert.False(t, state.CanGoBack())
	_, err = state.Back()
	assert.ErrorIs(t, err, ErrCannotGoBack)
}

func TestState_BackAndModify(t *testing.T) {
	state := NewState()

	// Add initial results
	state.AddResult("step1", NewResult("value1"))
	state.AddResult("step2", NewResult("value2"))
	state.AddResult("step3", NewResult("value3"))

	// Go back two steps
	_, _ = state.Back()
	_, _ = state.Back()
	assert.Equal(t, 1, state.CurrentStep())

	// Modify step2's result
	state.AddResult("step2", NewResult("new_value2"))

	// History should be updated, not appended
	history := state.History()
	assert.Len(t, history, 2)
	assert.Equal(t, "new_value2", history[1].Result.Value)

	// Current step should advance
	assert.Equal(t, 2, state.CurrentStep())
}

func TestState_Clear(t *testing.T) {
	state := NewState()

	// Add results
	state.AddResult("step1", NewResult("value1"))
	state.AddResult("step2", NewResult("value2"))

	// Clear state
	state.Clear()

	assert.Equal(t, 0, state.CurrentStep())
	assert.Empty(t, state.History())
	assert.False(t, state.CanGoBack())
}

func TestState_SetCurrentStep(t *testing.T) {
	state := NewState()

	// Set valid step
	err := state.SetCurrentStep(5)
	require.NoError(t, err)
	assert.Equal(t, 5, state.CurrentStep())

	// Set invalid step (negative)
	err = state.SetCurrentStep(-1)
	assert.ErrorIs(t, err, ErrInvalidStepIndex)

	// With max steps
	state = NewStateWithMax(3)
	err = state.SetCurrentStep(2)
	require.NoError(t, err)

	// Exceeds max
	err = state.SetCurrentStep(3)
	assert.ErrorIs(t, err, ErrInvalidStepIndex)
}

func TestState_History_IsCopy(t *testing.T) {
	state := NewState()

	state.AddResult("step1", NewResult("value1"))

	// Get history
	history := state.History()

	// Modify the returned history
	history[0].Result.Value = "modified"

	// Original history should be unchanged
	originalHistory := state.History()
	assert.Equal(t, "value1", originalHistory[0].Result.Value)
}

