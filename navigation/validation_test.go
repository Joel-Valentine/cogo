package navigation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequired(t *testing.T) {
	validator := ValidateRequired("name")

	// Valid values
	assert.NoError(t, validator("test"))
	assert.NoError(t, validator("  test  ")) // Trimmed

	// Invalid values
	assert.Error(t, validator(""))
	assert.Error(t, validator("   "))
	assert.Error(t, validator(nil))
}

func TestValidateLength(t *testing.T) {
	// Min and max
	validator := ValidateLength("name", 3, 10)

	assert.NoError(t, validator("test"))
	assert.NoError(t, validator("123"))
	assert.NoError(t, validator("1234567890"))

	assert.Error(t, validator("ab")) // Too short
	assert.Error(t, validator("12345678901")) // Too long

	// Only min
	validator = ValidateLength("name", 5, -1)
	assert.NoError(t, validator("12345"))
	assert.NoError(t, validator("123456789012345"))
	assert.Error(t, validator("1234"))

	// Only max
	validator = ValidateLength("name", -1, 5)
	assert.NoError(t, validator(""))
	assert.NoError(t, validator("12345"))
	assert.Error(t, validator("123456"))
}

func TestValidateDropletName(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"valid lowercase", "my-droplet", false},
		{"valid mixed case", "My-Droplet", false},
		{"valid with numbers", "droplet123", false},
		{"valid with underscore", "my_droplet", false},
		{"valid with period", "my.droplet", false},
		{"valid complex", "my-droplet_1.test", false},
		{"empty", "", true},
		{"too long", "this-is-a-very-long-droplet-name-that-exceeds-the-maximum-length-allowed", true},
		{"starts with hyphen", "-droplet", true},
		{"starts with period", ".droplet", true},
		{"invalid chars", "my droplet", true},
		{"invalid chars special", "my@droplet", true},
		{"not a string", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDropletName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, ErrValidationFailed))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRegex(t *testing.T) {
	validator := ValidateRegex("email", `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "valid email")

	assert.NoError(t, validator("test@example.com"))
	assert.NoError(t, validator("user.name+tag@example.co.uk"))

	assert.Error(t, validator("invalid"))
	assert.Error(t, validator("@example.com"))
	assert.Error(t, validator("test@"))
}

func TestValidateRange(t *testing.T) {
	validator := ValidateRange("count", 1, 10)

	// Valid
	assert.NoError(t, validator(1))
	assert.NoError(t, validator(5))
	assert.NoError(t, validator(10))
	assert.NoError(t, validator(int64(5)))
	assert.NoError(t, validator(float64(5)))

	// Invalid
	assert.Error(t, validator(0))
	assert.Error(t, validator(11))
	assert.Error(t, validator(-1))
	assert.Error(t, validator("not a number"))
}

func TestValidateOneOf(t *testing.T) {
	validator := ValidateOneOf("region", []string{"nyc1", "nyc3", "sfo3"})

	assert.NoError(t, validator("nyc1"))
	assert.NoError(t, validator("sfo3"))

	assert.Error(t, validator("invalid"))
	assert.Error(t, validator(""))
	assert.Error(t, validator("NYC1")) // Case sensitive
}

func TestCombine(t *testing.T) {
	validator := Combine(
		ValidateRequired("name"),
		ValidateLength("name", 3, 10),
	)

	// Valid
	assert.NoError(t, validator("test"))

	// Fails first validator (required)
	err := validator("")
	assert.Error(t, err)

	// Fails second validator (length)
	err = validator("ab")
	assert.Error(t, err)
}

func TestValidateResult(t *testing.T) {
	result := NewResult("test")
	validator := ValidateLength("value", 3, 10)

	// Valid
	assert.NoError(t, ValidateResult(result, validator))

	// Invalid
	result = NewResult("ab")
	assert.Error(t, ValidateResult(result, validator))

	// Nil validator
	assert.NoError(t, ValidateResult(result, nil))
}

func TestValidateInput(t *testing.T) {
	// Required, length, pattern
	validator := ValidateInput("username", true, 3, 20, `^[a-zA-Z][a-zA-Z0-9_]*$`, "alphanumeric starting with letter")

	// Valid
	assert.NoError(t, validator("user123"))
	assert.NoError(t, validator("John_Doe"))

	// Too short
	assert.Error(t, validator("ab"))

	// Invalid pattern
	assert.Error(t, validator("123user"))
	assert.Error(t, validator("user@test"))

	// Empty (required)
	assert.Error(t, validator(""))
}

