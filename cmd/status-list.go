package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/fikaworks/grgate/pkg/utils"
)

type statusListFlagsStruct struct {
	commitSha string
}

var statusListFlags statusListFlagsStruct

// commitStatusListCmd represents the status list command
var statusListCmd = &cobra.Command{
	Use:   "list [OWNER/REPOSITORY]",
	Short: "List statuses attached to a given commit",
	Long: `Example:
  # list statuses associated to a given commit
  grgate status list my-org/my-repo --commit 36a2dabd4cc732ccab2657392d4a1f8db2f9e19e`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		if utils.IsValidRepositoryName(args[0]) {
			return nil
		}
		return fmt.Errorf("invalid repository name specified: %s", args[0])
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Info().Msgf("Listing statuses for commit %s in repository %s",
			statusListFlags.commitSha, args[0])

		platform, err := newPlatform()
		if err != nil {
			return
		}

		statusList, err := platform.ListStatuses(utils.GetRepositoryOrganization(
			args[0]), utils.GetRepositoryName(args[0]), statusListFlags.commitSha)
		if err != nil {
			return
		}

		for _, status := range statusList {
			log.Info().Msgf("Found status %#v", status)
		}

		return
	},
}

func init() {
	statusCmd.AddCommand(statusListCmd)

	flags := statusListCmd.Flags()

	flags.StringVar(&statusListFlags.commitSha, "commit", "", "commit status sha")
	statusListCmd.MarkFlagRequired("commit")
}
