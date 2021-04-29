package utils

import  (
  "regexp"
  "strings"
)

var repositoryRegexp = regexp.MustCompile(`^[a-zA-Z-_]*/[a-zA-Z-_]*$`)

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
  return strings.Split(repository, "/")[1]
}
