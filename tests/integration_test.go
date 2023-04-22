//go:build integration || integrationgithub || integrationgitlab

package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/fikaworks/grgate/pkg/workers"
)

const repositoryPrefix = "grgate-integration"

func runTests(t *testing.T, platform platforms.Platform, owner, author string) {
	if _, err := config.NewGlobalConfig(""); err != nil {
		t.Errorf("Error not expected: %#v", err)
	}

	// force set author in order to look for issues by author during validation steps
	config.Main.Globals.Dashboard.Author = author

	runTest(t, platform, owner, disabledConfigTestCases)
	runTest(t, platform, owner, commitStatusTestCases)
	runTest(t, platform, owner, releaseNoteTestCases)
	runTest(t, platform, owner, dashboardTestCases)
}

func setup(platform platforms.Platform, owner string) (repository string, err error) {
	repository = generateRandomRepositoryName(repositoryPrefix)

	fmt.Printf("Creating repository %s/%s\n", owner, repository)
	err = platform.CreateRepository(owner, repository, "private")
	return
}

func tearDown(platform platforms.Platform, owner, repository string) {
	_ = platform.DeleteRepository(owner, repository)
}

// runTests prepare a repository and run GRGate against it
func runTest(t *testing.T, platform platforms.Platform, owner string, testCases map[string]*testCase) {
	for title, testCase := range testCases {
		t.Run(title, func(t *testing.T) {
			repository, err := setup(platform, owner)
			if err != nil {
				t.Errorf("Couldn't create repository: %#v", err)
				return
			}

			// fix flakky repository creation, it seems to have inconsistent delay
			time.Sleep(time.Second)

			defer tearDown(platform, owner, repository)

			if err := platform.CreateFile(owner, repository, ".grgate.yaml",
				"master", "init", testCase.withRepoConfig); err != nil {
				t.Errorf("Couldn't create file: %#v", err)
				return
			}

			// fix flakky CreateFile, it seems to have inconsistent delay
			time.Sleep(time.Second)

			release, err := platform.CreateRelease(owner, repository, &platforms.Release{
				CommitSha: "master",
				Tag:       testCase.withTag,
				Draft:     true,
			})
			if err != nil {
				t.Errorf("Couldn't create release: %#v", err)
				return
			}

			// fix flakky CreateRelease, it seems to have inconsistent delay
			time.Sleep(time.Second)

			job, err := workers.NewJob(platform, owner, repository)
			if err != nil {
				t.Errorf("Couldn't create job: %#v", err)
				return
			}

			for _, status := range testCase.withStatuses {
				status.CommitSha = release.CommitSha
				if err := platform.CreateStatus(owner, repository, status); err != nil {
					t.Errorf("Couldn't set status named %s to state %s and status %s : %#v",
						status.Name, status.State, status.Status, err)
					return
				}
			}

			// fix flakky CreateStatus, it seems to have inconsistent delay
			time.Sleep(time.Second)

			if err := job.Process(); err != nil && !testCase.expectErrorDuringProcess {
				t.Errorf("Couldn't process repository: %#v", err)
				return
			}

			// fix flakky Process, it seems to have inconsistent delay
			time.Sleep(time.Second)

			// validate issue dashboard
			issueList, err := platform.ListIssuesByAuthor(owner, repository, config.Main.Globals.Dashboard.Author)
			if err != nil {
				t.Errorf("Couldn't list issues from repository: %#v", err)
				return
			}

			issueExist := false
			for _, issue := range issueList {
				if issue.Title == testCase.expectedDashboardTitle {
					if diff := pretty.Compare(issue.Body, testCase.expectedDashboardBody); diff != "" {
						t.Errorf("diff: (-got +want)\n%s", diff)
						return
					}
					issueExist = true
					break
				}
			}

			if testCase.expectIssueToBeCreated && !issueExist {
				t.Errorf("Issue dashboard was not found in repository %s/%s", owner, repository)
				return
			}

			// validate release status
			releaseList, err := platform.ListReleases(owner, repository)
			if err != nil {
				t.Errorf("Couldn't list releases from repository: %#v", err)
				return
			}

			for _, release := range releaseList {
				if release.Tag == testCase.withTag {
					diff := pretty.Compare(release.ReleaseNote, testCase.expectedReleaseNote)
					if diff != "" && testCase.expectedReleaseNote != "" {
						t.Errorf("diff: (-got +want)\n%s", diff)
						return
					}

					if testCase.expectPublishedRelease == !release.Draft {
						return
					}

					t.Error("Expected release to be published")
					return
				}
			}

			t.Error("Release was not found")
		})
	}
}
