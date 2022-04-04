//go:build unit

package utils

import (
	"testing"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/kylelemons/godebug/pretty"
)

func TestRenderDashboard(t *testing.T) {
	data := &DashboardData{
		Enabled: true,
		Errors:  []string{},
	}

	expected := "GRGate is enabled for this repository."

	result, err := RenderDashboard(config.DefaultDashboardTemplate, data)
	if err != nil {
		t.Errorf("Error rendering dashboard: %#v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}

func TestRenderDashboardError(t *testing.T) {
	data := &DashboardData{
		Enabled: true,
		Errors:  []string{"error 1", "error 2"},
	}

	expected := `GRGate is enabled for this repository.

Incorrect configuration detected with the following error(s):
- error 1
- error 2`

	result, err := RenderDashboard(config.DefaultDashboardTemplate, data)
	if err != nil {
		t.Errorf("Error rendering dashboard: %#v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}
