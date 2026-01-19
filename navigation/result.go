package navigation

type Result struct {
	Value    interface{}
	Metadata map[string]interface{}
}

func NewResult(value interface{}) Result {
	return Result{
		Value:    value,
		Metadata: make(map[string]interface{}),
	}
}

func NewResultWithMetadata(value interface{}, metadata map[string]interface{}) Result {
	return Result{
		Value:    value,
		Metadata: metadata,
	}
}

func (r Result) WithMetadata(key string, value interface{}) Result {
	if r.Metadata == nil {
		r.Metadata = make(map[string]interface{})
	}
	r.Metadata[key] = value
	return r
}

func (r Result) GetMetadata(key string) (interface{}, bool) {
	if r.Metadata == nil {
		return nil, false
	}
	val, ok := r.Metadata[key]
	return val, ok
}

type StepResult struct {
	StepName string
	Result   Result
}

