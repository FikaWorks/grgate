package utils

import (
	"regexp"
	"strings"
)

var repositoryRegexp = regexp.MustCompile(`^[a-zA-Z-_0-9.]*/[a-zA-Z-_0-9.]*$`)

// IsValidRepositoryName returns true if the repository match the repository
// regexp
func IsValidRepositoryName(repository string) bool {
	return repositoryRegexp.MatchString(repository)
}

// GetRepositoryOrganization returns the organization from a given repository
func GetRepositoryOrganization(repository string) string {
	return strings.Split(repository, "/")[0]
}

// GetRepositoryName returns the name of a given repository
func GetRepositoryName(repository string) string {
  s := strings.Split(repository, "/")
  if len(s) > 1 {
    return s[1]
  }
	return ""
}
