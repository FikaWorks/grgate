package utils

import (
	"bytes"
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
		data.ReleaseNote = strings.Trim(data.ReleaseNote[0:start]+data.ReleaseNote[end+len(releaseNoteMarkerEnd):], "\n")
	}

	var b bytes.Buffer
	err = t.Execute(&b, &data)
	output = b.String()
	return
}
