package cmd

import (
	"fmt"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/platforms"
)

func newPlatform() (platform platforms.Platform, err error) {
	switch *config.Main.Platform {
	case config.GitlabPlatform:
		platform, err = platforms.NewGitlab(&platforms.GitlabConfig{
			Token: config.Main.Gitlab.Token,
		})
	case config.GithubPlatform:
		platform, err = platforms.NewGithub(&platforms.GithubConfig{
			AppID:          config.Main.Github.AppID,
			InstallationID: config.Main.Github.InstallationID,
			PrivateKeyPath: config.Main.Github.PrivateKeyPath,
		})
	default:
		err = fmt.Errorf("platform %s is not recognized", *config.Main.Platform)
	}
	return
}
