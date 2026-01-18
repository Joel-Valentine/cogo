package credentials

import (
	"context"
	"errors"
	
	"github.com/zalando/go-keyring"
)

const (
	keychainService = "cogo"
	keychainAccount = "digitalocean-token"
)

// KeychainProvider retrieves tokens from the OS keychain
type KeychainProvider struct{}

// NewKeychainProvider creates a new keychain-based credential provider
func NewKeychainProvider() *KeychainProvider {
	return &KeychainProvider{}
}

// GetToken retrieves the token from the OS keychain
func (p *KeychainProvider) GetToken(ctx context.Context) (string, error) {
	token, err := keyring.Get(keychainService, keychainAccount)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrTokenNotFound
		}
		return "", err
	}
	return token, nil
}

// SetToken stores the token in the OS keychain
func (p *KeychainProvider) SetToken(ctx context.Context, token string) error {
	return keyring.Set(keychainService, keychainAccount, token)
}

// DeleteToken removes the token from the OS keychain
func (p *KeychainProvider) DeleteToken(ctx context.Context) error {
	err := keyring.Delete(keychainService, keychainAccount)
	if err != nil && errors.Is(err, keyring.ErrNotFound) {
		return ErrTokenNotFound
	}
	return err
}

// Name returns the provider name
func (p *KeychainProvider) Name() string {
	return "keychain"
}

// Available returns true if the keychain is available on this system
func (p *KeychainProvider) Available() bool {
	// Try to access keychain - if it fails, it's not available
	// This handles headless Linux systems gracefully
	_, err := keyring.Get(keychainService, "_cogo_availability_test")
	return err == nil || errors.Is(err, keyring.ErrNotFound)
}

