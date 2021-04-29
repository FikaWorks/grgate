package cmd

import (
	"github.com/spf13/cobra"
	"github.com/rs/zerolog/log"
)

// versionCmd represents the status command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GGate",
  Run: func(cmd *cobra.Command, args []string) {
    log.Info().Msgf("Version %s", Version)
  },
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

