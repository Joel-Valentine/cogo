package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Joel-Valentine/cogo/credentials"
	do "github.com/Joel-Valentine/cogo/digitalocean"
	"github.com/Joel-Valentine/cogo/utils"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Cogo create, list, destroy wizard",
	Short: "For interacting with multiple cloud providers",
	Long:  `Cogo is a CLI tool used to interact easily as a wizard with multiple cloud providers`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(create)
	rootCmd.AddCommand(list)
	rootCmd.AddCommand(destroy)
	cobra.OnInitialize()
}

var create = &cobra.Command{
	Use:   "create",
	Short: "Creates a server in selected provider",
	Long:  `Will walk you through a wizard to create a server in a selected provider`,
	Run: func(cmd *cobra.Command, args []string) {

		selectedProvider, err := utils.AskForProvider()

		if err != nil {
			color.Yellow("Something went wrong asking for selected provider\n")
			return
		}

		if selectedProvider == "DO" {
			// Get DigitalOcean API token
			ctx := context.Background()
			credManager := credentials.NewManager() // Uses default providers
			token, _, err := credManager.GetToken(ctx)
			if err != nil {
				color.Red("✗ Error: Unable to get DigitalOcean API token\n")
				fmt.Println()
				fmt.Println("Run 'cogo config set-token' to configure your token.")
				return
			}

			// Create client
			client := godo.NewFromToken(token)

			// Execute create flow with back navigation support
			createdDroplet, createDropletError := do.ExecuteCreateFlow(client)

			if createDropletError != nil {
				color.Red("✗ Error: %v\n", createDropletError)
				return
			}

			if createdDroplet == nil {
				// Canceled or empty state (message already shown)
				return
			}

			color.Green("✓ Droplet [%s] was created!", createdDroplet.Name)
			color.Cyan("List your droplets in a couple of minutes to see the IP\n")
		}
	},
}

var list = &cobra.Command{
	Use:   "list",
	Short: "Lists servers created in selected provider",
	Long:  `Will show a list of servers that you currently have in a selected provider`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedProvider, err := utils.AskForProvider()

		if err != nil {
			color.Yellow("Something went wrong asking for selected provider\n")
			return
		}

		if selectedProvider == "DO" {
			do.DisplayDropletList()
		}
	},
}

var destroy = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys servers created in selected provider",
	Long: `Will show a list of servers that you currently have in a selected provider, 
	with the ability to select one and delete/destroy it.
	
	Be very careful here. There will be multiple warnings to make sure that you don't accidentally delete
	a crucial droplet`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedProvider, err := utils.AskForProvider()

		if err != nil {
			color.Yellow("Something went wrong asking for selected provider\n")
			return
		}

		if selectedProvider == "DO" {
			// Get DigitalOcean API token
			ctx := context.Background()
			credManager := credentials.NewManager() // Uses default providers
			token, _, err := credManager.GetToken(ctx)
			if err != nil {
				color.Red("✗ Error: Unable to get DigitalOcean API token\n")
				fmt.Println()
				fmt.Println("Run 'cogo config set-token' to configure your token.")
				return
			}

			// Create client
			client := godo.NewFromToken(token)

			// Execute destroy flow with back navigation support
			destroyedDroplet, err := do.ExecuteDestroyFlow(client)

			if err != nil {
				color.Red("✗ Error: %v\n", err)
				return
			}

			if destroyedDroplet == nil {
				// Canceled or empty state (message already shown)
				return
			}

			color.Green("✓ Droplet [%s] has been destroyed\n", destroyedDroplet.Name)
		}
	},
}
