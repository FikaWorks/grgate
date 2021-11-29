package utils

import (
	"testing"

	"github.com/fikaworks/grgate/pkg/platforms"
)

func TestRenderReleaseNoteAppendStatuses(t *testing.T) {
	data := &ReleaseNoteData{
		ReleaseNote: "This is a release note",
		Statuses: []*platforms.Status{
			{
				Name:   "e2e A",
				Status: "success",
			},
			{
				Name:   "e2e B",
				Status: "failed",
			},
		},
	}

	template := `{{ .ReleaseNote }}

<!-- GRGate start -->
<details><summary>Status check</summary>
{{- range .Statuses }}
- [{{ if eq .Status "success" }}x{{ else }} {{ end }}] {{ .Name }}
{{- end }}
</details>
<!-- GRGate end -->`

	expected := `This is a release note

<!-- GRGate start -->
<details><summary>Status check</summary>
- [x] e2e A
- [ ] e2e B
</details>
<!-- GRGate end -->`

	result, err := RenderReleaseNote(template, data)
	if err != nil {
		t.Error("Error rendering release note", err)
	}
	if result != expected {
		t.Errorf("Render release note got %s, expected %s", result, expected)
	}
}

func TestRenderReleaseNoteEditStatuses(t *testing.T) {
	data := &ReleaseNoteData{
		ReleaseNote: `This is a release note

<!-- GRGate start -->
<details><summary>Status check</summary>
- [ ] e2e A
- [ ] e2e B
</details>
<!-- GRGate end -->`,
		Statuses: []*platforms.Status{
			{
				Name:   "e2e A",
				Status: "success",
			},
			{
				Name:   "e2e B",
				Status: "success",
			},
		},
	}

	template := `{{ .ReleaseNote }}

<!-- GRGate start -->
<details><summary>Status check</summary>
{{- range .Statuses }}
- [{{ if eq .Status "success" }}x{{ else }} {{ end }}] {{ .Name }}
{{- end }}
</details>
<!-- GRGate end -->`

	expected := `This is a release note

<!-- GRGate start -->
<details><summary>Status check</summary>
- [x] e2e A
- [x] e2e B
</details>
<!-- GRGate end -->`

	result, err := RenderReleaseNote(template, data)
	if err != nil {
		t.Error("Error rendering release note", err)
	}
	if result != expected {
		t.Errorf("Render release note got %s, expected %s", result, expected)
	}
}
