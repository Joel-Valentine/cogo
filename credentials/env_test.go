package credentials

import (
	"context"
	"os"
	"testing"
)

func TestEnvProvider_GetToken(t *testing.T) {
	ctx := context.Background()
	
	// Save original env vars
	originalDO := os.Getenv("DIGITALOCEAN_TOKEN")
	originalCogo := os.Getenv("COGO_DIGITALOCEAN_TOKEN")
	defer func() {
		os.Setenv("DIGITALOCEAN_TOKEN", originalDO)
		os.Setenv("COGO_DIGITALOCEAN_TOKEN", originalCogo)
	}()
	
	tests := []struct {
		name          string
		doToken       string
		cogoToken     string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "DIGITALOCEAN_TOKEN set",
			doToken:       "do-token",
			cogoToken:     "",
			expectedToken: "do-token",
			expectError:   false,
		},
		{
			name:          "COGO_DIGITALOCEAN_TOKEN set",
			doToken:       "",
			cogoToken:     "cogo-token",
			expectedToken: "cogo-token",
			expectError:   false,
		},
		{
			name:          "both set, DIGITALOCEAN_TOKEN takes priority",
			doToken:       "do-token",
			cogoToken:     "cogo-token",
			expectedToken: "do-token",
			expectError:   false,
		},
		{
			name:        "neither set",
			doToken:     "",
			cogoToken:   "",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DIGITALOCEAN_TOKEN", tt.doToken)
			os.Setenv("COGO_DIGITALOCEAN_TOKEN", tt.cogoToken)
			
			provider := NewEnvProvider()
			token, err := provider.GetToken(ctx)
			
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("expected token %q, got %q", tt.expectedToken, token)
				}
			}
		})
	}
}

func TestEnvProvider_SetToken_NotSupported(t *testing.T) {
	ctx := context.Background()
	provider := NewEnvProvider()
	
	err := provider.SetToken(ctx, "test-token")
	if err != ErrNotSupported {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}

func TestEnvProvider_Available(t *testing.T) {
	// Save original env vars
	originalDO := os.Getenv("DIGITALOCEAN_TOKEN")
	defer os.Setenv("DIGITALOCEAN_TOKEN", originalDO)
	
	tests := []struct {
		name      string
		token     string
		available bool
	}{
		{
			name:      "token set",
			token:     "test-token",
			available: true,
		},
		{
			name:      "token not set",
			token:     "",
			available: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DIGITALOCEAN_TOKEN", tt.token)
			
			provider := NewEnvProvider()
			available := provider.Available()
			
			if available != tt.available {
				t.Errorf("expected Available() = %v, got %v", tt.available, available)
			}
		})
	}
}

