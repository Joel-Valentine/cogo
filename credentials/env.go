package credentials

import (
	"context"
	"os"
)

// EnvProvider retrieves tokens from environment variables
type EnvProvider struct {
	envVars []string
}

// NewEnvProvider creates a new environment variable credential provider
// It checks the provided environment variable names in order
func NewEnvProvider(envVars ...string) *EnvProvider {
	if len(envVars) == 0 {
		// Default environment variables
		envVars = []string{"DIGITALOCEAN_TOKEN", "COGO_DIGITALOCEAN_TOKEN"}
	}
	return &EnvProvider{
		envVars: envVars,
	}
}

// GetToken retrieves the token from environment variables
func (p *EnvProvider) GetToken(ctx context.Context) (string, error) {
	for _, envVar := range p.envVars {
		if token := os.Getenv(envVar); token != "" {
			return token, nil
		}
	}
	return "", ErrTokenNotFound
}

// SetToken is not supported for environment provider (read-only)
func (p *EnvProvider) SetToken(ctx context.Context, token string) error {
	return ErrNotSupported
}

// DeleteToken is not supported for environment provider (read-only)
func (p *EnvProvider) DeleteToken(ctx context.Context) error {
	return ErrNotSupported
}

// Name returns the provider name
func (p *EnvProvider) Name() string {
	return "environment"
}

// Available returns true if any of the environment variables are set
func (p *EnvProvider) Available() bool {
	for _, envVar := range p.envVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}
	return false
}
