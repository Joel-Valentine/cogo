package cmd

import (
	"fmt"
	"os"

	do "github.com/Midnight-Conqueror/cogo/digitalocean"
	"github.com/Midnight-Conqueror/cogo/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Cogo create, list, destroy wizard",
	Short: "For interacting with multiple cloud providers",
	Long:  `Cogo is a CLI tool used to intreact easily as a wizard with multiple cloud providers`,
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
	cobra.OnInitialize()
}

var create = &cobra.Command{
	Use:   "create",
	Short: "Creates a server in selected provider",
	Long:  `Will walk you through a wizard to create a server in a selected provider`,
	Run: func(cmd *cobra.Command, args []string) {

		selectedProvider, err := utils.AskForProvider()

		if err != nil {
			fmt.Print("Something went wrong asking for selected provider")
			return
		}

		if selectedProvider == "DO" {
			createdDroplet, createDropletError := do.CreateDroplet()

			if createDropletError != nil {
				fmt.Printf("Creation of droplet failed: %s\n\n", createDropletError)
			}

			if createdDroplet == nil {
				color.Cyan("Droplet was not created")
				return
			}

			color.Green("Droplet %s was created!", createdDroplet.Name)
			color.Cyan("List your droplets in a couple of minutes to see the IP")
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
			fmt.Print("Something went wrong asking for selected provider")
			return
		}

		if selectedProvider == "DO" {
			do.DisplayDropletList()
		}
	},
}
