//go:build unit

package utils

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestIsValidRepositoryName(t *testing.T) {
	testCases := map[string]bool{
		"organization/repository":         true,
		"OrgaNiza.tion123/Repo.Sitory123": true,
		"organization/ repository":        false,
		"singleword-no-slash-separated":   false,
	}
	for value, expected := range testCases {
		if result := IsValidRepositoryName(value); result != expected {
			t.Errorf("Repository name %s, got %t, expected %t", value, result, expected)
		}
	}
}

func TestGetRepositoryOrganization(t *testing.T) {
	testCases := map[string]string{
		"organization/repository":       "organization",
		"singleword-no-slash-separated": "singleword-no-slash-separated",
	}
	for value, expected := range testCases {
		if result := GetRepositoryOrganization(value); result != expected {
			t.Errorf("Repository name %s, got %s, expected %s", value, result, expected)
		}
	}
}

func TestGetRepositoryName(t *testing.T) {
	testCases := map[string]string{
		"organization/repository":       "repository",
		"singleword-no-slash-separated": "",
	}
	for value, expected := range testCases {
		if result := GetRepositoryName(value); result != expected {
			t.Errorf("Repository name %s, got %s, expected %s", value, result, expected)
		}
	}
}

func TestExtractRepository(t *testing.T) {
	testCases := []struct {
		input  string
		error  error
		output Repository
	}{
		{
			input: "https://github.com/my-org/my-repo",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "https",
				Platform: "github",
			},
		}, {
			input: "github.com/my-org/my-repo/",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "https",
				Platform: "github",
			},
		}, {
			input: "github.com/my-org/my-repo",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "https",
				Platform: "github",
			},
		}, {
			input: "gitlab.com/my-org/my-repo",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "https",
				Platform: "gitlab",
			},
		}, {
			input: "https://github.com/my-org/my-repo.git",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "https",
				Platform: "github",
			},
		}, {
			input: "git@github.com:my-org/my-repo",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "ssh",
				Platform: "github",
			},
		}, {
			input: "git@github.com:my-org/my-repo.git",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "ssh",
				Platform: "github",
			},
		}, {
			input: "ssh://git@github.com/my-org/my-repo",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "ssh",
				Platform: "github",
			},
		}, {
			input: "ssh://git@github.com/my-org/my-repo.git",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "ssh",
				Platform: "github",
			},
		}, {
			input: "https://github.com/my-org/my-repo_123",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo_123",
				Scheme:   "https",
				Platform: "github",
			},
		}, {
			input: "https://gitlab.com/my-org/my-repo",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "https",
				Platform: "gitlab",
			},
		}, {
			input: "my-org/my-repo",
			output: Repository{
				Owner:    "my-org",
				Name:     "my-repo",
				Scheme:   "https",
				Platform: "",
			},
		}, {
			input: "https://example.com",
			error: fmt.Errorf("cannot parse provided repository url or owner/name"),
		},
	}

	for _, test := range testCases {
		result, err := ExtractRepository(test.input)
		if test.error != nil {
			if diff := pretty.Compare(err, test.error); diff != "" {
				t.Errorf("input %s\ngot error: %s\nwant: %s", test.input, err, test.error)
			}
			continue
		}

		if diff := pretty.Compare(result, test.output); diff != "" {
			t.Errorf("input %s\ndiff: (-got +want)\n%s", test.input, diff)
		}
	}
}
