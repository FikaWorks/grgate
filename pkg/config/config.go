package config

var (
	// CommitSha from source repository used to build GRGate
	CommitSha string

	// Main contains the generated main configuration
	Main *MainConfig

	// Version of GRGate
	Version string
)

// PlatformType is the type of platform to run against (Github or Gitlab)
type PlatformType string

const (
	// DefaultPlatform is the default platform
	DefaultPlatform PlatformType = GithubPlatform

	// DefaultRepoConfigPath is the default path of the .grgate config stored in
	// the repository
	DefaultRepoConfigPath string = ".grgate.yaml"

	// DefaultServerListenAddress is the default main server listening address
	DefaultServerListenAddress string = "0.0.0.0:8080"

	// DefaultServerMetricsAddress is the default metric server listening address
	DefaultServerMetricsAddress string = "0.0.0.0:9101"

	// DefaultServerProbeAddress is the default probe server listening address
	DefaultServerProbeAddress string = "0.0.0.0:8086"

	// DefaultWorkers defined the default amount of workers
	DefaultWorkers int = 5

	// GithubPlatform represent the Github platform
	GithubPlatform PlatformType = "github"

	// GitlabPlatform represent the Gitlab platform
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
