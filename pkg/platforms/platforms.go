package platforms

import (
	"io"
)

// Platform interface Github and Gitlab
type Platform interface {
	CheckAllStatusSucceeded(string, string, string, []string) (bool, error)
	CreateStatus(string, string, *Status) error
	GetStatus(string, string, string, string) (*Status, error)
	ListReleases(string, string) ([]*Release, error)
	ListStatuses(string, string, string) ([]*Status, error)
	PublishRelease(string, string, interface{}) (bool, error)
	ReadFile(string, string, string) (io.Reader, error)
}

// Release represent a release regarding the platform
type Release struct {
	// CommitSha attached to the release
	CommitSha string

	// ID of the release, Github use an int, Gitlab use a string
	ID interface{}

	// Name of the release
	Name string

	// Platform, either github or gitlab
	Platform string

	// Tag associated to the release
	Tag string
}

// Status contains commit status informations
type Status struct {
	// CommitSha
	CommitSha string

	// Name of the status
	Name string

	// State is only used by Github checks, must be one of success or in_progress
	State string

	// Status the commit status:
	// For Github must be one of: queued, in_progress or completed
	// For Gitlab must be one of: pending, running, success, failed or cancelled
	Status string
}
