package platforms

import (
	"context"
	"io"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
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
	Token          string
}

type githubPlatform struct {
	config  *GithubConfig
	client  *github.Client
	context context.Context
}

// NewGithub returns an instance of platform
func NewGithub(config *GithubConfig) (Platform, error) {
	ctx := context.Background()

	platform := &githubPlatform{
		config:  config,
		context: ctx,
	}

	// GitHub private key or token based authentication
	if config.Token == "" {
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport,
			config.AppID, config.InstallationID, config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}

		platform.client = github.NewClient(&http.Client{Transport: itr})
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.Token},
		)
		tc := oauth2.NewClient(ctx, ts)

		platform.client = github.NewClient(tc)
	}

	return platform, nil
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
			draft := *release.Draft

			var releaseNote string
			if release.Body != nil {
				releaseNote = *release.Body
			}

			// TODO: if target commitish is branch, then get lastest commit from
			// branch
			releases = append(releases, &Release{
				CommitSha:   commit,
				ID:          id,
				Name:        name,
				Platform:    "github",
				ReleaseNote: releaseNote,
				Tag:         tag,
				Draft:       draft,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return releases, err
}

// ListDraftReleases from a Github repository
// Important: information about published releases are available to everyone.
// Only users with push access will receive listings for draft releases.
func (p *githubPlatform) ListDraftReleases(owner, repository string) (releases []*Release, err error) {
	releaseList, err := p.ListReleases(owner, repository)
	if err != nil {
		return
	}
	for _, release := range releaseList {
		if release.Draft {
			releases = append(releases, release)
		}
	}
	return
}

// UpdateRelease edit a release based on a provided releases ID and release note
func (p *githubPlatform) UpdateRelease(owner, repository string, release *Release) (err error) {
	r, _, err := p.client.Repositories.GetRelease(p.context, owner,
		repository, release.ID.(int64))
	if err != nil {
		return
	}

	r.Body = github.String(release.ReleaseNote)

	_, _, err = p.client.Repositories.EditRelease(p.context, owner, repository,
		release.ID.(int64), r)
	if err != nil {
		return
	}

	return
}

// PublishRelease publish a release
func (p *githubPlatform) PublishRelease(owner, repository string, release *Release) (published bool, err error) {
	r, _, err := p.client.Repositories.GetRelease(p.context, owner, repository,
		release.ID.(int64))
	if err != nil {
		return
	}

	r.Draft = github.Bool(false)

	_, _, err = p.client.Repositories.EditRelease(p.context, owner, repository,
		release.ID.(int64), r)
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
		return true, nil
	}

	opts := &github.ListCheckRunsOptions{
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: githubPerPage,
		},
	}

	succeededCheck := 0
	for {
		getCheckRun, resp, err := p.client.Checks.ListCheckRunsForRef(p.context,
			owner, repository, commitSha, opts)
		if err != nil {
			return false, err
		}
		// TODO: make sure all values in statuses are unique
		for _, check := range getCheckRun.CheckRuns {
			for _, status := range statuses {
				if *check.Name == status &&
					check.Status != nil &&
					*check.Status == completedStatusValue &&
					check.Conclusion != nil &&
					*check.Conclusion == successStatusValue {
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

// CreateFile create a file with content at a given path
// This function is only called by integration tests
func (p *githubPlatform) CreateFile(owner, repository, path, branch, commitMessage, body string) (err error) {
	opts := &github.RepositoryContentFileOptions{
		Branch:  github.String(branch),
		Content: []byte(body),
		Message: github.String(commitMessage),
	}

	_, _, err = p.client.Repositories.CreateFile(p.context, owner, repository, path, opts)

	return
}

// UpdateFile update a file with content at a given path
// This function is only called by integration tests
func (p *githubPlatform) UpdateFile(owner, repository, path, branch, commitMessage, body string) (err error) {
	opts := &github.RepositoryContentFileOptions{
		Branch:  github.String(branch),
		Content: []byte(body),
		Message: github.String(commitMessage),
	}

	_, _, err = p.client.Repositories.UpdateFile(p.context, owner, repository, path, opts)

	return
}

// CreateIssue create an issue
func (p *githubPlatform) CreateIssue(owner, repository string, issue *Issue) (err error) {
	issueRequest := &github.IssueRequest{
		Title: github.String(issue.Title),
		Body:  github.String(issue.Body),
	}
	_, _, err = p.client.Issues.Create(p.context, owner, repository, issueRequest)

	return
}

// CreateRelease create a release.
// This function is only called by integration tests
func (p *githubPlatform) CreateRelease(owner, repository string, release *Release) (*Release, error) {
	opts := &github.RepositoryRelease{
		Name:            github.String(release.Name),
		TargetCommitish: github.String(release.CommitSha),
		TagName:         github.String(release.Tag),
		Draft:           github.Bool(release.Draft),
		Body:            github.String(release.ReleaseNote),
	}
	r, _, err := p.client.Repositories.CreateRelease(p.context, owner, repository, opts)
	if err != nil {
		return nil, err
	}

	release.ID = *r.ID
	release.CommitSha = *r.TargetCommitish

	return release, err
}

// CreateRepository create a repository
// This function is only called by integration tests
func (p *githubPlatform) CreateRepository(owner, repository, visibility string) (err error) {
	opts := &github.Repository{
		Name:       github.String(repository),
		Visibility: github.String(visibility),
	}
	_, _, err = p.client.Repositories.Create(p.context, owner, opts)
	return
}

// CreateStatus returns the status of a specific commit matching a provided status name
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

// DeleteRepository delete a repository
// This function is only called by integration tests
func (p *githubPlatform) DeleteRepository(owner, repository string) (err error) {
	_, err = p.client.Repositories.Delete(p.context, owner, repository)
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

// ListIssuesByAuthor from a given repository
func (p *githubPlatform) ListIssuesByAuthor(owner, repository string,
	author interface{}) (issueList []*Issue, err error) {
	opts := &github.IssueListByRepoOptions{
		Creator: author.(string),
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: githubPerPage,
		},
	}

	for {
		issuesFromRepo, resp, err := p.client.Issues.ListByRepo(p.context, owner,
			repository, opts)
		if err != nil {
			return nil, err
		}

		for _, issue := range issuesFromRepo {
			issueList = append(issueList, &Issue{
				Body:  *issue.Body,
				ID:    *issue.Number,
				Title: *issue.Title,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opts.ListOptions.Page = resp.NextPage
	}

	return issueList, err
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

// UpdateIssue update an issue
func (p *githubPlatform) UpdateIssue(owner, repository string, issue *Issue) (err error) {
	issueRequest := &github.IssueRequest{
		Title: github.String(issue.Title),
		Body:  github.String(issue.Body),
	}
	_, _, err = p.client.Issues.Edit(p.context, owner, repository, issue.ID.(int), issueRequest)

	return
}
