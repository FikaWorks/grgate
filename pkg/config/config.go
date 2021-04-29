package config

// Main contains the generated main configuration
var Main *MainConfig

// MainConfig define the main configuration
type MainConfig struct {
  Github Github `mapstructure:"github"`
  Globals RepoConfig `mapstructure:"globals"`
  LogFormat string `mapstructure:"logFormat"`
  LogLevel string `mapstructure:"logLevel"`
  RepoConfigPath string `mapstructure:"repoConfigPath"`
  Server Server `mapstructure:"server"`
  Workers int `mapstructure:"workers"`
}

// Server define server configuration
type Server struct {
  ListenAddress string `mapstructure:"listenAddress"`
  MetricsAddress string `mapstructure:"metricsAddress"`
  ProbeAddress string `mapstructure:"probeAddress"`
}

// Github define Github configuration
type Github struct {
  AppID int64 `mapstructure:"appID"`
  InstallationID int64 `mapstructure:"installationID"`
  PrivateKeyPath string `mapstructure:"privateKeyPath"`
  WebhookSecret string `mapstructure:"webhookSecret"`
}

// RepoConfig define repository configuration
type RepoConfig struct {
  Enabled bool `mapstructure:"enabled"`
  TagRegexp string `mapstructure:"tagRegexp"`
  Statuses []string `mapstructure:"statuses"`
}
