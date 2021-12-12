//go:build integration || integrationgitlab

package tests

import (
	"os"
	"testing"

	"github.com/fikaworks/grgate/pkg/platforms"
)

func TestGitLabReleases(t *testing.T) {
	token := os.Getenv("GITLAB_TOKEN")
	owner := os.Getenv("GITLAB_OWNER")

	platform, err := platforms.NewGitlab(&platforms.GitlabConfig{
		Token: token,
	})

	if err != nil {
		return
	}

	runTests(t, platform, owner)
}
