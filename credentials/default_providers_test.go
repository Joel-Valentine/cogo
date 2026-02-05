package credentials

import (
	"testing"
)

func TestNewManager_DefaultProviders(t *testing.T) {
	// Create manager with no arguments
	manager := NewManager()

	// Should have default providers (5 of them)
	if len(manager.providers) != 5 {
		t.Errorf("Expected 5 default providers, got %d", len(manager.providers))
	}

	// Verify the providers are in the correct order
	expectedTypes := []string{"flag", "env", "keychain", "file", "prompt"}
	for i, provider := range manager.providers {
		var providerType string
		switch provider.(type) {
		case *FlagProvider:
			providerType = "flag"
		case *EnvProvider:
			providerType = "env"
		case *KeychainProvider:
			providerType = "keychain"
		case *FileProvider:
			providerType = "file"
		case *PromptProvider:
			providerType = "prompt"
		}

		if providerType != expectedTypes[i] {
			t.Errorf("Provider %d: expected %s, got %s", i, expectedTypes[i], providerType)
		}
	}
}

func TestNewManager_CustomProviders(t *testing.T) {
	// Create manager with custom providers
	custom := &mockProvider{
		name:      "custom",
		token:     "custom-token",
		available: true,
	}

	manager := NewManager(custom)

	// Should have only the custom provider
	if len(manager.providers) != 1 {
		t.Errorf("Expected 1 custom provider, got %d", len(manager.providers))
	}
}

