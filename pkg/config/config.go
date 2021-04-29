package config

import (
  "fmt"

	"github.com/spf13/viper"
	"github.com/rs/zerolog/log"

  "github.com/fikaworks/ggate/pkg/platforms"
)

var Main *MainConfig

type MainConfig struct {
  Github Github `mapstructure:"github"`
  Globals RepoConfig `mapstructure:"globals"`
  LogFormat string `mapstructure:"logFormat"`
  LogLevel string `mapstructure:"logLevel"`
  RepoConfigPath string `mapstructure:"repoConfigPath"`
  Server Server `mapstructure:"server"`
  Workers int `mapstructure:"workers"`
}

type Server struct {
  ListenAddress string `mapstructure:"listenAddress"`
  MetricsAddress string `mapstructure:"metricsAddress"`
  ProbeAddress string `mapstructure:"probeAddress"`
}

type Github struct {
  AppID int64 `mapstructure:"appID"`
  InstallationID int64 `mapstructure:"installationID"`
  PrivateKeyPath string `mapstructure:"privateKeyPath"`
  WebhookSecret string `mapstructure:"webhookSecret"`
}

type RepoConfig struct {
  Enabled bool `mapstructure:"enabled"`
  TagRegexp string `mapstructure:"tagRegexp"`
  Statuses []string `mapstructure:"statuses"`
}

// NewRepoConfig returns configuration defined in a repository
func NewRepoConfig(platform platforms.Platform, owner, repository string) (config *RepoConfig, err error) {
  cfg, err := platform.ReadFile(owner, repository, Main.RepoConfigPath)
  if err != nil {
    log.Info().Msgf("File %s not found in repository %s/%s, using default settings",
      Main.RepoConfigPath, owner, repository)

    config = &RepoConfig{}
  } else {
    log.Info().Msgf("Found file %s in repository %s/%s, overriding settings",
      Main.RepoConfigPath, owner, repository)

    v := viper.New()
    v.SetConfigName(".ggate.yaml")
    v.SetConfigType("yaml")
    v.ReadConfig(cfg)

    err = v.Unmarshal(&config)
    if err != nil {
      err = fmt.Errorf("couldn't unmarshal config \"%s\" from repository %s/%s:\n%s",
        Main.RepoConfigPath, owner, repository, err.Error())
      return
    }
  }

  if config.TagRegexp == "" {
    config.TagRegexp = Main.Globals.TagRegexp
  }

  if len(config.Statuses) == 0 {
    config.Statuses = Main.Globals.Statuses
  }

  return config, nil
}

// NewGlobalConfig define the viper configuration and set the global
// config.Main variable base on a config file
func NewGlobalConfig(path string) (v *viper.Viper, err error) {
	v = viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.AddConfigPath("/etc/ggate")
    v.SetConfigName("config.yaml")
    v.SetConfigType("yaml")
  }

  // Set defaults
  v.SetDefault("server.listenAddress", "0.0.0.0:8080")
  v.SetDefault("server.metricsAddress", "0.0.0.0:9101")
  v.SetDefault("server.proveAddress", "0.0.0.0:8086")
  v.SetDefault("workers", 5)
  v.SetDefault("repoConfigPath", ".ggate.yaml")
  v.SetDefault("globals.tagRegexp", ".*")
  v.SetDefault("globals.enabled", true)

	err = v.ReadInConfig()
  if err != nil {
    return
  }

  err = v.Unmarshal(&Main)
	return
}
