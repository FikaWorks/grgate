// +build integration integrationgitlab

package tests

import (
	"math/rand"
	"time"
	"os"
  "fmt"
  "testing"

	"github.com/xanzy/go-gitlab"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/fikaworks/grgate/pkg/workers"
)

type gitlabTest struct {
	client *gitlab.Client
	repository string
	pid string
}

func newGitlabTest(token, owner, repositoryPrefix string) (integrationTest integrationTest,
	client *gitlab.Client, pid string, repository string, err error) {
  rand.Seed(time.Now().UnixNano())

	repository = fmt.Sprintf("%s-%s", repositoryPrefix, randomString(5))
	pid = fmt.Sprintf("%s/%s", owner, repository)

  client, err = gitlab.NewClient(token)
	if err != nil {
		return
	}

	integrationTest = &gitlabTest{
		client: client,
		repository: repository,
		pid: pid,
	}

	return
}

func (g *gitlabTest) setup(rc *repoConfig, rr *repoRelease) (err error) {
  opts := &gitlab.CreateProjectOptions{
    Name: &g.repository,
		Visibility: gitlab.Visibility(gitlab.VisibilityValue("private")),
  }
	fmt.Printf("Creating repository %s\n", g.pid)
  // TODO: flakky test, might need to wait for project to be created before continuing
  _, _, err = g.client.Projects.CreateProject(opts, nil)

	fmt.Printf("Creating repository config %s\n", rc.path)
  fileOpts := &gitlab.CreateFileOptions{
    Branch: &rc.branch,
    Content: &rc.content,
    CommitMessage: &rc.commitMessage,
  }
  _, _, err = g.client.RepositoryFiles.CreateFile(g.pid, rc.path, fileOpts, nil)
  if err != nil {
    return
  }

	fmt.Printf("Creating release %s\n", rr.name)
  releaseOpts := &gitlab.CreateReleaseOptions{
    Name: &rr.name,
    Ref: &rr.ref,
    TagName: &rr.tag,
    ReleasedAt: &rr.releasedAt,
  }
  _, _, err = g.client.Releases.CreateRelease(g.pid, releaseOpts, nil)

	return
}

func (g *gitlabTest) teardown() error {
	fmt.Printf("Deleting repository %s\n", g.pid)
  _, err := g.client.Projects.DeleteProject(g.pid, nil)
	return err
}

func TestGitlabReleases(t *testing.T) {
  token := os.Getenv("GITLAB_TOKEN")
  owner := os.Getenv("GITLAB_OWNER")

  gt, client, pid, repository, err := newGitlabTest(token, owner,
    repositoryPrefix)
  if err != nil {
    t.Error("Couldn't setup test", err)
    return
  }

  config.Main = &config.MainConfig{
    RepoConfigPath: ".grgate.yaml",
    Globals: &config.RepoConfig{
      Enabled: true,
      TagRegexp: ".*",
      Statuses: []string{
        "e2e-happyflow",
      },
    },
  }

  rc := &repoConfig{
    branch: "master",
    commitMessage: "init",
    content: "content",
    path: "test-file.yaml",
  }
  rr := &repoRelease{
    name: "test release",
    ref: "master",
    tag: "v1.2.3",
    releasedAt: time.Now().UTC().Add(time.Hour),
  }

  err = gt.setup(rc, rr)
  if err != nil {
    t.Error("Couldn't setup test", err)
    return
  }

  // delete gitlab repository after tests
  defer gt.teardown()

  platform, err := platforms.NewGitlab(&platforms.GitlabConfig{
    Token: token,
  })
  if err != nil {
    t.Error("Couldn't create gitlab client", err)
    return
  }

  job, err := workers.NewJob(platform, owner, repository)
  if err != nil {
    t.Error("Couldn't create job", err)
    return
  }

  release, _, err := client.Releases.GetRelease(pid, rr.tag, nil)
  if err != nil {
    t.Error("Couldn't get release", err)
    return
  }

	t.Run("should not publish release when commit status are not defined", func(t *testing.T) {
		err = job.Process()
		if err != nil {
			t.Error("Couldn't process repository", err)
      return
		}

		release, _, err := client.Releases.GetRelease(pid, rr.tag, nil)
		if err != nil {
			t.Error("Couldn't get release", err)
			return
		}

		if release.ReleasedAt.Before(time.Now().UTC()) {
			t.Error("Release should not be published when status is not set")
			return
		}
  })

	t.Run("should not publish release when some commit status are still running", func(t *testing.T) {
		err = platform.CreateStatus(owner, repository, &platforms.Status{
			Name: "e2e-happyflow",
			Status: "running",
			CommitSha: release.Commit.ID,
		})
		if err != nil {
			t.Error("Couldn't set running status to commit", err)
			return
		}

		err = job.Process()
		if err != nil {
			t.Error("Couldn't process repository", err)
			return
		}

		release, _, err = client.Releases.GetRelease(pid, rr.tag, nil)
		if err != nil {
			t.Error("Couldn't get release", err)
			return
		}

		if release.ReleasedAt.Before(time.Now().UTC()) {
			t.Error("Release should not be published when the commit status are still running")
			return
		}
  })

	t.Run("should publish release if all status succeeded", func(t *testing.T) {
		err = platform.CreateStatus(owner, repository, &platforms.Status{
			Name: "e2e-happyflow",
			Status: "success",
			CommitSha: release.Commit.ID,
		})
		if err != nil {
			t.Error("Couldn't set success status to commit", err)
			return
		}

		err = job.Process()
		if err != nil {
			t.Error("Couldn't process repository", err)
			return
		}

		release, _, err = client.Releases.GetRelease(pid, rr.tag, nil)
		if err != nil {
			t.Error("Couldn't get release", err)
			return
		}

		if release.ReleasedAt.After(time.Now().UTC()) {
      t.Error("Release wasn't published after all commit status succeeded", err)
			return
		}
  })
}
