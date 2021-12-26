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
				Enabled: true,
				ReleaseNote: &ReleaseNote{
					Enabled:  true,
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
releaseNote:
  enabled: false
  template: |-
    no template
statuses:
  - happy-flow`), nil
					})

			expectedRepoConfig := RepoConfig{
				Enabled: true,
				ReleaseNote: &ReleaseNote{
					Enabled:  false,
					Template: "no template",
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
