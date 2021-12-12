package utils

import (
	"bytes"
	"sort"
	"strings"
	"text/template"

	"github.com/fikaworks/grgate/pkg/config"
	"github.com/fikaworks/grgate/pkg/platforms"
)

var (
	releaseNoteMarkerStart = config.DefaultReleaseNoteMarkerStart
	releaseNoteMarkerEnd   = config.DefaultReleaseNoteMarkerEnd
)

// ReleaseNoteData hold release data used to populate the release note template
type ReleaseNoteData struct {
	ReleaseNote string
	Statuses    []*platforms.Status
}

// RenderReleaseNote add/update status check from a release note  based on a
// template
func RenderReleaseNote(tpl string, data *ReleaseNoteData) (output string, err error) {
	t, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return
	}

	start := strings.Index(data.ReleaseNote, releaseNoteMarkerStart)
	end := strings.Index(data.ReleaseNote, releaseNoteMarkerEnd)

	// if markers already exist, remove content and render the template so it
	// looks like status check have been updated
	if start > -1 && end > -1 {
		data.ReleaseNote = data.ReleaseNote[0:start] + data.ReleaseNote[end+len(releaseNoteMarkerEnd):]
	}

	var b bytes.Buffer
	err = t.Execute(&b, &data)
	output = b.String()
	return
}

// MergeStatuses based on the repo/global config
func MergeStatuses(platformStatuses []*platforms.Status, configStatuses []string) []*platforms.Status {
	for _, configStatus := range configStatuses {
		present := false
		for _, platformStatus := range platformStatuses {
			if platformStatus.Name == configStatus {
				present = true
				break
			}
		}
		if !present {
			platformStatuses = append(platformStatuses, &platforms.Status{
				Name: configStatus,
			})
		}
	}

	// sort status by name
	sort.Slice(platformStatuses, func(i1, i2 int) bool {
		return platformStatuses[i1].Name < platformStatuses[i2].Name
	})

	return platformStatuses
}
