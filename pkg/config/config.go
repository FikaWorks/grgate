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
	// DefaultEnabled is the default value which define if GRGate should process
	// the repository or not
	DefaultEnabled bool = true

	// DefaultDashboardEnabled define if the issue dashboard your be enabled to
	// provide feedback on the state of GRGate
	DefaultDashboardEnabled bool = true

	// DefaultDashboardAuthor define the title of the issue dashboard
	DefaultDashboardAuthor string = "GRGate[bot]"

	// DefaultDashboardTitle define the title of the issue dashboard
	DefaultDashboardTitle string = "GRGate dashboard"

	// DefaultDashboardTemplate define the default template used to display
	// the state of GRGate in a GitHub/GitLab issue
	DefaultDashboardTemplate string = `GRGate is {{ if .Enabled }}enabled
{{- else }}disabled{{ end }} for this repository.
{{- if .Errors }}

Incorrect configuration detected with the following error(s):
{{- range .Errors }}
- {{ . }}
{{- end }}
{{- end }}

Last time GRGate processed this repository: {{ .LastExecutionTime }}`

	// DefaultTagRegexp is the default pattern used to match tags attached to
	// releases
	DefaultTagRegexp string = ".*"

	// DefaultPlatform is the default platform
	DefaultPlatform PlatformType = GithubPlatform

	// DefaultReleaseNoteEnabled define if the statuses should be added to the
	// release note
	DefaultReleaseNoteEnabled bool = true

	// DefaultReleaseNoteMarkerStart is the string that define the start of the
	// section appended to the release note
	DefaultReleaseNoteMarkerStart string = "<!-- GRGate start -->"

	// DefaultReleaseNoteMarkerEnd is the string that define the end of the
	// section appended to the release note
	DefaultReleaseNoteMarkerEnd string = "<!-- GRGate end -->"

	// DefaultReleaseNoteTemplate define the default template used to display
	// statuses in the release note
	DefaultReleaseNoteTemplate string = `{{- .ReleaseNote -}}
<!-- GRGate start -->
<details><summary>Status check</summary>
{{ range .Statuses }}
- [{{ if or (eq .Status "completed" ) (eq .Status "success") }}x{{ else }} {{ end }}] {{ .Name }}
{{- end }}

</details>
<!-- GRGate end -->`

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

// Dashboard define the issue dashboard configuration
type Dashboard struct {
	Enabled  bool   `mapstructure:"enabled"`
	Author   string `mapstructure:"author"`
	Title    string `mapstructure:"title"`
	Template string `mapstructure:"template"`
}

// ReleaseNote define the release note configuration
type ReleaseNote struct {
	Enabled  bool   `mapstructure:"enabled"`
	Template string `mapstructure:"template"`
}

// RepoConfig define repository configuration
type RepoConfig struct {
	Enabled     bool         `mapstructure:"enabled"`
	Dashboard   *Dashboard   `mapstructure:"dashboard"`
	ReleaseNote *ReleaseNote `mapstructure:"releaseNote"`
	Statuses    []string     `mapstructure:"statuses"`
	TagRegexp   string       `mapstructure:"tagRegexp"`
}

// Server define server configuration
type Server struct {
	ListenAddress  string `mapstructure:"listenAddress"`
	MetricsAddress string `mapstructure:"metricsAddress"`
	ProbeAddress   string `mapstructure:"probeAddress"`
	WebhookSecret  string `mapstructure:"webhookSecret"`
}
