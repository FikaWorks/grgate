package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/fikaworks/ggate/pkg/platforms"
)

// NewRepoConfig returns configuration defined in a repository
func NewRepoConfig(platform platforms.Platform, owner, repository string) (config *RepoConfig, err error) {
	cfg, err := platform.ReadFile(owner, repository, Main.RepoConfigPath)
	if err != nil {
		log.Info().
			Str("owner", owner).
			Str("repository", repository).
			Msgf("File \"%s\" not found in repository, using default settings",
				Main.RepoConfigPath)

		config = &RepoConfig{}
	} else {
		log.Info().
			Str("owner", owner).
			Str("repository", repository).
			Msgf("Found file \"%s\" in repository, overriding settings",
				Main.RepoConfigPath)

		v := viper.New()
		v.SetConfigName(".ggate.yaml")
		v.SetConfigType("yaml")
		if err = v.ReadConfig(cfg); err != nil {
			return
		}

		if err = v.Unmarshal(&config); err != nil {
			log.Error().
				Err(err).
				Str("owner", owner).
				Str("repository", repository).
				Msgf("couldn't unmarshal config \"%s\" from repository",
					Main.RepoConfigPath)
			return
		}
	}

	// set defaults from Globals if not defined in the repository configuration
	if config.TagRegexp == "" {
		config.TagRegexp = Main.Globals.TagRegexp
	}

	if len(config.Statuses) == 0 {
		config.Statuses = Main.Globals.Statuses
	}

	return config, nil
}
