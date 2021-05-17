// +build integration integrationgithub integrationgitlab

package tests

import (
  "time"
)

const (
	repositoryPrefix = "ggate-integration"
)


type repoConfig struct {
  branch string
  commitMessage string
  content string
  path string
}

type repoRelease struct {
  name string
  ref string
  tag string
  releasedAt time.Time
  draft bool
}

type integrationTest interface {
	setup(*repoConfig, *repoRelease) error
	teardown() error
}
