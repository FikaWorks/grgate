package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/fikaworks/grgate/pkg/utils"
)

type statusSetFlagsStruct struct {
	commitSha string
	name      string
	state     string
	status    string
}

var statusSetFlags statusSetFlagsStruct

// statusSetCmd represents the status set command
var statusSetCmd = &cobra.Command{
	Use:   "set [URL OR REPO/OWNER]",
	Short: "Set a status to a given commit",
	Long: `Examples:
  # set the e2e-happy-flow status to completed (github)
  grgate status set my-org/my-repo \
    --commit 36a2dabd4cc732ccab2657392d4a1f8db2f9e19e \
    --name e2e-happy-flow --status completed --state success

  # set the e2e-happy-flow status to success (gitlab)
  grgate status set gitlab.com/my-org/my-repo \
    --commit 36a2dabd4cc732ccab2657392d4a1f8db2f9e19e \
    --name e2e-happy-flow --status success`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		if _, err := utils.ExtractRepository(args[0]); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Info().Msgf("Setting commit %s with status %s = %s and state = %s",
			statusSetFlags.commitSha, statusSetFlags.name, statusSetFlags.status,
			statusSetFlags.state)

		platform, err := newPlatform()
		if err != nil {
			return
		}

		repository, err := utils.ExtractRepository(args[0])
		if err != nil {
			return err
		}

		err = platform.CreateStatus(repository.Owner, repository.Name,
			&platforms.Status{
				Name:      statusSetFlags.name,
				CommitSha: statusSetFlags.commitSha,
				Status:    statusSetFlags.status,
				State:     statusSetFlags.state,
			})
		if err != nil {
			return
		}

		log.Info().Msgf("Status \"%s\" with status \"%s\" created successfully",
			statusSetFlags.name, statusSetFlags.status)

		return
	},
}

func init() {
	statusCmd.AddCommand(statusSetCmd)

	flags := statusSetCmd.Flags()

	flags.StringVar(&statusSetFlags.commitSha, "commit", "", "commit status sha")
	statusGetCmd.MarkFlagRequired("commit")

	flags.StringVar(&statusSetFlags.name, "name", "", "commit status name")
	statusGetCmd.MarkFlagRequired("name")

	flags.StringVar(&statusSetFlags.status, "status", "",
		"for Github status must be one of: queued, in_progress or completed\n"+
			"for Gitlab, status must be one of: pending, running, success, failed\n"+
			"or canceled")
	statusGetCmd.MarkFlagRequired("status")

	flags.StringVar(&statusSetFlags.state, "state", "",
		"(Github only) commit status state is one of action_required, cancelled\n"+
			", failure, neutral, success, skipped, stale, timed_out")
}
