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
