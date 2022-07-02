//go:build unit

package utils

import (
	"testing"
	"time"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/kylelemons/godebug/pretty"
)

func TestRenderDashboard(t *testing.T) {
	currentTime := time.Now().UTC().Format(time.UnixDate)
	data := &DashboardData{
		Enabled:           true,
		Errors:            []string{},
		LastExecutionTime: currentTime,
	}

	expected := `GRGate is enabled for this repository.

Last time GRGate processed this repository: ` + currentTime

	result, err := RenderDashboard(config.DefaultDashboardTemplate, data)
	if err != nil {
		t.Errorf("Error rendering dashboard: %#v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}

func TestRenderDashboardError(t *testing.T) {
	currentTime := time.Now().UTC().Format(time.UnixDate)
	data := &DashboardData{
		Enabled:           true,
		Errors:            []string{"error 1", "error 2"},
		LastExecutionTime: currentTime,
	}

	expected := `GRGate is enabled for this repository.

Incorrect configuration detected with the following error(s):
- error 1
- error 2

Last time GRGate processed this repository: ` + currentTime

	result, err := RenderDashboard(config.DefaultDashboardTemplate, data)
	if err != nil {
		t.Errorf("Error rendering dashboard: %#v", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}
