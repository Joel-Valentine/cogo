package navigation

import (
	"fmt"
	"regexp"
	"strings"
)

type Validator func(interface{}) error

func ValidateRequired(fieldName string) Validator {
	return func(v interface{}) error {
		if v == nil {
			return NewValidationError(fmt.Sprintf("%s is required", fieldName))
		}

		switch val := v.(type) {
		case string:
			if strings.TrimSpace(val) == "" {
				return NewValidationError(fmt.Sprintf("%s cannot be empty", fieldName))
			}
		case []interface{}:
			if len(val) == 0 {
				return NewValidationError(fmt.Sprintf("%s cannot be empty", fieldName))
			}
		}

		return nil
	}
}

func ValidateLength(fieldName string, min, max int) Validator {
	return func(v interface{}) error {
		str, ok := v.(string)
		if !ok {
			return NewValidationError(fmt.Sprintf("%s must be a string", fieldName))
		}

		length := len(str)

		if min >= 0 && length < min {
			return NewValidationError(fmt.Sprintf("%s must be at least %d characters", fieldName, min))
		}

		if max >= 0 && length > max {
			return NewValidationError(fmt.Sprintf("%s must be at most %d characters", fieldName, max))
		}

		return nil
	}
}

// ValidateDropletName validates a DigitalOcean droplet name.
// Rules:
// - 1-63 characters
// - Letters, numbers, hyphens, underscores, periods
// - Must start with letter or number
func ValidateDropletName(v interface{}) error {
	name, ok := v.(string)
	if !ok {
		return NewValidationError("droplet name must be a string")
	}

	// Length check
	if len(name) == 0 {
		return NewValidationError("droplet name cannot be empty")
	}
	if len(name) > 63 {
		return NewValidationError("droplet name must be 63 characters or less")
	}

	// Character check (alphanumeric, hyphens, underscores, periods)
	validName := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)
	if !validName.MatchString(name) {
		return NewValidationError(
			"droplet name must start with a letter or number and contain only letters, numbers, hyphens, underscores, and periods",
		)
	}

	return nil
}

// ValidateRegex validates that a string matches a regex pattern.
func ValidateRegex(fieldName, pattern, description string) Validator {
	re := regexp.MustCompile(pattern)
	return func(v interface{}) error {
		str, ok := v.(string)
		if !ok {
			return NewValidationError(fmt.Sprintf("%s must be a string", fieldName))
		}

		if !re.MatchString(str) {
			return NewValidationError(fmt.Sprintf("%s must match pattern: %s", fieldName, description))
		}

		return nil
	}
}

// ValidateRange checks that a number is within min/max bounds (inclusive).
func ValidateRange(fieldName string, min, max int) Validator {
	return func(v interface{}) error {
		var num int
		switch val := v.(type) {
		case int:
			num = val
		case int64:
			num = int(val)
		case float64:
			num = int(val)
		default:
			return NewValidationError(fmt.Sprintf("%s must be a number", fieldName))
		}

		if num < min || num > max {
			return NewValidationError(fmt.Sprintf("%s must be between %d and %d", fieldName, min, max))
		}

		return nil
	}
}

// ValidateOneOf checks that a value is one of the allowed values.
func ValidateOneOf(fieldName string, allowed []string) Validator {
	return func(v interface{}) error {
		str, ok := v.(string)
		if !ok {
			return NewValidationError(fmt.Sprintf("%s must be a string", fieldName))
		}

		for _, a := range allowed {
			if str == a {
				return nil
			}
		}

		return NewValidationError(fmt.Sprintf("%s must be one of: %s", fieldName, strings.Join(allowed, ", ")))
	}
}

// Combine combines multiple validators into one.
// All validators must pass for the combined validator to pass.
// Returns the first validation error encountered.
func Combine(validators ...Validator) Validator {
	return func(v interface{}) error {
		for _, validator := range validators {
			if err := validator(v); err != nil {
				return err
			}
		}
		return nil
	}
}

// ValidateResult validates a Result using a Validator.
// Extracts the Value field and passes it to the validator.
func ValidateResult(result Result, validator Validator) error {
	if validator == nil {
		return nil
	}
	return validator(result.Value)
}

// ValidateInput is a helper for common input validation scenarios.
// Returns a Validator that combines required, length, and regex checks.
func ValidateInput(fieldName string, required bool, minLen, maxLen int, pattern, patternDesc string) Validator {
	validators := make([]Validator, 0)

	if required {
		validators = append(validators, ValidateRequired(fieldName))
	}

	if minLen >= 0 || maxLen >= 0 {
		validators = append(validators, ValidateLength(fieldName, minLen, maxLen))
	}

	if pattern != "" {
		validators = append(validators, ValidateRegex(fieldName, pattern, patternDesc))
	}

	return Combine(validators...)
}

