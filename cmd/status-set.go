package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/rs/zerolog/log"

	"github.com/fikaworks/ggate/pkg/config"
	"github.com/fikaworks/ggate/pkg/platforms"
	"github.com/fikaworks/ggate/pkg/utils"
)

type statusSetFlagsStruct struct {
  commitSha string
  name string
  state string
  status string
}

var statusSetFlags statusSetFlagsStruct

// statusSetCmd represents the status set command
var statusSetCmd = &cobra.Command{
	Use:   "set [OWNER/REPOSITORY]",
	Short: "Set a status to a given commit",
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
		log.Info().Msgf("Setting commit %s with status %s = %s and state = %s",
      statusSetFlags.commitSha, statusSetFlags.name, statusSetFlags.status,
      statusSetFlags.state)

    platform, err := platforms.NewGithub(&platforms.GithubConfig{
      AppID: config.Main.Github.AppID,
      InstallationID: config.Main.Github.InstallationID,
      PrivateKeyPath: config.Main.Github.PrivateKeyPath,
    })
    if err != nil {
      return
    }

		err = platform.CreateStatus(utils.GetRepositoryOrganization(args[0]),
      utils.GetRepositoryName(args[0]), &platforms.Status{
        Name: statusSetFlags.name,
        CommitSha: statusSetFlags.commitSha,
        Status: statusSetFlags.status,
        State: statusSetFlags.state,
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

	statusSetCmd.Flags().StringVar(&statusSetFlags.commitSha, "commit", "",
    "commit status sha")
  statusGetCmd.MarkFlagRequired("commit")

	statusSetCmd.Flags().StringVar(&statusSetFlags.name, "name", "",
    "commit status name")
  statusGetCmd.MarkFlagRequired("name")

	statusSetCmd.Flags().StringVar(&statusSetFlags.status, "status", "",
    "status, one of \"queued\", \"in_progress\", \"completed\"")
  statusGetCmd.MarkFlagRequired("status")

	statusSetCmd.Flags().StringVar(&statusSetFlags.state, "state", "",
    "commit status state is one of \"success\", \"in_progress\"")
}
