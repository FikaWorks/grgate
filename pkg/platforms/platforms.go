package platforms

import (
  "io"
)

// Platform
type Platform interface {
  CreateStatus(string, string, *Status) (error)
  GetStatus(string, string, string, string) (*Status, error)
  HasAllStatusSucceeded(string, string, string, []string) (bool, error)
  ListStatus(string, string, string) ([]*Status, error)
  ListReleases(string, string) ([]*Release, error)
  PublishRelease(string, string, int64) (bool, error)
  ReadFile(string, string, string) (io.ReadCloser, error)
}

// Release
type Release struct {
  CommitSha string
  ID int64
  Platform string
  Tag string
}

// Status
type Status struct {
	// CommitSha
  CommitSha string

	// Name of the status
  Name string

	// Status is one of queued, in_progress or completed
  Status string

	// State
  State string
}
