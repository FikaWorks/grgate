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
		v.AddConfigPath("/etc/ggate")
		v.SetConfigName("config.yaml")
		v.SetConfigType("yaml")
	}

	// Set defaults
	v.SetDefault("platform", "github")
	v.SetDefault("globals.enabled", true)
	v.SetDefault("globals.tagRegexp", ".*")
	v.SetDefault("repoConfigPath", ".ggate.yaml")
	v.SetDefault("server.listenAddress", "0.0.0.0:8080")
	v.SetDefault("server.metricsAddress", "0.0.0.0:9101")
	v.SetDefault("server.probeAddress", "0.0.0.0:8086")
	v.SetDefault("workers", 5)

	err = v.ReadInConfig()
	if err != nil {
		return
	}

	err = v.Unmarshal(&Main)
	return
}
