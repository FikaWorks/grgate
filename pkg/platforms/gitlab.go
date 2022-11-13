package platforms

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/xanzy/go-gitlab"
)

const (
	// default number of items per page to retrieve via the Gitlab API
	gitlabPerPage int = 100

	// default number of hours to set a future release to (aka draft release)
	// Gitlab doesn't distinguish between draft and published release but use
	// release date
	futureReleaseTime time.Duration = time.Hour * 24 * 365
)

// GitlabConfig hold the Gitlab configuration
type GitlabConfig struct {
	Token         string
	WebhookSecret string
}

type gitlabPlatform struct {
	config *GitlabConfig
	client *gitlab.Client
}

// NewGitlab returns an instance of platform
func NewGitlab(config *GitlabConfig) (platform Platform, err error) {
	client, err := gitlab.NewClient(config.Token)
	if err != nil {
		return
	}

	platform = &gitlabPlatform{
		config: config,
		client: client,
	}

	return
}

// ReadFile retrieve file located at the provided path in a given Gitlab repository
func (p *gitlabPlatform) ReadFile(owner, repository, path string) (content io.Reader, err error) {
	r, _, err := p.client.RepositoryFiles.GetRawFile(getPID(owner, repository),
		path, nil)
	if err != nil {
		return
	}
	content = bytes.NewBuffer(r)
	return
}

// ListReleases from a Gitlab repository
func (p *gitlabPlatform) ListReleases(owner, repository string) (releases []*Release, err error) {
	opts := &gitlab.ListReleasesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    0,
			PerPage: gitlabPerPage,
		},
	}

	for {
		releaseList, resp, err := p.client.Releases.ListReleases(getPID(owner,
			repository), opts, nil)
		if err != nil {
			return nil, err
		}

		for _, release := range releaseList {
			// if the release is in the future, then this is a "draft release"
			draft := release.ReleasedAt.After(time.Now().UTC())

			releases = append(releases, &Release{
				CommitSha:   release.Commit.ID,
				ID:          release.TagName,
				Name:        release.Name,
				Platform:    "gitlab",
				ReleaseNote: release.Description,
				Tag:         release.TagName,
				Draft:       draft,
			})
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opts.Page = resp.NextPage
	}

	return releases, err
}

// ListDraftReleases from a Gitlab repository
func (p *gitlabPlatform) ListDraftReleases(owner, repository string) (releases []*Release, err error) {
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
func (p *gitlabPlatform) UpdateRelease(owner, repository string, release *Release) (err error) {
	r, _, err := p.client.Releases.GetRelease(getPID(owner, repository),
		release.ID.(string), nil)
	if err != nil {
		return
	}

	opts := &gitlab.UpdateReleaseOptions{
		Description: &release.ReleaseNote,
		Name:        &release.Name,
		ReleasedAt:  r.ReleasedAt,
	}

	_, _, err = p.client.Releases.UpdateRelease(getPID(owner, repository),
		release.ID.(string), opts, nil)
	if err != nil {
		return
	}

	return
}

// PublishRelease publish a release
func (p *gitlabPlatform) PublishRelease(owner, repository string, release *Release) (published bool, err error) {
	releasedAt := time.Now().UTC()

	opts := &gitlab.UpdateReleaseOptions{
		ReleasedAt:  &releasedAt,
		Description: &release.ReleaseNote,
		Name:        &release.Name,
	}

	_, _, err = p.client.Releases.UpdateRelease(getPID(owner, repository),
		release.ID.(string), opts, nil)
	if err != nil {
		return
	}

	published = true
	return
}

// CheckAllStatusSucceeded checks that all the provided statuses succeeded
func (p *gitlabPlatform) CheckAllStatusSucceeded(owner, repository,
	commitSha string, statuses []string) (succeeded bool, err error) {
	if len(statuses) == 0 {
		return true, nil
	}

	opts := &gitlab.GetCommitStatusesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    0,
			PerPage: gitlabPerPage,
		},
	}

	succeededStatus := 0
	for {
		commitStatuses, resp, err := p.client.Commits.GetCommitStatuses(getPID(
			owner, repository), commitSha, opts, nil)
		if err != nil {
			return false, err
		}

		// for all statuses, check if the provided one are all successful
		for _, commitStatus := range commitStatuses {
			for _, status := range statuses {
				if commitStatus.Name == status && commitStatus.Status == successStatusValue {
					succeededStatus++
				}
			}
		}

		succeeded = succeededStatus == len(statuses)

		if resp.NextPage == 0 {
			break
		}

		opts.ListOptions.Page = resp.NextPage
	}

	return succeeded, err
}

// CreateFile create a file with content at a given path
// This function is only called by integration tests
func (p *gitlabPlatform) CreateFile(owner, repository, path, branch, commitMessage, body string) (err error) {
	opts := &gitlab.CreateFileOptions{
		Branch:        &branch,
		Content:       &body,
		CommitMessage: &commitMessage,
	}
	_, _, err = p.client.RepositoryFiles.CreateFile(getPID(owner, repository), path, opts, nil)
	return
}

// UpdateFile update a file with content at a given path
// This function is only called by integration tests
func (p *gitlabPlatform) UpdateFile(owner, repository, path, branch, commitMessage, body string) (err error) {
	opts := &gitlab.UpdateFileOptions{
		Branch:        &branch,
		Content:       &body,
		CommitMessage: &commitMessage,
	}
	_, _, err = p.client.RepositoryFiles.UpdateFile(getPID(owner, repository), path, opts, nil)
	return
}

// CreateIssue create an issue
func (p *gitlabPlatform) CreateIssue(owner, repository string, issue *Issue) (err error) {
	opts := &gitlab.CreateIssueOptions{
		Title:       &issue.Title,
		Description: &issue.Body,
	}
	_, _, err = p.client.Issues.CreateIssue(getPID(owner, repository), opts, nil)

	return
}

// CreateRelease create a release.
// This function is only called by integration tests
func (p *gitlabPlatform) CreateRelease(owner, repository string, release *Release) (*Release, error) {
	opts := &gitlab.CreateReleaseOptions{
		Name:        &release.Name,
		Ref:         &release.CommitSha,
		TagName:     &release.Tag,
		Description: &release.ReleaseNote,
	}

	if release.Draft {
		// if draft release, set releasedAt to 1 year from now
		future := time.Now().UTC().Add(futureReleaseTime)
		opts.ReleasedAt = &future
	}

	r, _, err := p.client.Releases.CreateRelease(getPID(owner, repository), opts, nil)
	if err != nil {
		return nil, err
	}

	release.ID = r.TagName
	release.CommitSha = r.Commit.ID

	return release, err
}

// CreateRepository create a repository
// This function is only called by integration tests
func (p *gitlabPlatform) CreateRepository(owner, repository, visibility string) (err error) {
	opts := &gitlab.CreateProjectOptions{
		Name:       gitlab.String(repository),
		Visibility: gitlab.Visibility(gitlab.VisibilityValue(visibility)),
	}
	_, _, err = p.client.Projects.CreateProject(opts, nil)
	return
}

// CreateStatus for a given commit
func (p *gitlabPlatform) CreateStatus(owner, repository string, status *Status) (err error) {
	// safely map Github to Gitlab state
	state := mapGithubStatusToGitlabStatus(status.Status)

	opts := &gitlab.SetCommitStatusOptions{
		State: *gitlab.BuildState(gitlab.BuildStateValue(state)),
		Name:  &status.Name,
	}

	_, _, err = p.client.Commits.SetCommitStatus(getPID(owner, repository),
		status.CommitSha, opts, nil)

	return
}

// DeleteRepository delete a repository
// This function is only called by integration tests
func (p *gitlabPlatform) DeleteRepository(owner, repository string) (err error) {
	_, err = p.client.Projects.DeleteProject(getPID(owner, repository), nil)
	return
}

// ListIssuesByAuthor from a given repository
func (p *gitlabPlatform) ListIssuesByAuthor(owner, repository string,
	author interface{}) (issueList []*Issue, err error) {
	opts := &gitlab.ListProjectIssuesOptions{
		AuthorUsername: gitlab.String(author.(string)),
		ListOptions: gitlab.ListOptions{
			Page:    0,
			PerPage: gitlabPerPage,
		},
	}

	for {
		issuesFromRepo, resp, err := p.client.Issues.ListProjectIssues(getPID(owner, repository), opts, nil)
		if err != nil {
			return nil, err
		}

		for _, issue := range issuesFromRepo {
			issueList = append(issueList, &Issue{
				Body:  issue.Description,
				ID:    issue.ID,
				Title: issue.Title,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opts.ListOptions.Page = resp.NextPage
	}

	return issueList, err
}

// GetStatus returns the status of a specific commit matching a provided status name
func (p *gitlabPlatform) GetStatus(owner, repository, commitSha, statusName string) (status *Status, err error) {
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
func (p *gitlabPlatform) ListStatuses(owner, repository, commitSha string) (statusList []*Status, err error) {
	commitStatuses, _, err := p.client.Commits.GetCommitStatuses(getPID(owner,
		repository), commitSha, nil, nil)
	if err != nil {
		return nil, err
	}

	for _, commitStatus := range commitStatuses {
		cr := &Status{
			CommitSha: commitStatus.SHA,
			Name:      commitStatus.Name,
			Status:    commitStatus.Status,
		}

		statusList = append(statusList, cr)
	}

	return statusList, err
}

// UpdateIssue update an issue
func (p *gitlabPlatform) UpdateIssue(owner, repository string, issue *Issue) (err error) {
	opts := &gitlab.UpdateIssueOptions{
		Title:       &issue.Title,
		Description: &issue.Body,
	}
	_, _, err = p.client.Issues.UpdateIssue(getPID(owner, repository), issue.ID.(int), opts)

	return
}

func getPID(owner, repository string) string {
	return fmt.Sprintf("%s/%s", owner, repository)
}
