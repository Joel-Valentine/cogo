package credentials

import (
	"context"
)

// FlagProvider retrieves tokens from command-line flags
type FlagProvider struct {
	token string
}

// NewFlagProvider creates a new flag-based credential provider
func NewFlagProvider(token string) *FlagProvider {
	return &FlagProvider{
		token: token,
	}
}

// GetToken returns the token from the flag
func (p *FlagProvider) GetToken(ctx context.Context) (string, error) {
	if p.token == "" {
		return "", ErrTokenNotFound
	}
	return p.token, nil
}

// SetToken is not supported for flag provider
func (p *FlagProvider) SetToken(ctx context.Context, token string) error {
	return ErrNotSupported
}

// DeleteToken is not supported for flag provider
func (p *FlagProvider) DeleteToken(ctx context.Context) error {
	return ErrNotSupported
}

// Name returns the provider name
func (p *FlagProvider) Name() string {
	return "flag"
}

// Available returns true if a token was provided via flag
func (p *FlagProvider) Available() bool {
	return p.token != ""
}
