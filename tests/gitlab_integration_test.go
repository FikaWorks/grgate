//go:build integration || integrationgitlab

package tests

import (
	"os"
	"testing"

	"github.com/fikaworks/grgate/pkg/platforms"
)

func TestGitLabReleases(t *testing.T) {
	author := os.Getenv("GITLAB_AUTHOR")
	owner := os.Getenv("GITLAB_OWNER")
	token := os.Getenv("GITLAB_TOKEN")

	platform, err := platforms.NewGitlab(&platforms.GitlabConfig{
		Token: token,
	})

	if err != nil {
		return
	}

	runTests(t, platform, owner, author)
}
