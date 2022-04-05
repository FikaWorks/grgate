//go:build integration || integrationgithub || integrationgitlab

package tests

import (
	"github.com/fikaworks/grgate/pkg/platforms"
)

type testCase struct {
	expectErrorDuringProcess bool
	expectIssueToBeCreated   bool
	expectedDashboardBody    string
	expectedDashboardTitle   string
	expectPublishedRelease   bool
	expectedReleaseNote      string
	withRepoConfig           string
	withStatuses             []*platforms.Status
	withTag                  string
}

var (
	disabledConfigTestCases = map[string]*testCase{
		"should not process the release when disabled by config": {
			withRepoConfig:         "enabled: false",
			withTag:                "v1.2.3",
			expectPublishedRelease: false,
		},
	}

	commitStatusTestCases = map[string]*testCase{
		"should not publish release when commit status are not defined": {
			withRepoConfig: `enabled: true
tagRegexp: v\d*\.\d*\.\d*
releaseNote:
  enabled: false`,
			withStatuses:           []*platforms.Status{},
			withTag:                "v1.2.3",
			expectPublishedRelease: false,
		},
		"should not publish release when commit status are still running": {
			withRepoConfig: `enabled: true
tagRegexp: v\d*\.\d*\.\d*
releaseNote:
  enabled: false
statuses:
- e2e-happyflow
- e2e-featureflow`,
			withStatuses: []*platforms.Status{
				{
					Name:   "e2e-happyflow",
					Status: "in_progress",
				},
			},
			withTag:                "v1.2.3",
			expectPublishedRelease: false,
		},
		"should publish release when all status succeeded": {
			withRepoConfig: `enabled: true
tagRegexp: v\d*\.\d*\.\d*
releaseNote:
  enabled: false
statuses:
- e2e-happyflow
- e2e-featureflow`,
			withStatuses: []*platforms.Status{
				{
					Name:   "e2e-happyflow",
					State:  "success",
					Status: "completed",
				}, {
					Name:   "e2e-featureflow",
					State:  "success",
					Status: "completed",
				},
			},
			withTag:                "v1.2.3",
			expectPublishedRelease: true,
		},
	}

	releaseNoteTestCases = map[string]*testCase{
		"should update release note with statuses": {
			withRepoConfig: `enabled: true
tagRegexp: v\d*\.\d*\.\d*
releaseNote:
  enabled: true
  template: |-
    {{- .ReleaseNote -}}
    <!-- GRGate start -->
    <details><summary>GRGate status check</summary>
    {{ range .Statuses }}
    - [{{ if or (eq .Status "completed" ) (eq .Status "success") }}x{{ else }} {{ end }}] {{ .Name }}
    {{- end }}

    </details>
    <!-- GRGate end -->
statuses:
- e2e-happyflow
- e2e-featureflow-a
- e2e-featureflow-b`,
			withStatuses:           []*platforms.Status{},
			withTag:                "v1.2.3",
			expectPublishedRelease: false,
			expectedReleaseNote: `<!-- GRGate start -->
<details><summary>GRGate status check</summary>

- [ ] e2e-featureflow-a
- [ ] e2e-featureflow-b
- [ ] e2e-happyflow

</details>
<!-- GRGate end -->`,
		},
		"should publish release if all status succeeded": {
			withRepoConfig: `enabled: true
tagRegexp: v\d*\.\d*\.\d*
releaseNote:
  enabled: true
  template: |-
    {{- .ReleaseNote -}}
    <!-- GRGate start -->
    <details><summary>GRGate status check</summary>
    {{ range .Statuses }}
    - [{{ if or (eq .Status "completed" ) (eq .Status "success") }}x{{ else }} {{ end }}] {{ .Name }}
    {{- end }}

    </details>
    <!-- GRGate end -->
statuses:
- e2e-happyflow
- e2e-featureflow-a
- e2e-featureflow-b`,
			withStatuses: []*platforms.Status{
				{
					Name:   "e2e-happyflow",
					State:  "success",
					Status: "completed",
				}, {
					Name:   "e2e-featureflow-a",
					State:  "success",
					Status: "completed",
				}, {
					Name:   "e2e-featureflow-b",
					State:  "success",
					Status: "completed",
				},
			},
			withTag:                "v1.2.3",
			expectPublishedRelease: true,
			expectedReleaseNote: `<!-- GRGate start -->
<details><summary>GRGate status check</summary>

- [x] e2e-featureflow-a
- [x] e2e-featureflow-b
- [x] e2e-happyflow

</details>
<!-- GRGate end -->`,
		},
	}

	dashboardTestCases = map[string]*testCase{
		"should not create issue dashboard when dashboard is disabled": {
			withRepoConfig: `enabled: false
dashboard:
  enabled: false
statuses:
- happy flow`,
			withTag:                  "v1.2.3",
			expectIssueToBeCreated:   false,
			expectErrorDuringProcess: false,
		},
		"should create issue dashboard": {
			withRepoConfig: `enabled: false
tagRegexp: v\d*\.\d*\.\d*
dashboard:
  enabled: true
  title: GRGate dashboard
  template: GRGate is enabled
statuses:
- happy flow`,
			withTag:                  "v1.2.3",
			expectIssueToBeCreated:   true,
			expectErrorDuringProcess: false,
			expectedDashboardTitle:   "GRGate dashboard",
			expectedDashboardBody:    "GRGate is enabled",
		},
		"should report back to dashboard when statuses are not defined": {
			withRepoConfig: `enabled: true
tagRegexp: v\d*\.\d*\.\d*
dashboard:
  enabled: true
  title: GRGate dashboard
  template: |-
    GRGate is {{ if .Enabled }}enabled
    {{- else }}disabled{{ end }} for this repository.
    {{- if .Errors }}

    Incorrect configuration detected with the following error(s):
    {{- range .Errors }}
    - {{ . }}
    {{- end }}
    {{- end }}`,
			withTag:                  "v1.2.3",
			expectIssueToBeCreated:   true,
			expectErrorDuringProcess: false,
			expectedDashboardTitle:   "GRGate dashboard",
			expectedDashboardBody: `GRGate is enabled for this repository.

Incorrect configuration detected with the following error(s):
- Statuses are undefined in .grgate.yaml`,
		},
		"should report back to dashboard when tag regexp cannot be parsed": {
			withRepoConfig: `enabled: true
tagRegexp: "[["
dashboard:
  enabled: true
  title: GRGate dashboard
  template: |-
    GRGate is {{ if .Enabled }}enabled
    {{- else }}disabled{{ end }} for this repository.
    {{- if .Errors }}

    Incorrect configuration detected with the following error(s):
    {{- range .Errors }}
    - {{ . }}
    {{- end }}
    {{- end }}
statuses:
- e2e happy flow`,
			withTag:                  "v1.2.3",
			expectIssueToBeCreated:   true,
			expectErrorDuringProcess: true,
			expectedDashboardTitle:   "GRGate dashboard",
			expectedDashboardBody: `GRGate is enabled for this repository.

Incorrect configuration detected with the following error(s):
- Couldn't compile regexp "[["`,
		},
	}
)
