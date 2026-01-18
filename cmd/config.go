package cmd

import (
	"context"
	"fmt"
	"os"
	
	"github.com/Joel-Valentine/cogo/credentials"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	useKeychain bool
	useFile     bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cogo configuration and credentials",
	Long: `Manage cogo configuration including secure credential storage.
	
Credentials are stored securely in your OS keychain by default (macOS Keychain,
Windows Credential Manager, or Linux Secret Service).

You can also use environment variables or legacy file-based storage.`,
}

// setTokenCmd sets the DigitalOcean API token
var setTokenCmd = &cobra.Command{
	Use:   "set-token [token]",
	Short: "Set your DigitalOcean API token",
	Long: `Store your DigitalOcean API token securely.

By default, tokens are stored in your OS keychain. You can also store
in a configuration file using the --file flag (not recommended).

Example:
  cogo config set-token dop_v1_xxx
  cogo config set-token --file dop_v1_xxx
  cogo config set-token  (will prompt for token)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSetToken,
}

// getTokenCmd displays the current token (masked)
var getTokenCmd = &cobra.Command{
	Use:   "get-token",
	Short: "Display your stored API token (masked)",
	Long: `Display your currently stored DigitalOcean API token with masking.
	
Only the first and last 4 characters are shown for security.`,
	RunE: runGetToken,
}

// deleteTokenCmd removes the stored token
var deleteTokenCmd = &cobra.Command{
	Use:   "delete-token",
	Short: "Delete your stored API token",
	Long: `Remove your DigitalOcean API token from all storage locations.
	
This will delete the token from:
- OS keychain
- Configuration files
- All other storage locations`,
	RunE: runDeleteToken,
}

// statusCmd shows credential configuration status
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show credential configuration status",
	Long: `Display information about your current credential configuration,
including where credentials are stored and if they're accessible.`,
	RunE: runStatus,
}

// migrateCmd migrates tokens from file to keychain
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate credentials from file to keychain",
	Long: `Migrate your DigitalOcean API token from the legacy file storage
to secure OS keychain storage.

This will:
1. Read your token from ~/.cogo
2. Store it securely in your OS keychain
3. Optionally remove the plain-text file`,
	RunE: runMigrate,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setTokenCmd)
	configCmd.AddCommand(getTokenCmd)
	configCmd.AddCommand(deleteTokenCmd)
	configCmd.AddCommand(statusCmd)
	configCmd.AddCommand(migrateCmd)
	
	// Flags
	setTokenCmd.Flags().BoolVar(&useKeychain, "keychain", true, "Store in OS keychain (default)")
	setTokenCmd.Flags().BoolVar(&useFile, "file", false, "Store in config file (not recommended)")
}

func runSetToken(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	var token string
	
	// Get token from args or prompt
	if len(args) > 0 {
		token = args[0]
	} else {
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
		
		var err error
		token, err = prompt.Run()
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
	}
	
	// Determine which provider to use
	var provider credentials.Provider
	var providerName string
	
	if useFile {
		provider = credentials.NewFileProvider()
		providerName = "file"
	} else {
		provider = credentials.NewKeychainProvider()
		providerName = "keychain"
	}
	
	// Store the token
	if err := provider.SetToken(ctx, token); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}
	
	color.Green("✓ Token successfully stored in %s", providerName)
	
	if useFile {
		color.Yellow("\n⚠️  WARNING: Token stored in plain text file")
		color.Yellow("   Consider using keychain storage for better security:")
		color.Yellow("   $ cogo config set-token --keychain\n")
	}
	
	return nil
}

func runGetToken(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	
	manager := createManager("", false)
	token, source, err := manager.GetToken(ctx)
	if err != nil {
		if err == credentials.ErrTokenNotFound {
			color.Red("✗ No token found")
			fmt.Println("\nTo set a token, run:")
			fmt.Println("  $ cogo config set-token")
			return nil
		}
		return fmt.Errorf("failed to retrieve token: %w", err)
	}
	
	fmt.Printf("Token: %s\n", credentials.MaskToken(token))
	fmt.Printf("Source: %s\n", source.Provider)
	
	if !source.Secure {
		color.Yellow("\n⚠️  WARNING: Token is stored insecurely")
		color.Yellow("   Consider migrating to keychain storage:")
		color.Yellow("   $ cogo config migrate\n")
	}
	
	return nil
}

func runDeleteToken(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	
	// Confirm deletion
	prompt := promptui.Prompt{
		Label:     "Are you sure you want to delete your stored token?",
		IsConfirm: true,
	}
	
	if _, err := prompt.Run(); err != nil {
		fmt.Println("Cancelled")
		return nil
	}
	
	manager := createManager("", false)
	if err := manager.DeleteToken(ctx); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	
	color.Green("✓ Token deleted from all storage locations")
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	
	fmt.Println("Credential Configuration Status")
	fmt.Println("================================\n")
	
	// Check each provider
	providers := []credentials.Provider{
		credentials.NewEnvProvider(),
		credentials.NewKeychainProvider(),
		credentials.NewFileProvider(),
	}
	
	for _, provider := range providers {
		status := "✗ Not available"
		var details string
		
		if provider.Available() {
			if token, err := provider.GetToken(ctx); err == nil && token != "" {
				status = "✓ Token found"
				details = fmt.Sprintf("(%s)", credentials.MaskToken(token))
			} else {
				status = "○ Available (no token)"
			}
		}
		
		fmt.Printf("%-15s: %s %s\n", provider.Name(), status, details)
	}
	
	// Show current effective token
	fmt.Println("\nEffective Token")
	fmt.Println("---------------")
	
	manager := createManager("", false)
	token, source, err := manager.GetToken(ctx)
	if err != nil {
		if err == credentials.ErrTokenNotFound {
			color.Yellow("No token configured")
			fmt.Println("\nTo set a token, run:")
			fmt.Println("  $ cogo config set-token")
		} else {
			color.Red("Error: %v", err)
		}
	} else {
		fmt.Printf("Token: %s\n", credentials.MaskToken(token))
		fmt.Printf("Source: %s\n", source.Provider)
		
		if !source.Secure {
			color.Yellow("\n⚠️  WARNING: Using insecure storage")
			color.Yellow("   Run 'cogo config migrate' to upgrade\n")
		} else {
			color.Green("✓ Using secure storage")
		}
	}
	
	// Show environment variable info
	fmt.Println("\nEnvironment Variables")
	fmt.Println("--------------------")
	if os.Getenv("DIGITALOCEAN_TOKEN") != "" {
		color.Green("DIGITALOCEAN_TOKEN: Set")
	} else {
		fmt.Println("DIGITALOCEAN_TOKEN: Not set")
	}
	if os.Getenv("COGO_DIGITALOCEAN_TOKEN") != "" {
		color.Green("COGO_DIGITALOCEAN_TOKEN: Set")
	} else {
		fmt.Println("COGO_DIGITALOCEAN_TOKEN: Not set")
	}
	
	return nil
}

func runMigrate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	
	fileProvider := credentials.NewFileProvider()
	keychainProvider := credentials.NewKeychainProvider()
	
	// Check if keychain is available
	if !keychainProvider.Available() {
		color.Red("✗ OS keychain is not available on this system")
		fmt.Println("\nConsider using environment variables instead:")
		fmt.Println("  export DIGITALOCEAN_TOKEN=your_token_here")
		return nil
	}
	
	// Check if file has token
	if !fileProvider.Available() {
		color.Yellow("No config file found to migrate")
		return nil
	}
	
	token, err := fileProvider.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to read token from file: %w", err)
	}
	
	fmt.Printf("Found token in file: %s\n", credentials.MaskToken(token))
	
	// Store in keychain
	if err := keychainProvider.SetToken(ctx, token); err != nil {
		return fmt.Errorf("failed to store token in keychain: %w", err)
	}
	
	color.Green("✓ Token successfully stored in keychain")
	
	// Ask if they want to delete the file
	prompt := promptui.Prompt{
		Label:     "Delete the plain-text config file?",
		IsConfirm: true,
	}
	
	if _, err := prompt.Run(); err == nil {
		if err := fileProvider.DeleteToken(ctx); err != nil {
			color.Yellow("⚠  Failed to delete config file: %v", err)
		} else {
			color.Green("✓ Plain-text config file deleted")
		}
	} else {
		color.Yellow("⚠  Keeping plain-text config file")
		fmt.Println("  You can manually delete it at: ~/.cogo")
	}
	
	color.Green("\n✓ Migration complete!")
	fmt.Println("Your token is now stored securely in your OS keychain.")
	
	return nil
}

// createManager creates a credential manager with the standard provider chain
// flagToken is an optional token from CLI flag
// includePrompt determines whether to include the interactive prompt provider
func createManager(flagToken string, includePrompt bool) *credentials.Manager {
	providers := []credentials.Provider{
		credentials.NewFlagProvider(flagToken),
		credentials.NewEnvProvider(),
		credentials.NewKeychainProvider(),
		credentials.NewFileProvider(),
	}
	
	if includePrompt {
		providers = append(providers, credentials.NewPromptProvider())
	}
	
	return credentials.NewManager(providers...)
}

