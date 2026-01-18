package credentials

import (
	"context"
	"fmt"
	
	"github.com/manifoldco/promptui"
)

// PromptProvider retrieves tokens interactively from user input
type PromptProvider struct {
	prompted bool
	token    string
}

// NewPromptProvider creates a new interactive prompt credential provider
func NewPromptProvider() *PromptProvider {
	return &PromptProvider{}
}

// GetToken prompts the user to enter their token
func (p *PromptProvider) GetToken(ctx context.Context) (string, error) {
	if p.prompted {
		return p.token, nil
	}
	
	prompt := promptui.Prompt{
		Label: "Enter your DigitalOcean API Token",
		Mask:  '*',
		Validate: func(input string) error {
			if len(input) == 0 {
				return fmt.Errorf("token cannot be empty")
			}
			return nil
		},
	}
	
	token, err := prompt.Run()
	if err != nil {
		return "", err
	}
	
	p.token = token
	p.prompted = true
	
	return token, nil
}

// SetToken is not supported for prompt provider
func (p *PromptProvider) SetToken(ctx context.Context, token string) error {
	return ErrNotSupported
}

// DeleteToken is not supported for prompt provider
func (p *PromptProvider) DeleteToken(ctx context.Context) error {
	return ErrNotSupported
}

// Name returns the provider name
func (p *PromptProvider) Name() string {
	return "prompt"
}

// Available returns true (prompt is always available)
func (p *PromptProvider) Available() bool {
	return true
}

