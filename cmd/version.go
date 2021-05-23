package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/fikaworks/grgate/pkg/config"
)

// versionCmd represents the status command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GRGate",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("Version %s", config.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
