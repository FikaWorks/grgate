//go:build unit

package config

import (
	"errors"
	"io"
	"strings"
	"testing"

	mock_platforms "github.com/fikaworks/grgate/pkg/platforms/mocks"

	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
)

func TestNewRepo(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Run("should return default global value if repo config file is not present",
		func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			mockPlatforms := mock_platforms.NewMockPlatform(ctrl)

			mockPlatforms.EXPECT().ReadFile(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ string, _ string, _ string) (io.Reader, error) {
						return nil, errors.New("file not found")
					})

			expectedRepoConfig := RepoConfig{
				Enabled: DefaultEnabled,
				Dashboard: &Dashboard{
					Enabled:  DefaultDashboardEnabled,
					Author:   DefaultDashboardAuthor,
					Title:    DefaultDashboardTitle,
					Template: DefaultDashboardTemplate,
				},
				ReleaseNote: &ReleaseNote{
					Enabled:  DefaultReleaseNoteEnabled,
					Template: DefaultReleaseNoteTemplate,
				},
				Statuses:  []string{},
				TagRegexp: ".*",
			}

			_, _ = NewGlobalConfig("")

			repoConfig, err := NewRepoConfig(mockPlatforms, "owner", "repository")
			if err != nil {
				t.Errorf("Error not expected: %#v", err)
			}

			if diff := pretty.Compare(repoConfig, expectedRepoConfig); diff != "" {
				t.Errorf("diff: (-got +want)\n%s", diff)
			}
		})

	t.Run("should return correct data if repo config file is present",
		func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			mockPlatforms := mock_platforms.NewMockPlatform(ctrl)

			mockPlatforms.EXPECT().ReadFile(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ string, _ string, _ string) (io.Reader, error) {
						return strings.NewReader(`enabled: true
dashboard:
  enabled: false
  author: some author
  title: some title
  template: |-
    some template
releaseNote:
  enabled: false
  template: |-
    some template
statuses:
  - happy-flow`), nil
					})

			expectedRepoConfig := RepoConfig{
				Enabled: true,
				Dashboard: &Dashboard{
					Enabled:  false,
					Author:   "some author",
					Title:    "some title",
					Template: "some template",
				},
				ReleaseNote: &ReleaseNote{
					Enabled:  false,
					Template: "some template",
				},
				Statuses:  []string{"happy-flow"},
				TagRegexp: ".*",
			}

			_, _ = NewGlobalConfig("")

			repoConfig, err := NewRepoConfig(mockPlatforms, "owner", "repository")
			if err != nil {
				t.Errorf("Error not expected: %#v", err)
			}

			if diff := pretty.Compare(repoConfig, expectedRepoConfig); diff != "" {
				t.Errorf("diff: (-got +want)\n%s", diff)
			}
		})
}
