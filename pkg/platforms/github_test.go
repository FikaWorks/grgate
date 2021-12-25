//go:build unit

package platforms

import (
	"context"
	"testing"

	"github.com/google/go-github/v41/github"
	"github.com/kylelemons/godebug/pretty"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestGithubListReleases(t *testing.T) {
	t.Run("should list releases", func(t *testing.T) {
		expected := []*Release{
			{
				CommitSha: "master",
				Draft:     true,
				ID:        123,
				Name:      "draft",
				Platform:  "github",
				Tag:       "v1.2.3",
			},
			{
				CommitSha:   "master",
				Draft:       false,
				ID:          456,
				Name:        "published",
				Platform:    "github",
				Tag:         "v1.2.3",
				ReleaseNote: "release note",
			},
		}

		mockedHTTPClient := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposReleasesByOwnerByRepo,
				[]*github.RepositoryRelease{
					{
						ID:              github.Int64(123),
						Draft:           github.Bool(true),
						Name:            github.String("draft"),
						TagName:         github.String("v1.2.3"),
						TargetCommitish: github.String("master"),
					},
					{
						ID:              github.Int64(456),
						Body:            github.String("release note"),
						Draft:           github.Bool(false),
						Name:            github.String("published"),
						TagName:         github.String("v1.2.3"),
						TargetCommitish: github.String("master"),
					},
				},
			),
		)

		gh := &githubPlatform{
			client:  github.NewClient(mockedHTTPClient),
			context: context.Background(),
		}

		result, err := gh.ListReleases("a", "a")
		if err != nil {
			t.Error("Error listing draft releases", err)
		}
		if diff := pretty.Compare(result, expected); diff != "" {
			t.Errorf("diff: (-got +want)\n%s", diff)
		}
	})
}

func TestGithubListDraftReleases(t *testing.T) {
	t.Run("should list draft releases", func(t *testing.T) {
		expected := []*Release{
			{
				CommitSha: "master",
				Draft:     true,
				ID:        123,
				Name:      "draft",
				Platform:  "github",
				Tag:       "v1.2.3",
			},
		}

		mockedHTTPClient := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposReleasesByOwnerByRepo,
				[]*github.RepositoryRelease{
					{
						Draft:           github.Bool(true),
						ID:              github.Int64(123),
						Name:            github.String("draft"),
						TagName:         github.String("v1.2.3"),
						TargetCommitish: github.String("master"),
					},
					{
						Draft:           github.Bool(false),
						ID:              github.Int64(456),
						Name:            github.String("published"),
						TagName:         github.String("v1.2.3"),
						TargetCommitish: github.String("master"),
					},
				},
			),
		)

		gh := &githubPlatform{
			client:  github.NewClient(mockedHTTPClient),
			context: context.Background(),
		}

		result, err := gh.ListDraftReleases("a", "a")
		if err != nil {
			t.Error("Error listing draft releases", err)
		}
		if diff := pretty.Compare(result, expected); diff != "" {
			t.Errorf("diff: (-got +want)\n%s", diff)
		}
	})
}
