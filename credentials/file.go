package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// FileProvider retrieves tokens from legacy config files (deprecated but supported)
type FileProvider struct {
	configPath string
	warned     bool
}

// NewFileProvider creates a new file-based credential provider
func NewFileProvider() *FileProvider {
	return &FileProvider{}
}

// GetToken retrieves the token from the config file
func (p *FileProvider) GetToken(ctx context.Context) (string, error) {
	v := viper.New()
	v.SetConfigName(".cogo")
	v.SetConfigType("json")
	v.AddConfigPath("$HOME")
	v.AddConfigPath("$HOME/.config/")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return "", ErrTokenNotFound
	}

	p.configPath = v.ConfigFileUsed()

	token := v.GetString("digitaloceantoken")
	if token == "" {
		// Try alternate key
		token = v.GetString("digitalOceanToken")
	}

	if token == "" {
		return "", ErrTokenNotFound
	}

	// Show warning about insecure storage (only once)
	if !p.warned {
		fmt.Fprintf(os.Stderr, "\n⚠️  WARNING: Token stored in plain text file: %s\n", p.configPath)
		fmt.Fprintf(os.Stderr, "   Consider migrating to secure keychain storage:\n")
		fmt.Fprintf(os.Stderr, "   $ cogo config migrate\n\n")
		p.warned = true
	}

	return token, nil
}

// SetToken stores the token in the config file (deprecated)
func (p *FileProvider) SetToken(ctx context.Context, token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".cogo")

	// Read existing config if it exists
	config := make(map[string]interface{})
	if data, readErr := os.ReadFile(configPath); readErr == nil {
		if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
			// Ignore unmarshal errors for existing files, we'll overwrite
		}
	}

	// Update token
	config["digitaloceantoken"] = token

	// Write back
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write with restrictive permissions
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "\n⚠️  WARNING: Token stored in plain text file: %s\n", configPath)
	fmt.Fprintf(os.Stderr, "   Consider using keychain storage instead:\n")
	fmt.Fprintf(os.Stderr, "   $ cogo config set-token --keychain\n\n")

	return nil
}

// DeleteToken removes the token from the config file
func (p *FileProvider) DeleteToken(ctx context.Context) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".cogo")

	// Check if file exists
	if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
		return ErrTokenNotFound
	}

	// Read existing config
	config := make(map[string]interface{})
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if unmarshalErr := json.Unmarshal(data, &config); unmarshalErr != nil {
		return unmarshalErr
	}

	// Remove token keys
	delete(config, "digitaloceantoken")
	delete(config, "digitalOceanToken")

	// If config is now empty, delete the file
	if len(config) == 0 {
		return os.Remove(configPath)
	}

	// Otherwise, write back without token
	data, err = json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// Name returns the provider name
func (p *FileProvider) Name() string {
	return "file"
}

// Available returns true if a config file exists
func (p *FileProvider) Available() bool {
	v := viper.New()
	v.SetConfigName(".cogo")
	v.SetConfigType("json")
	v.AddConfigPath("$HOME")
	v.AddConfigPath("$HOME/.config/")
	v.AddConfigPath(".")

	return v.ReadInConfig() == nil
}
