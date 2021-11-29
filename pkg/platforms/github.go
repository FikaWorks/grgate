package platforms

import (
	"context"
	"io"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v41/github"
)

const (
	// number of items per page to retrieve via the Github API
	githubPerPage int = 100
)

// GithubConfig hold the Github configuration
type GithubConfig struct {
	AppID          int64
	InstallationID int64
	PrivateKeyPath string
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

// ReadFile retrieve file located at the provided path in a given Github repository
func (p *githubPlatform) ReadFile(owner, repository, path string) (content io.Reader, err error) {
	content, _, err = p.client.Repositories.DownloadContents(p.context, owner,
		repository, path, nil)
	return
}

// ListReleases from a Github repository
// Important: information about published releases are available to everyone.
// Only users with push access will receive listings for draft releases.
func (p *githubPlatform) ListReleases(owner, repository string) (releases []*Release, err error) {
	opts := &github.ListOptions{
		Page:    0,
		PerPage: githubPerPage,
	}

	for {
		releaseList, resp, err := p.client.Repositories.ListReleases(p.context,
			owner, repository, opts)
		if err != nil {
			return nil, err
		}

		// for all statuses, check if the provided one are all successful
		for _, release := range releaseList {
			id := *release.ID
			tag := *release.TagName
			name := *release.Name
			commit := *release.TargetCommitish
			releaseNote := *release.Body

			// TODO: if target commitish is branch, then get lastest commit from
			// branch
			if *release.Draft {
				releases = append(releases, &Release{
					ID:          id,
					CommitSha:   commit,
					Name:        name,
					Tag:         tag,
					ReleaseNote: releaseNote,
					Platform:    "github",
				})
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return releases, err
}

// UpdateRelease edit a release based on a provided releases ID and release note
func (p *githubPlatform) UpdateRelease(owner, repository string, id interface{}, releaseNote string) (err error) {
	release, _, err := p.client.Repositories.GetRelease(p.context, owner,
		repository, id.(int64))
	if err != nil {
		return
	}

	release.Body = github.String(releaseNote)

	_, _, err = p.client.Repositories.EditRelease(p.context, owner, repository,
		id.(int64), release)
	if err != nil {
		return
	}

	return
}

// PublishRelease publish a release based on a provided releases ID
func (p *githubPlatform) PublishRelease(owner, repository string, id interface{}) (published bool, err error) {
	release, _, err := p.client.Repositories.GetRelease(p.context, owner,
		repository, id.(int64))
	if err != nil {
		return
	}

	release.Draft = github.Bool(false)

	_, _, err = p.client.Repositories.EditRelease(p.context, owner, repository,
		id.(int64), release)
	if err != nil {
		return
	}

	published = true
	return
}

// CheckAllStatusSucceeded checks that all the provided statuses succeeded
func (p *githubPlatform) CheckAllStatusSucceeded(owner, repository,
	commitSha string, statuses []string) (succeeded bool, err error) {
	if len(statuses) == 0 {
		return
	}

	opts := &github.ListCheckRunsOptions{
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: githubPerPage,
		},
	}

	for {
		getCheckRun, resp, err := p.client.Checks.ListCheckRunsForRef(p.context,
			owner, repository, commitSha, opts)
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

		opts.ListOptions.Page = resp.NextPage
	}

	return succeeded, err
}

// GetStatus returns the status of a specific commit matching a provided status name
func (p *githubPlatform) CreateStatus(owner, repository string, status *Status) (err error) {
	opts := github.CreateCheckRunOptions{
		Name:      status.Name,
		HeadSHA:   status.CommitSha,
		Status:    github.String(status.Status),
		StartedAt: &github.Timestamp{},
	}

	if status.State != "" {
		opts.Conclusion = github.String(status.State)
	}

	_, _, err = p.client.Checks.CreateCheckRun(p.context, owner, repository, opts)

	return
}

// GetStatus from provided commit and status name
func (p *githubPlatform) GetStatus(owner, repository, commitSha, statusName string) (status *Status, err error) {
	statusList, err := p.ListStatuses(owner, repository, commitSha)
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

// ListStatuses attached to a given commit sha
func (p *githubPlatform) ListStatuses(owner, repository, commitSha string) (statusList []*Status, err error) {
	opts := &github.ListCheckRunsOptions{
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: githubPerPage,
		},
	}

	for {
		getCheckRun, resp, err := p.client.Checks.ListCheckRunsForRef(p.context,
			owner, repository, commitSha, opts)
		if err != nil {
			return nil, err
		}

		for _, checkRun := range getCheckRun.CheckRuns {
			cr := &Status{
				CommitSha: *checkRun.HeadSHA,
				Name:      *checkRun.Name,
				Status:    *checkRun.Status,
			}

			if checkRun.Conclusion != nil {
				cr.State = *checkRun.Conclusion
			}
			statusList = append(statusList, cr)
		}

		if resp.NextPage == 0 {
			break
		}

		opts.ListOptions.Page = resp.NextPage
	}

	return statusList, err
}
