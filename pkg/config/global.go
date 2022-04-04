package config

import (
	"github.com/spf13/viper"
)

// NewGlobalConfig define the viper configuration and set the global
// config.Main variable base on a config file
func NewGlobalConfig(path string) (v *viper.Viper, err error) {
	v = viper.New()

	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.AddConfigPath("/etc/grgate")
		v.SetConfigName("config.yaml")
		v.SetConfigType("yaml")
	}

	v.SetEnvPrefix("grgate")

	// Set defaults
	v.SetDefault("globals.enabled", DefaultEnabled)
	v.SetDefault("globals.dashboard.enabled", DefaultDashboardEnabled)
	v.SetDefault("globals.dashboard.author", DefaultDashboardAuthor)
	v.SetDefault("globals.dashboard.title", DefaultDashboardTitle)
	v.SetDefault("globals.dashboard.template", DefaultDashboardTemplate)
	v.SetDefault("globals.releaseNote.enabled", DefaultReleaseNoteEnabled)
	v.SetDefault("globals.releaseNote.template", DefaultReleaseNoteTemplate)
	v.SetDefault("globals.tagRegexp", DefaultTagRegexp)
	v.SetDefault("platform", DefaultPlatform)
	v.SetDefault("repoConfigPath", DefaultRepoConfigPath)
	v.SetDefault("server.listenAddress", DefaultServerListenAddress)
	v.SetDefault("server.metricsAddress", DefaultServerMetricsAddress)
	v.SetDefault("server.probeAddress", DefaultServerProbeAddress)
	v.SetDefault("workers", DefaultWorkers)

	err = v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// config file not found, use fallback to default config
			return
		}
	}

	return v, v.Unmarshal(&Main)
}
