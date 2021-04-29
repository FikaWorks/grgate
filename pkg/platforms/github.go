package platforms

import (
	"context"
	"io"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v34/github"
)

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

func NewGithub(config *GithubConfig) (platform Platform, err error) {
	ctx := context.Background()

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, config.AppID,
		config.InstallationID, config.PrivateKeyPath)
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

func (p *githubPlatform) ReadFile(owner, repository, path string) (content io.ReadCloser, err error) {
	content, _, err = p.client.Repositories.DownloadContents(p.context, owner,
    repository, path, nil)
	return
}

func (p *githubPlatform) ListReleases(owner, repository string) (releases []*Release, err error) {
	// opt := &github.ListOptions{Page: 2}
	// Information about published releases are available to everyone. Only users with push access will receive listings for draft releases.
	releaseList, _, err := p.client.Repositories.ListReleases(p.context, owner,
    repository, nil)
	if err != nil {
		return
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

	return
}

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

func (p *githubPlatform) CheckAllStatusSucceeded(owner, repository, commitSha string, statusNameList []string) (succeeded bool, err error) {
	if len(statusNameList) == 0 {
		return
	}

	// TODO: paginate
	getCheckRun, _, err := p.client.Checks.ListCheckRunsForRef(p.context, owner,
    repository, commitSha, nil)
	if err != nil {
		return
	}

	succeededCheck := 0
	// TODO: make sure all values in statusNameList are unique
	for _, check := range getCheckRun.CheckRuns {
		for _, status := range statusNameList {
			if *check.Name == status && check.Conclusion != nil && *check.Conclusion == "success" {
				succeededCheck++
			}
		}
	}

	succeeded = succeededCheck == len(statusNameList)
	return
}

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

func (p *githubPlatform) ListStatus(owner, repository, commitSha string) (statusList []*Status, err error) {
	// TODO: paginate
	getCheckRun, _, err := p.client.Checks.ListCheckRunsForRef(p.context, owner,
		repository, commitSha, nil)
	if err != nil {
		return
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

	return
}
