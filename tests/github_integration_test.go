// +build integration integrationgithub

package tests

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"time"
  "fmt"
  "strconv"
  "testing"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v34/github"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/fikaworks/grgate/pkg/workers"
)

type githubTest struct {
	client *github.Client
  owner string
	repository string
	context context.Context
}

func newGithubTest(appID, installationID int64, privateKeyPath, owner, repositoryPrefix string) (integrationTest integrationTest,
	client *github.Client, ctx context.Context, repository string, err error) {
  rand.Seed(time.Now().UnixNano())

	repository = fmt.Sprintf("%s-%s", repositoryPrefix, randomString(5))

	ctx = context.Background()

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID,
    installationID, privateKeyPath)
	if err != nil {
		return
	}

  client = github.NewClient(&http.Client{Transport: itr})

	integrationTest = &githubTest{
		client: client,
		repository: repository,
		owner: owner,
    context: ctx,
	}

	return
}

func (g *githubTest) setup(rc *repoConfig, rr *repoRelease) (err error) {
  opts := &github.Repository{
    Name: &g.repository,
    Visibility: github.String("private"),
  }
	fmt.Printf("Creating repository %s\n", g.repository)
  _, _, err = g.client.Repositories.Create(g.context, g.owner, opts)

	fmt.Printf("Creating repository config %s\n", rc.path)
  fileOpts := &github.RepositoryContentFileOptions{
    Branch: github.String(rc.branch),
    Content: []byte(rc.content),
    Message: github.String(rc.commitMessage),
  }
  _, _, err = g.client.Repositories.CreateFile(g.context, g.owner,
    g.repository, rc.path, fileOpts)
  if err != nil {
    return
  }

	fmt.Printf("Creating release %s\n", rr.name)
  releaseOpts := &github.RepositoryRelease{
    Name: github.String(rr.name),
    TargetCommitish: github.String(rr.ref),
    TagName: github.String(rr.tag),
    Draft: github.Bool(rr.draft),
  }
  _, _, err = g.client.Repositories.CreateRelease(g.context, g.owner,
    g.repository, releaseOpts)

	return
}

func (g *githubTest) teardown() error {
	fmt.Printf("Deleting repository %s\n", g.repository)
  _, err := g.client.Repositories.Delete(g.context, g.owner, g.repository)
	return err
}

func TestGithubReleases(t *testing.T) {
  githubAppID, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
  githubInstallationID, _ := strconv.ParseInt(os.Getenv("GITHUB_INSTALLATION_ID"), 10, 64)
  githubPrivateKeyPath := os.Getenv("GITHUB_PRIVATE_KEY_PATH")
  owner := os.Getenv("GITHUB_OWNER")

  gt, client, ctx, repository, err := newGithubTest(githubAppID,
    githubInstallationID, githubPrivateKeyPath, owner,
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
    draft: true,
  }

  err = gt.setup(rc, rr)
  if err != nil {
    t.Error("Couldn't setup test", err)
    return
  }

  // delete github repository after tests
  defer gt.teardown()

  platform, err := platforms.NewGithub(&platforms.GithubConfig{
    AppID: githubAppID,
    InstallationID: githubInstallationID,
    PrivateKeyPath: githubPrivateKeyPath,
  })
  if err != nil {
    t.Error("Couldn't create github client", err)
    return
  }

  job, err := workers.NewJob(platform, owner, repository)
  if err != nil {
    t.Error("Couldn't create job", err)
    return
  }

	t.Run("should not publish release when commit status are not defined", func(t *testing.T) {
		err = job.Process()
		if err != nil {
			t.Error("Couldn't process repository", err)
      return
		}

		releaseList, _, err := client.Repositories.ListReleases(ctx, owner,
      repository, nil)
		if err != nil {
			t.Error("Couldn't get releases", err)
			return
		}

		// check that release hasn't been published, if still draft then the test
    // is successful
		for _, release := range releaseList {
			if *release.TagName == rr.tag && *release.Draft {
        return
      }
		}

    t.Error("Release should not be published when status is not set")
  })

	t.Run("should not publish release when some commit status are still running", func(t *testing.T) {
		err = platform.CreateStatus(owner, repository, &platforms.Status{
			Name: "e2e-happyflow",
			Status: "in_progress",
			CommitSha: "master",
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

		releaseList, _, err := client.Repositories.ListReleases(ctx, owner,
      repository, nil)
		if err != nil {
			t.Error("Couldn't get releases", err)
			return
		}

		// check that release hasn't been published, if still draft then the test
    // is successful
		for _, release := range releaseList {
			if *release.TagName == rr.tag && *release.Draft {
        return
      }
		}

    t.Error("Release should not be published when the commit status are still running")
  })

	t.Run("should publish release if all status succeeded", func(t *testing.T) {
		err = platform.CreateStatus(owner, repository, &platforms.Status{
			Name: "e2e-happyflow",
			State: "success",
			Status: "completed",
			CommitSha: "master",
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

		releaseList, _, err := client.Repositories.ListReleases(ctx, owner,
      repository, nil)
		if err != nil {
			t.Error("Couldn't get releases", err)
			return
		}

		// check that release has been published, if still draft then test failed
		for _, release := range releaseList {
			if *release.TagName == rr.tag && !*release.Draft {
        return
      }
		}

    t.Error("Release wasn't published after all commit status succeeded", err)
  })
}
