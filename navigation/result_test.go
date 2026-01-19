package navigation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResult(t *testing.T) {
	result := NewResult("test_value")

	assert.Equal(t, "test_value", result.Value)
	assert.NotNil(t, result.Metadata)
	assert.Empty(t, result.Metadata)
}

func TestNewResultWithMetadata(t *testing.T) {
	metadata := map[string]interface{}{
		"index": 0,
		"name":  "test",
	}

	result := NewResultWithMetadata("value", metadata)

	assert.Equal(t, "value", result.Value)
	assert.Equal(t, metadata, result.Metadata)
}

func TestResult_WithMetadata(t *testing.T) {
	result := NewResult("value")

	result = result.WithMetadata("key1", "val1")
	result = result.WithMetadata("key2", 42)

	val1, ok := result.GetMetadata("key1")
	assert.True(t, ok)
	assert.Equal(t, "val1", val1)

	val2, ok := result.GetMetadata("key2")
	assert.True(t, ok)
	assert.Equal(t, 42, val2)
}

func TestResult_GetMetadata(t *testing.T) {
	result := NewResult("value")

	// Get non-existent key
	_, ok := result.GetMetadata("nonexistent")
	assert.False(t, ok)

	// Add and get existing key
	result = result.WithMetadata("exists", "yes")
	val, ok := result.GetMetadata("exists")
	assert.True(t, ok)
	assert.Equal(t, "yes", val)
}

func TestResult_GetMetadata_NilMetadata(t *testing.T) {
	// Create result with nil metadata directly
	result := Result{
		Value:    "test",
		Metadata: nil,
	}

	_, ok := result.GetMetadata("key")
	assert.False(t, ok)
}

func TestResult_WithMetadata_NilMetadata(t *testing.T) {
	// Create result with nil metadata
	result := Result{
		Value:    "test",
		Metadata: nil,
	}

	// WithMetadata should initialize the map
	result = result.WithMetadata("key", "value")

	val, ok := result.GetMetadata("key")
	assert.True(t, ok)
	assert.Equal(t, "value", val)
}

