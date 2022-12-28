package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Check if a flag is present
func flagPresent(cmd *cobra.Command, flagName string) bool {
	isPresent, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		log.Fatal(err)
	}
	return isPresent
}

// Cobra root command — this is the entrypoint for the CLI
var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("sh <(curl %s) [flags]", Config.InstallURL),
	Short: Config.HelpDescription,
	Run: func(cmd *cobra.Command, args []string) {
		options := map[string]bool{}
		for k := range Config.Installers {
			options[k] = flagPresent(cmd, k)
		}
		tui(options)
	},
}

// Execute the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Add flags to the root command
func init() {
	// Update default help message
	rootCmd.Flags().BoolP("help", "h", false, "Show this help message")

	// Add flag for temporary install
	rootCmd.Flags().BoolP("tmp", "", false, "Install temporarily to "+Config.TmpDir)

	// Add flags for all installers
	for flag, v := range Config.Installers {
		long_help_message := v.HelpMessage
		if flag == "full" {
			long_help_message = "Full system config"
		}
		rootCmd.Flags().BoolP(flag, "", false, long_help_message)
	}
}
