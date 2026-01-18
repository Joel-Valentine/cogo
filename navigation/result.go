package navigation

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

// NewResult creates a Result with a value and empty metadata.
func NewResult(value interface{}) Result {
	return Result{
		Value:    value,
		Metadata: make(map[string]interface{}),
	}
}

// NewResultWithMetadata creates a Result with a value and metadata.
func NewResultWithMetadata(value interface{}, metadata map[string]interface{}) Result {
	return Result{
		Value:    value,
		Metadata: metadata,
	}
}

// WithMetadata adds a metadata key-value pair and returns the Result (builder pattern).
func (r Result) WithMetadata(key string, value interface{}) Result {
	if r.Metadata == nil {
		r.Metadata = make(map[string]interface{})
	}
	r.Metadata[key] = value
	return r
}

// GetMetadata retrieves a metadata value by key.
// Returns (value, found bool).
func (r Result) GetMetadata(key string) (interface{}, bool) {
	if r.Metadata == nil {
		return nil, false
	}
	val, ok := r.Metadata[key]
	return val, ok
}

// StepResult pairs a Step with its Result for history tracking.
type StepResult struct {
	StepName string
	Result   Result
}

