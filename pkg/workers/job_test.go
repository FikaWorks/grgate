//go:build unit

package workers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/platforms"
	mock_platforms "github.com/fikaworks/grgate/pkg/platforms/mocks"
)

func TestProcess(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Run("should publish release with all status succeeded that match tag regexp",
		func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPlatforms := mock_platforms.NewMockPlatform(ctrl)

			mockPlatforms.EXPECT().ListDraftReleases(gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ string, _ string) ([]*platforms.Release, error) {
						return []*platforms.Release{
							{
								ID:          1,
								Tag:         "v1.2.3",
								ReleaseNote: "",
							},
							{
								ID:          2,
								Tag:         "v1.2.3-beta.0",
								ReleaseNote: "",
							},
						}, nil
					})

			mockPlatforms.EXPECT().CheckAllStatusSucceeded(gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string, _ []string) (bool, error) {
					return true, nil
				})

			mockPlatforms.EXPECT().PublishRelease(gomock.Any(), gomock.Any(),
				gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ interface{}) (bool, error) {
					return true, nil
				})

			job := &Job{
				Platform: mockPlatforms,
				Config: &config.RepoConfig{
					Enabled:   true,
					TagRegexp: "^v\\d+\\.\\d+\\.\\d+$",
					ReleaseNote: &config.ReleaseNote{
						Enabled: false,
					},
				},
			}

			if err := job.Process(); err != nil {
				t.Error("error not expected")
			}
		})

	t.Run("should not publish release when RepoConfig.Enabled is false",
		func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPlatforms := mock_platforms.NewMockPlatform(ctrl)

			mockPlatforms.EXPECT().ListDraftReleases(gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ string, _ string) ([]*platforms.Release, error) {
						return []*platforms.Release{
							{
								ID:          1,
								Tag:         "v1.2.3",
								ReleaseNote: "",
							},
						}, nil
					})

			mockPlatforms.EXPECT().CheckAllStatusSucceeded(gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string, _ []string) (bool, error) {
					return true, nil
				})

			job := &Job{
				Platform: mockPlatforms,
				Config: &config.RepoConfig{
					Enabled:   false,
					TagRegexp: ".*",
					ReleaseNote: &config.ReleaseNote{
						Enabled: false,
					},
				},
			}

			if err := job.Process(); err != nil {
				t.Error("error not expected")
			}
		})

	t.Run("should not publish release if not all the status check succeeded",
		func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPlatforms := mock_platforms.NewMockPlatform(ctrl)

			mockPlatforms.EXPECT().ListDraftReleases(gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ string, _ string) ([]*platforms.Release, error) {
						return []*platforms.Release{
							{
								ID:          1,
								Tag:         "v1.2.3",
								ReleaseNote: "",
							},
						}, nil
					})

			mockPlatforms.EXPECT().CheckAllStatusSucceeded(gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string, _ []string) (bool, error) {
					return false, nil
				})

			job := &Job{
				Platform: mockPlatforms,
				Config: &config.RepoConfig{
					Enabled:   true,
					TagRegexp: ".*",
					ReleaseNote: &config.ReleaseNote{
						Enabled: false,
					},
				},
			}

			if err := job.Process(); err != nil {
				t.Error("error not expected")
			}
		})

	t.Run("should update release note with correct status check when RepoConfig.ReleaseNote.Enabled is true",
		func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPlatforms := mock_platforms.NewMockPlatform(ctrl)

			mockPlatforms.EXPECT().ListDraftReleases(gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ string, _ string) ([]*platforms.Release, error) {
						return []*platforms.Release{
							{
								ID:          1,
								Tag:         "v1.2.3",
								ReleaseNote: "This is a release note",
							},
						}, nil
					})

			mockPlatforms.EXPECT().CheckAllStatusSucceeded(gomock.Any(), gomock.Any(),
				gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string, _ []string) (bool, error) {
					return false, nil
				})

			mockPlatforms.EXPECT().ListStatuses(gomock.Any(), gomock.Any(),
				gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string) ([]*platforms.Status, error) {
					return []*platforms.Status{
						{
							Name:   "e2e A",
							Status: "success",
						},
						{
							Name:   "e2e B",
							Status: "pending",
						},
						{
							Name:   "e2e C",
							Status: "running",
						},
						{
							Name:   "e2e D",
							Status: "failed",
						},
						{
							Name:   "e2e E",
							Status: "cancelled",
						},
						{
							Name:   "e2e F",
							Status: "in_progress",
						},
					}, nil
				})

			mockPlatforms.EXPECT().UpdateRelease(gomock.Any(), gomock.Any(),
				gomock.Any()).DoAndReturn(
				func(_ string, _ string, release *platforms.Release) error {
					expectedReleaseNote := `This is a release note
<!-- GRGate start -->
<details><summary>Status check</summary>

- [x] e2e A
- [ ] e2e B
- [ ] e2e C
- [ ] e2e D
- [ ] e2e E
- [ ] e2e F

</details>
<!-- GRGate end -->`

					if diff := pretty.Compare(release.ReleaseNote, expectedReleaseNote); diff != "" {
						t.Errorf("diff: (-got +want)\n%s", diff)
					}
					return nil
				})

			job := &Job{
				Platform: mockPlatforms,
				Config: &config.RepoConfig{
					Enabled:   true,
					TagRegexp: ".*",
					ReleaseNote: &config.ReleaseNote{
						Enabled: true,
						Template: `{{- .ReleaseNote }}
<!-- GRGate start -->
<details><summary>Status check</summary>
{{ range .Statuses }}
- [{{ if eq .Status "success" }}x{{ else }} {{ end }}] {{ .Name }}
{{- end }}

</details>
<!-- GRGate end -->`,
					},
				},
			}

			if err := job.Process(); err != nil {
				t.Error("error not expected")
			}
		})
}
