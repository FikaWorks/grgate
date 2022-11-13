package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/fikaworks/grgate/pkg/utils"
)

type statusGetFlagsStruct struct {
	commitSha string
	name      string
}

var statusGetFlags statusGetFlagsStruct

// statusGetCmd represents the status get command
var statusGetCmd = &cobra.Command{
	Use:   "get [URL OR REPO/OWNER]",
	Short: "Get a status attached to a given commit by name",
	Long: `Example:
  # get the e2e-happy-flow status associated to a given commit
  grgate status get my-org/my-repo \
    --commit 36a2dabd4cc732ccab2657392d4a1f8db2f9e19e \
    --name e2e-happy-flow`,
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
		log.Info().Msgf("Retrieving commit status for commit %s in repository %s",
			statusGetFlags.commitSha, args[0])

		platform, err := newPlatform()
		if err != nil {
			return
		}

		repository, err := utils.ExtractRepository(args[0])
		if err != nil {
			return err
		}

		status, err := platform.GetStatus(repository.Owner, repository.Name,
			statusGetFlags.commitSha, statusGetFlags.name)
		if err != nil {
			return
		}

		if status == nil {
			return fmt.Errorf("specified commit status name not found: %s",
				statusGetFlags.name)
		}

		log.Info().Msgf("Found commit status %#v", status)

		return
	},
}

func init() {
	statusCmd.AddCommand(statusGetCmd)

	flags := statusGetCmd.Flags()

	flags.StringVar(&statusGetFlags.commitSha, "commit", "", "commit status sha")
	statusCmd.MarkFlagRequired("commit")

	flags.StringVar(&statusGetFlags.name, "name", "", "commit status name")
	statusCmd.MarkFlagRequired("name")
}
