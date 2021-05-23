package config

var (
	// Main contains the generated main configuration
	Main *MainConfig

	// Version of GRGate
	Version string
)

type PlatformType string

const (
	GithubPlatform PlatformType = "github"
	GitlabPlatform PlatformType = "gitlab"
)

// MainConfig define the main configuration
type MainConfig struct {
	Github         *Github       `mapstructure:"github"`
	Gitlab         *Gitlab       `mapstructure:"gitlab"`
	Globals        *RepoConfig   `mapstructure:"globals"`
	LogFormat      string        `mapstructure:"logFormat"`
	LogLevel       string        `mapstructure:"logLevel"`
	Platform       *PlatformType `mapstructure:"platform"`
	RepoConfigPath string        `mapstructure:"repoConfigPath"`
	Server         *Server       `mapstructure:"server"`
	Workers        int           `mapstructure:"workers"`
}

// Server define server configuration
type Server struct {
	ListenAddress  string `mapstructure:"listenAddress"`
	MetricsAddress string `mapstructure:"metricsAddress"`
	ProbeAddress   string `mapstructure:"probeAddress"`
	WebhookSecret  string `mapstructure:"webhookSecret"`
}

// Github define Github configuration
type Github struct {
	AppID          int64  `mapstructure:"appID"`
	InstallationID int64  `mapstructure:"installationID"`
	PrivateKeyPath string `mapstructure:"privateKeyPath"`
}

// Gitlab define Gitlab configuration
type Gitlab struct {
	Token string `mapstructure:"token"`
}

// RepoConfig define repository configuration
type RepoConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	TagRegexp string   `mapstructure:"tagRegexp"`
	Statuses  []string `mapstructure:"statuses"`
}
