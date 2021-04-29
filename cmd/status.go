package cmd

import (
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [COMMANDS]",
	Short: "Interact with commit status",
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

