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

	// Set defaults
	v.SetDefault("platform", DefaultPlatform)
	v.SetDefault("globals.enabled", true)
	v.SetDefault("globals.tagRegexp", ".*")
	v.SetDefault("globals.releaseNote.enabled", true)
	v.SetDefault("globals.releaseNote.template", DefaultReleaseNoteTemplate)
	v.SetDefault("repoConfigPath", DefaultRepoConfigPath)
	v.SetDefault("server.listenAddress", DefaultServerListenAddress)
	v.SetDefault("server.metricsAddress", DefaultServerMetricsAddress)
	v.SetDefault("server.probeAddress", DefaultServerProbeAddress)
	v.SetDefault("workers", DefaultWorkers)

	err = v.ReadInConfig()
	if err != nil {
		return
	}

	err = v.Unmarshal(&Main)
	return
}
