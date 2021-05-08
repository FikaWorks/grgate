package cmd

import (
	"github.com/fikaworks/ggate/pkg/config"
	"github.com/fikaworks/ggate/pkg/platforms"
)

func newPlatform() (platform platforms.Platform, err error) {
	if config.Main.Gitlab.Token != "" {
		platform, err = platforms.NewGitlab(&platforms.GitlabConfig{
			Token: config.Main.Gitlab.Token,
		})
	} else {
		platform, err = platforms.NewGithub(&platforms.GithubConfig{
			AppID:          config.Main.Github.AppID,
			InstallationID: config.Main.Github.InstallationID,
			PrivateKeyPath: config.Main.Github.PrivateKeyPath,
		})
	}
	return
}
