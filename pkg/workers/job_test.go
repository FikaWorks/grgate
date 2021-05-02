package workers

import (
  "testing"

  "github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"github.com/fikaworks/ggate/pkg/config"
	"github.com/fikaworks/ggate/pkg/platforms"
	mock_platforms "github.com/fikaworks/ggate/pkg/platforms/mocks"
)

func TestProcess(t *testing.T) {
  zerolog.SetGlobalLevel(zerolog.Disabled)

  t.Run("publish release with all status succeeded that match tag regexp",
    func(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockPlatforms := mock_platforms.NewMockPlatform(ctrl)

    mockPlatforms.EXPECT().ListReleases(gomock.Any(), gomock.Any()).
      DoAndReturn(
        func (_ string, _ string) ([]*platforms.Release, error) {
          return []*platforms.Release{
            &platforms.Release{
              ID: 1,
              Tag: "v1.2.3",
            },
            &platforms.Release{
              ID: 2,
              Tag: "v1.2.3-beta.0",
            },
          }, nil
        })

    mockPlatforms.EXPECT().CheckAllStatusSucceeded(gomock.Any(), gomock.Any(),
      gomock.Any(), gomock.Any()).DoAndReturn(
        func (_ string, _ string, _ string, _ []string) (bool, error) {
          return true, nil
        })

    mockPlatforms.EXPECT().PublishRelease(gomock.Any(), gomock.Any(),
      gomock.Any()).DoAndReturn(
        func (_ string, _ string, _ int64) (bool, error) {
          return true, nil
        })

    job := &Job{
      Platform: mockPlatforms,
      Config: &config.RepoConfig{
        Enabled: true,
        TagRegexp: "^v\\d+\\.\\d+\\.\\d+$",
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

    mockPlatforms.EXPECT().ListReleases(gomock.Any(), gomock.Any()).
      DoAndReturn(
        func (_ string, _ string) ([]*platforms.Release, error) {
          return []*platforms.Release{
            &platforms.Release{
              ID: 1,
              Tag: "v1.2.3",
            },
          }, nil
        })

    mockPlatforms.EXPECT().CheckAllStatusSucceeded(gomock.Any(), gomock.Any(),
      gomock.Any(), gomock.Any()).DoAndReturn(
        func (_ string, _ string, _ string, _ []string) (bool, error) {
          return true, nil
        })

    job := &Job{
      Platform: mockPlatforms,
      Config: &config.RepoConfig{
        Enabled: false,
        TagRegexp: ".*",
      },
    }

    if err := job.Process(); err != nil {
      t.Error("error not expected")
    }
  })
}
