package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type Repository struct {
	Name     string
	Owner    string
	Platform string
	Scheme   string
}

// Regexp matching a valid repository name: "owner/name"
var repositoryRegexp = regexp.MustCompile(`^[a-zA-Z-_0-9.]*/[a-zA-Z-_0-9.]*$`)

// Regexp matching a valid Git URI
var uriRegexp = regexp.MustCompile(`^((((?P<scheme>https?|ssh):\/\/)?([^@]+@)?(?P<platform>github|gitlab)(.com)?(?P<separator>[\/:])?))?(?P<owner>[a-zA-Z-_0-9.]*)\/(?P<name>[a-zA-Z-_0-9.]*)/?$`)

// Number of named group captured in the above uriRegexp
const captureGroupNumber = 5

// IsValidRepositoryName returns true if the input match the repository
// regexp
func IsValidRepositoryName(input string) bool {
	return repositoryRegexp.MatchString(input)
}

// GetRepositoryOrganization returns the organization from a given string
func GetRepositoryOrganization(input string) string {
	return strings.Split(input, "/")[0]
}

// GetRepositoryName returns the name of a given string
func GetRepositoryName(input string) string {
	s := strings.Split(input, "/")
	if len(s) > 1 {
		return s[1]
	}
	return ""
}

// ExtractRepository extract values from a Git Url and returns a Repository
// struct
func ExtractRepository(input string) (repository *Repository, err error) {
	result := make(map[string]string, captureGroupNumber)

	find := uriRegexp.FindStringSubmatch(input)
	if find == nil {
		err = fmt.Errorf("Cannot parse provided repository uri or owner/name")
		return
	}

	// store all named capture into result map
	for i, name := range uriRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = find[i]
		}
	}

	// attempt detecting scheme if not defined
	if result["scheme"] == "" {
		result["scheme"] = "https"

		// scheme is ssh if separator between domain name and organization is :
		if result["separator"] == ":" {
			result["scheme"] = "ssh"
		}
	}

	repository = &Repository{
		Name:     strings.TrimRight(result["name"], ".git"),
		Owner:    result["owner"],
		Platform: result["platform"],
		Scheme:   result["scheme"],
	}

	return
}
