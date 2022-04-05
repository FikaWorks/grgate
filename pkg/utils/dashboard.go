package utils

import (
	"bytes"
	"text/template"
)

// DashboardData hold issue data used to populate the issue dashboard template
type DashboardData struct {
	Errors  []string
	Enabled bool
}

func RenderDashboard(tpl string, data *DashboardData) (output string, err error) {
	t, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return
	}

	var b bytes.Buffer
	err = t.Execute(&b, &data)
	output = b.String()
	return
}
