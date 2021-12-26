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
			t.Errorf("Error listing releases: %#v", err)
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
			t.Errorf("Error listing draft releases: %#v", err)
		}
		if diff := pretty.Compare(result, expected); diff != "" {
			t.Errorf("diff: (-got +want)\n%s", diff)
		}
	})
}

func TestGithubCheckAllStatusSucceeded(t *testing.T) {
	testCases := []struct {
		name      string
		checkRuns []*github.CheckRun
		statuses  []string
		expected  bool
	}{
		{
			name: "should return true if all required check runs completed and conclusion set to success",
			checkRuns: []*github.CheckRun{
				{
					Name:       github.String("happy flow"),
					Status:     github.String("completed"),
					Conclusion: github.String("success"),
				},
				{
					Name:       github.String("no required check run"),
					Status:     github.String("pending"),
					Conclusion: github.String(""),
				},
				{
					Name:       github.String("feature B"),
					Status:     github.String("completed"),
					Conclusion: github.String("success"),
				},
			},
			statuses: []string{"happy flow", "feature B"},
			expected: true,
		},
		{
			name: "should return false if all check runs completed and not all conclusion are set to success",
			checkRuns: []*github.CheckRun{
				{
					Name:       github.String("happy flow"),
					Status:     github.String("completed"),
					Conclusion: github.String("success"),
				},
				{
					Name:       github.String("feature A"),
					Status:     github.String("completed"),
					Conclusion: github.String("skipped"),
				},
				{
					Name:       github.String("feature B"),
					Status:     github.String("completed"),
					Conclusion: github.String("failure"),
				},
			},
			statuses: []string{"happy flow", "feature A", "feature B"},
			expected: false,
		},
		{
			name: "should return false if not all check runs completed",
			checkRuns: []*github.CheckRun{
				{
					Name:       github.String("happy flow"),
					Status:     github.String("completed"),
					Conclusion: github.String("success"),
				},
				{
					Name:       github.String("feature A"),
					Status:     github.String("pending"),
					Conclusion: github.String(""),
				},
				{
					Name:       github.String("feature B"),
					Status:     github.String("completed"),
					Conclusion: github.String("success"),
				},
			},
			statuses: []string{"happy flow", "feature A", "feature B"},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatch(
					mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
					github.ListCheckRunsResults{CheckRuns: testCase.checkRuns},
				),
			)

			gh := githubPlatform{
				client:  github.NewClient(mockedHTTPClient),
				context: context.Background(),
			}

			result, err := gh.CheckAllStatusSucceeded("a", "a", "a", testCase.statuses)
			if err != nil {
				t.Errorf("Error checking status check: %#v", err)
			}
			if result != testCase.expected {
				t.Errorf("Expected %t, got %t", testCase.expected, result)
			}
		})
	}
}

func TestGithubListStatuses(t *testing.T) {
	t.Run("should list statuses", func(t *testing.T) {
		expected := []*Status{
			{
				CommitSha: "abcd1234",
				Name:      "happy flow",
				Status:    "completed",
				State:     "success",
			},
			{
				CommitSha: "abcd1234",
				Name:      "feature A",
				Status:    "queued",
				State:     "",
			},
		}

		mockedHTTPClient := mock.NewMockedHTTPClient(
			mock.WithRequestMatch(
				mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
				github.ListCheckRunsResults{CheckRuns: []*github.CheckRun{
					{
						HeadSHA:    github.String("abcd1234"),
						Name:       github.String("happy flow"),
						Status:     github.String("completed"),
						Conclusion: github.String("success"),
					},
					{
						HeadSHA: github.String("abcd1234"),
						Name:    github.String("feature A"),
						Status:  github.String("queued"),
					},
				},
				},
			),
		)

		gh := &githubPlatform{
			client:  github.NewClient(mockedHTTPClient),
			context: context.Background(),
		}

		result, err := gh.ListStatuses("a", "a", "a")
		if err != nil {
			t.Errorf("Error listing statuses: %#v", err)
		}
		if diff := pretty.Compare(result, expected); diff != "" {
			t.Errorf("diff: (-got +want)\n%s", diff)
		}
	})
}
