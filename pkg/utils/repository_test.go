package utils

import (
	"testing"
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
