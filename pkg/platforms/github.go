package platforms

import (
	"context"
	"io"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v34/github"
)

// GithubConfig hold the Github configuration
type GithubConfig struct {
	AppID          int64
	InstallationID int64
	PrivateKeyPath string
	WebhookSecret  string
}

type githubPlatform struct {
	config  *GithubConfig
	client  *github.Client
	context context.Context
}

// NewGithub returns an instance of platform
func NewGithub(config *GithubConfig) (platform Platform, err error) {
	ctx := context.Background()

  itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport,
    config.AppID, config.InstallationID, config.PrivateKeyPath)
	if err != nil {
		return
	}

	platform = &githubPlatform{
		config:  config,
		client:  github.NewClient(&http.Client{Transport: itr}),
		context: ctx,
	}

	return
}

// ReadFile located at the provided path in a given Github repository
func (p *githubPlatform) ReadFile(owner, repository, path string) (content io.ReadCloser, err error) {
	content, _, err = p.client.Repositories.DownloadContents(p.context, owner,
    repository, path, nil)
	return
}

// ListReleases from a Github repository
// Important: information about published releases are available to everyone.
// Only users with push access will receive listings for draft releases.
func (p *githubPlatform) ListReleases(owner, repository string) (releases []*Release, err error) {
	opt := &github.ListOptions{
    Page: 0,
    PerPage: 100,
  }

  for {
    releaseList, resp, err := p.client.Repositories.ListReleases(p.context,
      owner, repository, opt)
    if err != nil {
      return nil, err
    }

    for _, release := range releaseList {
      id := *release.ID
      tag := *release.TagName
      name := *release.Name
      commit := *release.TargetCommitish

      // TODO: if target commitish is branch, then get lastest commit from branch
      if *release.Draft {
        releases = append(releases, &Release{
          ID:        id,
          CommitSha: commit,
          Name:      name,
          Tag:       tag,
          Platform:  "github",
        })
      }
    }

    if resp.NextPage == 0 {
      break
    }

    opt.Page = resp.NextPage
  }

	return
}

// PublishRelease publish a release based on a provided releases ID
func (p *githubPlatform) PublishRelease(owner, repository string, id int64) (published bool, err error) {
	release, _, err := p.client.Repositories.GetRelease(p.context, owner,
    repository, id)
	if err != nil {
		return
	}

	release.Draft = github.Bool(false)
	p.client.Repositories.EditRelease(p.context, owner, repository, id, release)

	published = true
	return
}

// CheckAllStatusSucceeded checks that all the provided statuses succeeded
func (p *githubPlatform) CheckAllStatusSucceeded(owner, repository, commitSha string, statuses []string) (succeeded bool, err error) {
	if len(statuses) == 0 {
		return
	}

	opt := &github.ListCheckRunsOptions{
    ListOptions: github.ListOptions{
      Page: 0,
      PerPage: 100,
    },
  }

  for {
    getCheckRun, resp, err := p.client.Checks.ListCheckRunsForRef(p.context,
      owner, repository, commitSha, opt)
    if err != nil {
      return false, err
    }

    succeededCheck := 0
    // TODO: make sure all values in statuses are unique
    for _, check := range getCheckRun.CheckRuns {
      for _, status := range statuses {
        if *check.Name == status && check.Conclusion != nil && *check.Conclusion == "success" {
          succeededCheck++
        }
      }
    }

    succeeded = succeededCheck == len(statuses)

    if resp.NextPage == 0 {
      break
    }

    opt.ListOptions.Page = resp.NextPage
  }

	return
}

// CreateStatus for a given commit
func (p *githubPlatform) CreateStatus(owner, repository string, status *Status) (err error) {
	checkRunOpt := github.CreateCheckRunOptions{
		Name:      status.Name,
		HeadSHA:   status.CommitSha,
		Status:    github.String(status.Status),
		StartedAt: &github.Timestamp{},
	}

	if status.State != "" {
		checkRunOpt.Conclusion = github.String(status.State)
	}

	_, _, err = p.client.Checks.CreateCheckRun(p.context, owner, repository,
    checkRunOpt)

	return
}

// GetStatus from provided commit and status name
func (p *githubPlatform) GetStatus(owner, repository, commitSha, statusName string) (status *Status, err error) {
	statusList, err := p.ListStatus(owner, repository, commitSha)
	if err != nil {
		return
	}

	for _, cr := range statusList {
		if cr.Name == statusName {
			status = cr
			return
		}
	}

	return
}

// ListStatus from provided commit
func (p *githubPlatform) ListStatus(owner, repository, commitSha string) (statusList []*Status, err error) {
	opt := &github.ListCheckRunsOptions{
    ListOptions: github.ListOptions{
      Page: 0,
      PerPage: 100,
    },
  }

  for {
    getCheckRun, resp, err := p.client.Checks.ListCheckRunsForRef(p.context,
      owner, repository, commitSha, opt)
    if err != nil {
      return nil, err
    }

    for _, checkRun := range getCheckRun.CheckRuns {
      cr := &Status{
        Name:   *checkRun.Name,
        Status: *checkRun.Status,
      }
      if checkRun.Conclusion != nil {
        cr.State = *checkRun.Conclusion
      }
      statusList = append(statusList, cr)
    }

    if resp.NextPage == 0 {
      break
    }

    opt.ListOptions.Page = resp.NextPage
  }

	return
}
