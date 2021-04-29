package workers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/fikaworks/ggate/pkg/platforms"
	"github.com/fikaworks/ggate/pkg/config"
)

// Job
type Job struct {
  Platform platforms.Platform
	Owner string
	Repository string
  Config *config.RepoConfig
}

// NewJob return a Job to be processed by a worker
func NewJob(platform platforms.Platform, owner, repository string) (job Job, err error) {
  repoConfig, err := config.NewRepoConfig(platform, owner, repository)
  if err != nil {
    return
  }

	job = Job{
		Platform: platform,
		Owner: owner,
		Repository: repository,
    Config: repoConfig,
	}
  return
}

// Process job
func (j *Job) Process() error {
	log.Info().
    Str("repository", j.Repository).
    Str("owner", j.Owner).
    Msgf("Processing")
  log.Info().
    Str("repository", j.Repository).
    Str("owner", j.Owner).
    Msgf("Dry run: %t", j.Config.Enabled)
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
    return fmt.Errorf("couldn't compile regexp \"%s\": %s",
      j.Config.TagRegexp, err.Error())
  }

  releaseList, err := j.Platform.ListReleases(j.Owner, j.Repository)
  if err != nil {
    return fmt.Errorf("couldn't list releases from %s/%s: %s", j.Owner,
      j.Repository, err.Error())
  }

  log.Info().
    Str("repository", j.Repository).
    Str("owner", j.Owner).
    Msgf("Found %d release(s) marked as draft", len(releaseList))

  for _, release := range releaseList {
    if ! tagRegexp.MatchString(release.Tag) {
      log.Debug().
        Str("repository", j.Repository).
        Str("owner", j.Owner).
        Msgf("Release %s do not match target tag %s", release.Tag,
          j.Config.TagRegexp)
      continue
    }

    log.Debug().
      Str("repository", j.Repository).
      Str("owner", j.Owner).
      Msgf("Release %s match target tag %s", release.Tag,
        j.Config.TagRegexp)

    succeeded, err := j.Platform.HasAllStatusSucceeded(j.Owner,
      j.Repository, release.CommitSha, j.Config.Statuses)
    if err != nil {
      return fmt.Errorf("couldn't check all status check: %s", err.Error())
    }

    log.Debug().
      Str("repository", j.Repository).
      Str("owner", j.Owner).
      Msgf("Release %s @ %s passed all tests: %t", release.Tag,
        release.CommitSha, succeeded)

    if succeeded {
      if !j.Config.Enabled {
        log.Info().
          Str("repository", j.Repository).
          Str("owner", j.Owner).
          Msgf("Would publish release %s with tag %s@%s",
            release.Name, release.Tag, release.CommitSha)
        continue
      }

      log.Info().
        Str("repository", j.Repository).
        Str("owner", j.Owner).
        Msgf("All status succeeded, publishing release with tag %s and commit %s",
          release.Tag, release.CommitSha)
      _, err := j.Platform.PublishRelease(j.Owner, j.Repository, release.ID)
      if err != nil {
        return fmt.Errorf("couldn't publish release: %s", err.Error())
      }

      log.Info().
        Str("repository", j.Repository).
        Str("owner", j.Owner).
        Msgf("Successfully published release with tag %s and commit %s",
          release.Tag, release.CommitSha)
    }
  }

  return nil
}
