package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/fikaworks/ggate/pkg/config"
	"github.com/fikaworks/ggate/pkg/platforms"
	"github.com/fikaworks/ggate/pkg/utils"
	"github.com/fikaworks/ggate/pkg/workers"
)

type runCmdFlagsStruct struct {
	dryRun    bool
	tagRegexp string
	statuses  []string
}

var runCmdFlags runCmdFlagsStruct

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [OWNER/REPOSITORY]",
	Short: "Run GGate against a repository",
	Long: `The run command list all the draft/unpublished releases from a given
repository that match the provided tag. From this list, if all the status check
are completed and successful and match the list of provided status, then the
release is published.`,
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
		owner := utils.GetRepositoryOrganization(args[0])
		repository := utils.GetRepositoryName(args[0])

		if runCmdFlags.dryRun {
			log.Info().Msg("Executing command with dry-run mode enabled")
		}

		platform, err := platforms.NewGithub(&platforms.GithubConfig{
			AppID:          config.Main.Github.AppID,
			InstallationID: config.Main.Github.InstallationID,
			PrivateKeyPath: config.Main.Github.PrivateKeyPath,
		})
		if err != nil {
			return err
		}

		job, err := workers.NewJob(platform, owner, repository)
		if err != nil {
			return err
		}

		// override status from command line if defined
		if len(runCmdFlags.statuses) > 0 {
			job.Config.Statuses = runCmdFlags.statuses
		}
		job.Config.TagRegexp = runCmdFlags.tagRegexp
		job.Config.Enabled = !runCmdFlags.dryRun

		// process the repository
		return job.Process()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	flags := runCmd.Flags()

	flags.BoolVar(&runCmdFlags.dryRun, "dry-run", false, "dry run")
	flags.StringVar(&runCmdFlags.tagRegexp, "tag-regexp", ".*", "tag regexp")
	flags.StringArrayVarP(&runCmdFlags.statuses, "status", "s", []string{},
		"List of status to succeed")
}
