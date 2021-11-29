package workers

import (
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

// Process job by getting all the draft/unpublished releases, for each release
// check that all the required status succeeded then publish the release
func (j *Job) Process() error {
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

	tagRegexp, err := regexp.Compile(j.Config.TagRegexp)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", j.Owner).
			Str("repository", j.Repository).
			Msgf("Couldn't compile regexp \"%s\"", j.Config.TagRegexp)
		return err
	}

	releaseList, err := j.Platform.ListReleases(j.Owner, j.Repository)
	if err != nil {
		log.Error().
			Err(err).
			Str("owner", j.Owner).
			Str("repository", j.Repository).
			Msg("Couldn't list releases")
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

		if j.Config.ReleaseNote.Enabled {
			log.Info().
				Str("repository", j.Repository).
				Str("owner", j.Owner).
				Str("releaseCommit", release.CommitSha).
				Str("releaseTag", release.Tag).
				Str("releaseName", release.Name).
				Msg("Updating statuses list in release note")

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
				return err
			}

			releaseNoteData := &utils.ReleaseNoteData{
				ReleaseNote: release.ReleaseNote,
				Statuses:    statusList,
			}
			releaseNote, err := utils.RenderReleaseNote(j.Config.ReleaseNote.Template,
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
				return err
			}
			if j.Config.Enabled {
				err = j.Platform.UpdateRelease(j.Owner, j.Repository, release.ID,
					releaseNote)
				if err != nil {
					log.Error().
						Err(err).
						Str("owner", j.Owner).
						Str("repository", j.Repository).
						Str("releaseCommit", release.CommitSha).
						Str("releaseTag", release.Tag).
						Str("releaseName", release.Name).
						Msg("Couldn't update release")
					return err
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

			_, err := j.Platform.PublishRelease(j.Owner, j.Repository, release.ID)
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
