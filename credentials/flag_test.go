package credentials

import (
	"context"
	"testing"
)

func TestFlagProvider_GetToken(t *testing.T) {
	ctx := context.Background()
	
	tests := []struct {
		name        string
		flagToken   string
		expectError bool
	}{
		{
			name:        "token provided",
			flagToken:   "flag-token",
			expectError: false,
		},
		{
			name:        "no token",
			flagToken:   "",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewFlagProvider(tt.flagToken)
			token, err := provider.GetToken(ctx)
			
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token != tt.flagToken {
					t.Errorf("expected token %q, got %q", tt.flagToken, token)
				}
			}
		})
	}
}

func TestFlagProvider_Available(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		available bool
	}{
		{
			name:      "token provided",
			token:     "test-token",
			available: true,
		},
		{
			name:      "no token",
			token:     "",
			available: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewFlagProvider(tt.token)
			available := provider.Available()
			
			if available != tt.available {
				t.Errorf("expected Available() = %v, got %v", tt.available, available)
			}
		})
	}
}

func TestFlagProvider_SetToken_NotSupported(t *testing.T) {
	ctx := context.Background()
	provider := NewFlagProvider("")
	
	err := provider.SetToken(ctx, "test-token")
	if err != ErrNotSupported {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

