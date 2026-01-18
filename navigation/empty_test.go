package navigation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{"nil", nil, true},
		{"empty string slice", []string{}, true},
		{"non-empty string slice", []string{"item"}, false},
		{"empty int slice", []int{}, true},
		{"non-empty int slice", []int{1}, false},
		{"empty interface slice", []interface{}{}, true},
		{"non-empty interface slice", []interface{}{1}, false},
		{"empty map", map[string]interface{}{}, true},
		{"non-empty map", map[string]interface{}{"key": "val"}, false},
		{"empty string", "", true},
		{"non-empty string", "test", false},
		{"other type", 42, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsEmpty(tt.value))
		})
	}
}

func TestEmptyStateHandler_Check(t *testing.T) {
	handler := &EmptyStateHandler{
		ResourceName: "droplets",
	}

	// Empty
	err := handler.Check([]string{})
	assert.ErrorIs(t, err, ErrEmptyState)

	// Not empty
	err = handler.Check([]string{"item"})
	assert.NoError(t, err)
}

func TestEmptyStateHandler_CheckAndDisplay(t *testing.T) {
	handler := &EmptyStateHandler{
		ResourceName:     "droplets",
		Context:          "in your account",
		SuggestedCommand: "cogo create",
	}

	// Empty - should display message
	err := handler.CheckAndDisplay([]string{})
	assert.ErrorIs(t, err, ErrEmptyState)

	// Not empty - should not display
	err = handler.CheckAndDisplay([]string{"item"})
	assert.NoError(t, err)
}

func TestEmptyStateMessage(t *testing.T) {
	msg := EmptyStateMessage("droplets", "in your account", "Run 'cogo create'")

	assert.Contains(t, msg, "No droplets found")
	assert.Contains(t, msg, "in your account")
	assert.Contains(t, msg, "Run 'cogo create'")
}

func TestEmptyStateError(t *testing.T) {
	err := EmptyStateError("droplets", "in your account", "Run 'cogo create'")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEmptyState)
	assert.Contains(t, err.Error(), "No droplets found")
}

func TestErrorMessage(t *testing.T) {
	msg := ErrorMessage("Failed to create", "API returned 422", "Try a different region")

	assert.Contains(t, msg, "Error:")
	assert.Contains(t, msg, "Failed to create")
	assert.Contains(t, msg, "API returned 422")
	assert.Contains(t, msg, "Try a different region")
}

func TestErrorMessage_Minimal(t *testing.T) {
	msg := ErrorMessage("Failed to create", "", "")

	assert.Contains(t, msg, "Error:")
	assert.Contains(t, msg, "Failed to create")
	assert.NotContains(t, msg, "\n\n")
}

