package platforms

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/xanzy/go-gitlab"
)

const (
	// number of items per page to retrieve via the Gitlab API
	gitlabePerPage int = 100
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
		Page:    0,
		PerPage: gitlabePerPage,
	}

	for {
		releaseList, resp, err := p.client.Releases.ListReleases(getPID(owner,
			repository), opts, nil)
		if err != nil {
			return nil, err
		}

		for _, release := range releaseList {
			tag := release.TagName
			name := release.Name
			commit := release.Commit.ID

			// If the release is in the future, then this is a "draft release"
			if release.ReleasedAt.After(time.Now().UTC()) {
				releases = append(releases, &Release{
					CommitSha: commit,
					Name:      name,
					Platform:  "gitlab",
					Tag:       tag,
					ID:        tag,
				})
			}
		}

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opts.Page = resp.NextPage
	}

	return releases, err
}

// PublishRelease publish a release based on a provided release ID (tag name)
func (p *gitlabPlatform) PublishRelease(owner, repository string, id interface{}) (published bool, err error) {
	releasedAt := time.Now().UTC()

	opts := &gitlab.UpdateReleaseOptions{
		ReleasedAt: &releasedAt,
	}

	_, _, err = p.client.Releases.UpdateRelease(getPID(owner, repository),
		id.(string), opts, nil)
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
		return
	}

	opts := &gitlab.GetCommitStatusesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    0,
			PerPage: gitlabePerPage,
		},
	}

	for {
		commitStatuses, resp, err := p.client.Commits.GetCommitStatuses(getPID(
			owner, repository), commitSha, opts, nil)
		if err != nil {
			return false, err
		}

		// for all statuses, check if the provided one are all successful
		succeededStatus := 0
		for _, commitStatus := range commitStatuses {
			for _, status := range statuses {
				if commitStatus.Name == status && commitStatus.Status == "success" {
					succeededStatus++
				}
			}
		}

		succeeded = succeededStatus == len(statuses)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return succeeded, err
}

// CreateStatus for a given commit
func (p *gitlabPlatform) CreateStatus(owner, repository string, status *Status) (err error) {
	opts := &gitlab.SetCommitStatusOptions{
		State: *gitlab.BuildState(gitlab.BuildStateValue(status.Status)),
		Name:  &status.Name,
	}

	_, _, err = p.client.Commits.SetCommitStatus(getPID(owner, repository),
		status.CommitSha, opts, nil)

	return
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

func getPID(owner, repository string) string {
	return fmt.Sprintf("%s/%s", owner, repository)
}
