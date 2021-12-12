//go:build integration || integrationgithub

package tests

import (
	"os"
	"strconv"
	"testing"

	"github.com/fikaworks/grgate/pkg/platforms"
)

func TestGithubReleases(t *testing.T) {
	githubAppID, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	githubInstallationID, _ := strconv.ParseInt(os.Getenv("GITHUB_INSTALLATION_ID"), 10, 64)
	githubPrivateKeyPath := os.Getenv("GITHUB_PRIVATE_KEY_PATH")
	owner := os.Getenv("GITHUB_OWNER")

	platform, err := platforms.NewGithub(&platforms.GithubConfig{
		AppID:          githubAppID,
		InstallationID: githubInstallationID,
		PrivateKeyPath: githubPrivateKeyPath,
	})

	if err != nil {
		return
	}

	runTests(t, platform, owner)
}
