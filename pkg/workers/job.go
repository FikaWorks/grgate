package workers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/fikaworks/grgate/pkg/utils"
)

// Job define information about the job to process
type Job struct {
	Platform   platforms.Platform
	Owner      string
	Repository string
	Config     *config.RepoConfig
}

// NewJob return a Job to be processed by a worker
func NewJob(platform platforms.Platform, owner, repository string) (job *Job, err error) {
	repoConfig, err := config.NewRepoConfig(platform, owner, repository)
	if err != nil {
		return
	}

	job = &Job{
		Platform:   platform,
		Owner:      owner,
		Repository: repository,
		Config:     repoConfig,
	}
	return
}

// processReleaseNote update releases description with statuses based on the
// release template defined in config
func (j *Job) processReleaseNote(release *platforms.Release) (err error) {
	if !j.Config.ReleaseNote.Enabled {
		return
	}
	log.Info().
		Str("repository", j.Repository).
		Str("owner", j.Owner).
		Str("releaseCommit", release.CommitSha).
		Str("releaseTag", release.Tag).
		Str("releaseName", release.Name).
		Msg("Updating status list in release note")

	statusList, err := j.Platform.ListStatuses(j.Owner, j.Repository,
		release.CommitSha)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", j.Owner).
			Str("repository", j.Repository).
			Str("releaseCommit", release.CommitSha).
			Str("releaseTag", release.Tag).
			Str("releaseName", release.Name).
			Msg("Couldn't list release statuses")
		return
	}

	releaseNoteData := &utils.ReleaseNoteData{
		ReleaseNote: release.ReleaseNote,
		Statuses:    utils.MergeStatuses(statusList, j.Config.Statuses),
	}
	release.ReleaseNote, err = utils.RenderReleaseNote(j.Config.ReleaseNote.Template,
		releaseNoteData)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", j.Owner).
			Str("repository", j.Repository).
			Str("releaseCommit", release.CommitSha).
			Str("releaseTag", release.Tag).
			Str("releaseName", release.Name).
			Msg("Couldn't render release note")
		return
	}

	if j.Config.Enabled {
		err = j.Platform.UpdateRelease(j.Owner, j.Repository, release)
		if err != nil {
			log.Error().
				Err(err).
				Str("owner", j.Owner).
				Str("repository", j.Repository).
				Str("releaseCommit", release.CommitSha).
				Str("releaseTag", release.Tag).
				Str("releaseName", release.Name).
				Msg("Couldn't update release")
			return
		}
	} else {
		log.Info().
			Str("repository", j.Repository).
			Str("owner", j.Owner).
			Str("releaseCommit", release.CommitSha).
			Str("releaseTag", release.Tag).
			Str("releaseName", release.Name).
			Msgf("Would update release note with statuses [dry-run]")
	}

	return nil
}

// findIssueDashboard look for issue matching dashboard title
func (j *Job) findIssueDashboard(issueList []*platforms.Issue) *platforms.Issue {
	for _, issue := range issueList {
		if issue.Title == j.Config.Dashboard.Title {
			return issue
		}
	}
	return nil
}

// processDashboard look for each issues created by the author in a repository,
// then update issue with current GRGate state of the first issue matching the
// dashboard title
func (j *Job) processDashboard(errorList []string) {
	if !j.Config.Dashboard.Enabled {
		return
	}
	log.Info().
		Str("repository", j.Repository).
		Str("owner", j.Owner).
		Msg("Updating dashboard")

	issueList, err := j.Platform.ListIssuesByAuthor(j.Owner, j.Repository, j.Config.Dashboard.Author)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", j.Owner).
			Str("repository", j.Repository).
			Msg("Couldn't list issues")
		return
	}
	log.Debug().
		Str("owner", j.Owner).
		Str("repository", j.Repository).
		Msgf("Found %d dashboard issue(s)", len(issueList))

	body, err := utils.RenderDashboard(
		j.Config.Dashboard.Template,
		&utils.DashboardData{
			Enabled: j.Config.Enabled,
			Errors:  errorList,
		},
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", j.Owner).
			Str("repository", j.Repository).
			Msg("Couldn't render dashboard")
		return
	}

	issue := j.findIssueDashboard(issueList)
	if issue != nil {
		issue.Body = body
		err = j.Platform.UpdateIssue(j.Owner, j.Repository, issue)
		if err != nil {
			log.Error().
				Err(err).
				Str("owner", j.Owner).
				Str("repository", j.Repository).
				Msg("Couldn't update issue")
			return
		}
	} else {
		issue = &platforms.Issue{
			Title: j.Config.Dashboard.Title,
			Body:  body,
		}
		err = j.Platform.CreateIssue(j.Owner, j.Repository, issue)
		if err != nil {
			log.Error().
				Err(err).
				Str("owner", j.Owner).
				Str("repository", j.Repository).
				Msg("Couldn't create issue")
			return
		}
	}
}

// Process job by getting all the draft/unpublished releases, for each release
// check that all the required status succeeded then publish the release
func (j *Job) Process() (err error) {
	var errorDashboardList []string

	defer func() {
		j.processDashboard(errorDashboardList)
	}()

	log.Info().
		Str("repository", j.Repository).
		Str("owner", j.Owner).
		Msgf("Processing")
	log.Info().
		Str("repository", j.Repository).
		Str("owner", j.Owner).
		Msgf("Dry run: %t", !j.Config.Enabled)
	log.Info().
		Str("repository", j.Repository).
		Str("owner", j.Owner).
		Msgf("Matching statuses: %s", strings.Join(j.Config.Statuses, ", "))
	log.Info().
		Str("repository", j.Repository).
		Str("owner", j.Owner).
		Msgf("Matching tag regexp: %s", j.Config.TagRegexp)

	if len(j.Config.Statuses) == 0 {
		log.Info().
			Str("repository", j.Repository).
			Str("owner", j.Owner).
			Msg("Statuses are undefined in config, skipping process")
		errorDashboardList = append(errorDashboardList, "Statuses are undefined in .grgate.yaml")
		return nil
	}

	tagRegexp, err := regexp.Compile(j.Config.TagRegexp)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", j.Owner).
			Str("repository", j.Repository).
			Msgf("Couldn't compile regexp \"%s\"", j.Config.TagRegexp)
		errorDashboardList = append(errorDashboardList,
			fmt.Sprintf("Couldn't compile regexp \"%s\"", j.Config.TagRegexp))
		return err
	}

	releaseList, err := j.Platform.ListDraftReleases(j.Owner, j.Repository)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", j.Owner).
			Str("repository", j.Repository).
			Msg("Couldn't list draft releases")
		return err
	}

	log.Info().
		Str("repository", j.Repository).
		Str("owner", j.Owner).
		Msgf("Found %d release(s) marked as draft", len(releaseList))

	for _, release := range releaseList {
		if !tagRegexp.MatchString(release.Tag) {
			log.Debug().
				Str("repository", j.Repository).
				Str("owner", j.Owner).
				Str("releaseCommit", release.CommitSha).
				Str("releaseTag", release.Tag).
				Str("releaseName", release.Name).
				Msgf("Release do not match provided target tag %s", j.Config.TagRegexp)
			continue
		}

		log.Debug().
			Str("repository", j.Repository).
			Str("owner", j.Owner).
			Str("releaseCommit", release.CommitSha).
			Str("releaseTag", release.Tag).
			Str("releaseName", release.Name).
			Msgf("Release match provided target tag %s", j.Config.TagRegexp)

		succeeded, err := j.Platform.CheckAllStatusSucceeded(j.Owner,
			j.Repository, release.CommitSha, j.Config.Statuses)
		if err != nil {
			log.Error().
				Err(err).
				Str("owner", j.Owner).
				Str("repository", j.Repository).
				Str("releaseCommit", release.CommitSha).
				Str("releaseTag", release.Tag).
				Str("releaseName", release.Name).
				Msg("Couldn't check all status check")
			return err
		}

		if err = j.processReleaseNote(release); err != nil {
			return err
		}

		log.Trace().
			Str("repository", j.Repository).
			Str("owner", j.Owner).
			Str("releaseCommit", release.CommitSha).
			Str("releaseTag", release.Tag).
			Str("releaseName", release.Name).
			Msgf("CheckAllStatusSucceeded: %t", succeeded)

		if succeeded {
			if !j.Config.Enabled {
				log.Info().
					Str("repository", j.Repository).
					Str("owner", j.Owner).
					Str("releaseCommit", release.CommitSha).
					Str("releaseTag", release.Tag).
					Str("releaseName", release.Name).
					Msgf("All required status succeeded, would publish release [dry-run]")
				continue
			}

			log.Debug().
				Str("repository", j.Repository).
				Str("owner", j.Owner).
				Str("releaseCommit", release.CommitSha).
				Str("releaseTag", release.Tag).
				Str("releaseName", release.Name).
				Msg("All required status succeeded, publishing release...")

			_, err := j.Platform.PublishRelease(j.Owner, j.Repository, release)
			if err != nil {
				log.Error().
					Err(err).
					Str("owner", j.Owner).
					Str("repository", j.Repository).
					Str("releaseCommit", release.CommitSha).
					Str("releaseTag", release.Tag).
					Str("releaseName", release.Name).
					Msg("Couldn't publish release")
				return err
			}

			log.Info().
				Str("repository", j.Repository).
				Str("owner", j.Owner).
				Str("releaseCommit", release.CommitSha).
				Str("releaseTag", release.Tag).
				Str("releaseName", release.Name).
				Msg("Successfully published release")
		}
	}

	return nil
}
