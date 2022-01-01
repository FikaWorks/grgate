//go:build integration || integrationgithub || integrationgitlab

package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/kylelemons/godebug/pretty"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/fikaworks/grgate/pkg/workers"
)

const (
	repositoryPrefix = "grgate-integration"
	hashSize         = 5
)

func runTests(t *testing.T, platform platforms.Platform, owner string) {
	if _, err := config.NewGlobalConfig(""); err != nil {
		t.Errorf("Error not expected: %#v", err)
	}

	testDisabledConfig(t, platform, owner)
	testCommitStatus(t, platform, owner)
	testReleaseNote(t, platform, owner)
}

func setup(platform platforms.Platform, owner string) (repository string, err error) {
	rand.Seed(time.Now().UnixNano())
	repository = generateRandomRepositoryName(repositoryPrefix)
	fmt.Printf("Creating repository %s/%s\n", owner, repository)
	err = platform.CreateRepository(owner, repository, "private")
	return
}

func tearDown(platform platforms.Platform, owner, repository string) {
	_ = platform.DeleteRepository(owner, repository)
}

func testDisabledConfig(t *testing.T, platform platforms.Platform, owner string) {
	repoConfig := `enabled: false
statuses:
- e2e-happyflow`
	tag := "v1.2.3"

	repository, err := setup(platform, owner)
	if err != nil {
		t.Errorf("Couldn't create repository: %#v", err)
		return
	}
	defer tearDown(platform, owner, repository)

	if err := platform.CreateFile(owner, repository, ".grgate.yaml", "master", "init", repoConfig); err != nil {
		t.Errorf("Couldn't create file: %#v", err)
		return
	}

	job, err := workers.NewJob(platform, owner, repository)
	if err != nil {
		t.Errorf("Couldn't create job: %#v", err)
		return
	}

	currentRelease, err := platform.CreateRelease(owner, repository, &platforms.Release{
		CommitSha: "master",
		Tag:       tag,
		Draft:     true,
	})
	if err != nil {
		t.Errorf("Couldn't create release: %#v", err)
		return
	}

	if err := platform.CreateStatus(owner, repository, &platforms.Status{
		Name:      "e2e-happyflow",
		State:     "success",
		Status:    "completed",
		CommitSha: currentRelease.CommitSha,
	}); err != nil {
		t.Errorf("Couldn't set running status to commit: %#v", err)
		return
	}

	t.Run("should not process the release when disabled by config", func(t *testing.T) {
		if err := job.Process(); err != nil {
			t.Errorf("Couldn't process repository: %#v", err)
			return
		}

		releaseList, err := platform.ListReleases(owner, repository)
		if err != nil {
			t.Errorf("Couldn't list releases from repository: %#v", err)
			return
		}

		// check that release hasn't been published, if still draft then the test
		// is successful
		for _, release := range releaseList {
			if release.Tag == tag && release.Draft {
				return
			}
		}

		t.Error("Release should not be published when disabled by config")
	})
}

func testCommitStatus(t *testing.T, platform platforms.Platform, owner string) {
	repoConfig := `enabled: true
tagRegexp: v\d*\.\d*\.\d*-beta\.\d*
statuses:
- e2e-happyflow
- e2e-featureflow`

	tag := "v1.2.3-beta.0"

	repository, err := setup(platform, owner)
	if err != nil {
		t.Errorf("Couldn't create repository: %#v", err)
		return
	}
	defer tearDown(platform, owner, repository)

	if err := platform.CreateFile(owner, repository, ".grgate.yaml", "master", "init", repoConfig); err != nil {
		t.Errorf("Couldn't create file: %#v", err)
		return
	}

	job, err := workers.NewJob(platform, owner, repository)
	if err != nil {
		t.Errorf("Couldn't create job: %#v", err)
		return
	}

	currentRelease, err := platform.CreateRelease(owner, repository, &platforms.Release{
		CommitSha: "master",
		Tag:       tag,
		Draft:     true,
	})
	if err != nil {
		t.Errorf("Couldn't create release: %#v", err)
		return
	}

	t.Run("should not publish release when commit status are not defined", func(t *testing.T) {
		if err := job.Process(); err != nil {
			t.Errorf("Couldn't process repository: %#v", err)
			return
		}

		releaseList, err := platform.ListReleases(owner, repository)
		if err != nil {
			t.Errorf("Couldn't list releases from repository: %#v", err)
			return
		}

		// check that release hasn't been published, if still draft then the test
		// is successful
		for _, release := range releaseList {
			if release.Tag == tag && release.Draft {
				return
			}
		}

		t.Error("Release should not be published when status is not set")
	})

	t.Run("should not publish release when some commit status are still running", func(t *testing.T) {
		if err := platform.CreateStatus(owner, repository, &platforms.Status{
			Name:      "e2e-happyflow",
			Status:    "in_progress",
			CommitSha: currentRelease.CommitSha,
		}); err != nil {
			t.Errorf("Couldn't set running status to commit: %#v", err)
			return
		}

		if err := job.Process(); err != nil {
			t.Errorf("Couldn't process repository: %#v", err)
			return
		}

		releaseList, err := platform.ListReleases(owner, repository)
		if err != nil {
			t.Errorf("Couldn't list releases from repository: %#v", err)
			return
		}

		// check that release hasn't been published, if still draft then the test
		// is successful
		for _, release := range releaseList {
			if release.Tag == tag && release.Draft {
				return
			}
		}

		t.Error("Release should not be published when the commit status are still running")
	})

	t.Run("should publish release if all status succeeded", func(t *testing.T) {
		if err := platform.CreateStatus(owner, repository, &platforms.Status{
			Name:      "e2e-happyflow",
			State:     "success",
			Status:    "completed",
			CommitSha: currentRelease.CommitSha,
		}); err != nil {
			t.Errorf("Couldn't set success status to commit: %#v", err)
			return
		}

		if err := platform.CreateStatus(owner, repository, &platforms.Status{
			Name:      "e2e-featureflow",
			State:     "success",
			Status:    "completed",
			CommitSha: currentRelease.CommitSha,
		}); err != nil {
			t.Errorf("Couldn't set success status to commit: %#v", err)
			return
		}

		if err := job.Process(); err != nil {
			t.Errorf("Couldn't process repository: %#v", err)
			return
		}

		releaseList, err := platform.ListReleases(owner, repository)
		if err != nil {
			t.Errorf("Couldn't list releases from repository: %#v", err)
			return
		}

		// check that release hasn't been published, if still draft then the test
		// is successful
		for _, release := range releaseList {
			if release.Tag == tag && !release.Draft {
				return
			}
		}

		t.Error("Release wasn't published after all commit status succeeded")
	})
}

func testReleaseNote(t *testing.T, platform platforms.Platform, owner string) {
	repoConfig := `enabled: true
tagRegexp: v\d*\.\d*\.\d*-beta\.\d*
releaseNote:
  enabled: true
  template: |-
    {{- .ReleaseNote -}}
    <!-- GRGate start -->
    <details><summary>GRGate status check</summary>
    {{ range .Statuses }}
    - [{{ if or (eq .Status "completed" ) (eq .Status "success") }}x{{ else }} {{ end }}] {{ .Name }}
    {{- end }}

    </details>
    <!-- GRGate end -->
statuses:
- e2e-happyflow
- e2e-featureflow-a
- e2e-featureflow-b`

	tag := "v1.2.3-beta.1"

	repository, err := setup(platform, owner)
	if err != nil {
		t.Errorf("Couldn't create repository: %#v", err)
		return
	}
	defer tearDown(platform, owner, repository)

	if err := platform.CreateFile(owner, repository, ".grgate.yaml", "master", "init", repoConfig); err != nil {
		t.Errorf("Couldn't create file: %#v", err)
		return
	}

	job, err := workers.NewJob(platform, owner, repository)
	if err != nil {
		t.Errorf("Couldn't create job: %#v", err)
		return
	}

	currentRelease, err := platform.CreateRelease(owner, repository, &platforms.Release{
		CommitSha: "master",
		Tag:       tag,
		Draft:     true,
	})
	if err != nil {
		t.Errorf("Couldn't create release: %#v", err)
		return
	}

	t.Run("should update release note with statuses", func(t *testing.T) {
		expectedReleaseNote := `<!-- GRGate start -->
<details><summary>GRGate status check</summary>

- [ ] e2e-featureflow-a
- [ ] e2e-featureflow-b
- [ ] e2e-happyflow

</details>
<!-- GRGate end -->`

		if err := job.Process(); err != nil {
			t.Errorf("Couldn't process repository: %#v", err)
			return
		}

		releaseList, err := platform.ListDraftReleases(owner, repository)
		if err != nil {
			t.Errorf("Couldn't list releases from repository: %#v", err)
			return
		}

		// check that the draft release hasn't been published
		for _, release := range releaseList {
			if release.Tag == tag {
				currentRelease = release

				if diff := pretty.Compare(currentRelease.ReleaseNote, expectedReleaseNote); diff != "" {
					t.Errorf("diff: (-got +want)\n%s", diff)
					return
				}
				return
			}
		}

		t.Error("Release should not be published when statuses are not set")
	})

	t.Run("should publish release if all status succeeded", func(t *testing.T) {
		expectedReleaseNote := `<!-- GRGate start -->
<details><summary>GRGate status check</summary>

- [x] e2e-featureflow-a
- [x] e2e-featureflow-b
- [x] e2e-happyflow

</details>
<!-- GRGate end -->`

		if err := platform.CreateStatus(owner, repository, &platforms.Status{
			Name:      "e2e-happyflow",
			State:     "success",
			Status:    "completed",
			CommitSha: "master",
		}); err != nil {
			t.Errorf("Couldn't set success status to commit: %#v", err)
			return
		}

		if err := platform.CreateStatus(owner, repository, &platforms.Status{
			Name:      "e2e-featureflow-a",
			State:     "success",
			Status:    "completed",
			CommitSha: "master",
		}); err != nil {
			t.Errorf("Couldn't set success status to commit: %#v", err)
			return
		}

		if err := platform.CreateStatus(owner, repository, &platforms.Status{
			Name:      "e2e-featureflow-b",
			State:     "success",
			Status:    "completed",
			CommitSha: "master",
		}); err != nil {
			t.Errorf("Couldn't set success status to commit: %#v", err)
			return
		}

		if err := job.Process(); err != nil {
			t.Errorf("Couldn't process repository: %#v", err)
			return
		}

		releaseList, err := platform.ListReleases(owner, repository)
		if err != nil {
			t.Errorf("Couldn't list releases from repository: %#v", err)
			return
		}

		// check that release has correctly been published
		for _, release := range releaseList {
			if release.Tag == tag {
				if release.Draft {
					t.Error("Expect release to be published")
					return
				}
				if diff := pretty.Compare(release.ReleaseNote, expectedReleaseNote); diff != "" {
					t.Errorf("diff: (-got +want)\n%s", diff)
					return
				}
				return
			}
		}

		t.Error("Release wasn't published after all commit status succeeded")
	})
}
