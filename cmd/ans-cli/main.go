package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/kyma-project/ans-manager/internal/cli/command"
)

var (
	debug bool
)

var rootCmd = &cobra.Command{
	Use:   "ans-cli",
	Short: "ANS Manager CLI - Send notifications to the ANS service",
	Long: `ANS Manager CLI is a command line tool for sending events and notifications 
to the SAP BTP Alert Notification Service (ANS). 

It supports OAuth authentication and multiple regions.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode to show request payloads")

	rootCmd.AddCommand(command.NewNotifyCommand())
}

func GetDebugFlag() bool {
	return debug
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
