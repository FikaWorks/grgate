package config

import (
	"bytes"

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
		cfg = bytes.NewBuffer([]byte{})
	} else {
		log.Info().
			Str("owner", owner).
			Str("repository", repository).
			Msgf("Found file \"%s\" in repository, overriding settings",
				Main.RepoConfigPath)
	}

	v := viper.New()
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("Enabled", Main.Globals.Enabled)
	v.SetDefault("Statuses", Main.Globals.Statuses)
	v.SetDefault("TagRegexp", Main.Globals.TagRegexp)

	if err = v.ReadConfig(cfg); err != nil {
		return
	}

	if err = v.Unmarshal(&config); err != nil {
		log.Error().
			Err(err).
			Str("owner", owner).
			Str("repository", repository).
			Msg("couldn't unmarshal repo config")
		return
	}

	return config, nil
}
