//go:build unit

package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGlobalConfig(t *testing.T) {
	t.Run("should return default global value if no config file is used",
		func(t *testing.T) {
			expectedValues := map[string]interface{}{
				"globals.enabled":              DefaultEnabled,
				"globals.releaseNote.enabled":  DefaultReleaseNoteEnabled,
				"globals.releaseNote.template": DefaultReleaseNoteTemplate,
				"globals.tagRegexp":            DefaultTagRegexp,
				"platform":                     DefaultPlatform,
				"repoConfigPath":               DefaultRepoConfigPath,
				"server.listenAddress":         DefaultServerListenAddress,
				"server.metricsAddress":        DefaultServerMetricsAddress,
				"server.probeAddress":          DefaultServerProbeAddress,
				"workers":                      DefaultWorkers,
			}

			v, err := NewGlobalConfig("")
			if err != nil {
				t.Errorf("Error not expected: %#v", err)
			}

			for key, expected := range expectedValues {
				if result := v.Get(key); result != expected {
					t.Errorf("Expected %#v, got %#v", expected, result)
				}
			}
		})

	t.Run("should override default settings if config file is provided",
		func(t *testing.T) {
			currentDir, err := os.Getwd()
			if err != nil {
				t.Errorf("Error not expected: %#v", err)
			}

			file, err := ioutil.TempFile(currentDir, "test-config.*.yaml")
			if err != nil {
				t.Errorf("Error not expected: %#v", err)
			}

			defer os.Remove(file.Name())

			if _, err := file.Write([]byte(`globals:
  enabled: true
  releaseNote:
    enabled: false
    template: |-
      some template
  tagRegexp: v\d*\.\d*\.\d*
platform: gitlab
`)); err != nil {
				t.Errorf("Error not expected: %#v", err)
			}

			expectedValues := map[string]interface{}{
				"globals.enabled":              DefaultEnabled,
				"globals.releaseNote.enabled":  false,
				"globals.releaseNote.template": "some template",
				"globals.tagRegexp":            "v\\d*\\.\\d*\\.\\d*",
				"platform":                     "gitlab",
				"repoConfigPath":               DefaultRepoConfigPath,
				"server.listenAddress":         DefaultServerListenAddress,
				"server.metricsAddress":        DefaultServerMetricsAddress,
				"server.probeAddress":          DefaultServerProbeAddress,
				"workers":                      DefaultWorkers,
			}

			v, err := NewGlobalConfig(file.Name())
			if err != nil {
				t.Errorf("Error not expected: %#v", err)
			}

			for key, expected := range expectedValues {
				if result := v.Get(key); result != expected {
					t.Errorf("Expected %#v, got %#v", expected, result)
				}
			}
		})
}
