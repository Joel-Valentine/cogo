// Package credentials provides secure credential management for cogo
// with support for multiple storage backends in priority order.
package credentials

import (
	"context"
	"errors"
)

// Common errors
var (
	ErrTokenNotFound = errors.New("token not found")
	ErrNotSupported  = errors.New("operation not supported by this provider")
)

// Provider defines the interface for credential storage backends
type Provider interface {
	// GetToken retrieves the token from this provider
	GetToken(ctx context.Context) (string, error)

	// SetToken stores the token using this provider
	SetToken(ctx context.Context, token string) error

	// DeleteToken removes the token from this provider
	DeleteToken(ctx context.Context) error

	// Name returns a human-readable name for this provider
	Name() string

	// Available returns true if this provider is available on the system
	Available() bool
}

// Source describes where a token was retrieved from
type Source struct {
	Provider string
	Location string
	Secure   bool
}

// Manager orchestrates multiple credential providers with priority ordering
type Manager struct {
	providers []Provider
}

// NewManager creates a new credential manager with the given providers in priority order
func NewManager(providers ...Provider) *Manager {
	return &Manager{
		providers: providers,
	}
}

// GetToken retrieves a token by trying each provider in order
// Returns the token and the source it came from
func (m *Manager) GetToken(ctx context.Context) (string, *Source, error) {
	for _, provider := range m.providers {
		if !provider.Available() {
			continue
		}

		token, err := provider.GetToken(ctx)
		if err == nil && token != "" {
			return token, &Source{
				Provider: provider.Name(),
				Secure:   isSecureProvider(provider),
				Location: provider.Name(),
			}, nil
		}

		// Continue to next provider on error
		if !errors.Is(err, ErrTokenNotFound) && !errors.Is(err, ErrNotSupported) {
			// Log unexpected errors but continue
			continue
		}
	}

	return "", nil, ErrTokenNotFound
}

// SetToken stores a token using the first available writable provider
func (m *Manager) SetToken(ctx context.Context, token string, providerName string) error {
	for _, provider := range m.providers {
		if !provider.Available() {
			continue
		}

		if providerName != "" && provider.Name() != providerName {
			continue
		}

		if err := provider.SetToken(ctx, token); err != nil {
			if errors.Is(err, ErrNotSupported) {
				continue
			}
			return err
		}
		return nil
	}

	return errors.New("no writable provider available")
}

// DeleteToken removes a token from all providers
func (m *Manager) DeleteToken(ctx context.Context) error {
	var lastErr error
	deleted := false

	for _, provider := range m.providers {
		if !provider.Available() {
			continue
		}

		if err := provider.DeleteToken(ctx); err != nil {
			if !errors.Is(err, ErrNotSupported) && !errors.Is(err, ErrTokenNotFound) {
				lastErr = err
			}
			continue
		}
		deleted = true
	}

	if !deleted && lastErr != nil {
		return lastErr
	}

	return nil
}

// isSecureProvider returns true if the provider is considered secure
func isSecureProvider(p Provider) bool {
	switch p.Name() {
	case "keychain", "environment":
		return true
	case "file", "prompt":
		return false
	default:
		return false
	}
}

// MaskToken returns a masked version of the token showing only first and last 4 characters
func MaskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
