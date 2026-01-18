package credentials

import (
	"context"
	"testing"
)

func TestMaskToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "standard token",
			token:    "fake_token_1234567890abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnop",
			expected: "fake...mnop",
		},
		{
			name:     "short token",
			token:    "short",
			expected: "***",
		},
		{
			name:     "exactly 8 chars",
			token:    "12345678",
			expected: "***",
		},
		{
			name:     "9 chars",
			token:    "123456789",
			expected: "1234...6789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskToken(tt.token)
			if result != tt.expected {
				t.Errorf("MaskToken(%q) = %q, want %q", tt.token, result, tt.expected)
			}
		})
	}
}

func TestManager_GetToken_Priority(t *testing.T) {
	ctx := context.Background()

	// Create mock providers
	highPriority := &mockProvider{
		name:  "high",
		token: "high-token",
	}
	lowPriority := &mockProvider{
		name:  "low",
		token: "low-token",
	}

	manager := NewManager(highPriority, lowPriority)

	token, source, err := manager.GetToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "high-token" {
		t.Errorf("expected token from high priority provider, got %q", token)
	}

	if source.Provider != "high" {
		t.Errorf("expected source 'high', got %q", source.Provider)
	}
}

func TestManager_GetToken_Fallback(t *testing.T) {
	ctx := context.Background()

	// Create providers where first one fails
	failingProvider := &mockProvider{
		name:      "failing",
		available: false,
	}
	workingProvider := &mockProvider{
		name:  "working",
		token: "working-token",
	}

	manager := NewManager(failingProvider, workingProvider)

	token, source, err := manager.GetToken(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "working-token" {
		t.Errorf("expected token from working provider, got %q", token)
	}

	if source.Provider != "working" {
		t.Errorf("expected source 'working', got %q", source.Provider)
	}
}

func TestManager_GetToken_NotFound(t *testing.T) {
	ctx := context.Background()

	emptyProvider := &mockProvider{
		name:  "empty",
		token: "",
	}

	manager := NewManager(emptyProvider)

	_, _, err := manager.GetToken(ctx)
	if err != ErrTokenNotFound {
		t.Errorf("expected ErrTokenNotFound, got %v", err)
	}
}

// mockProvider is a test implementation of Provider
type mockProvider struct {
	name     string
	token    string
	setError error
	available bool
}

func (m *mockProvider) GetToken(ctx context.Context) (string, error) {
	if m.token == "" {
		return "", ErrTokenNotFound
	}
	return m.token, nil
}

func (m *mockProvider) SetToken(ctx context.Context, token string) error {
	if m.setError != nil {
		return m.setError
	}
	m.token = token
	return nil
}

func (m *mockProvider) DeleteToken(ctx context.Context) error {
	m.token = ""
	return nil
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Available() bool {
	if m.available {
		return true
	}
	// Default to available if token is set
	return m.token != ""
}
