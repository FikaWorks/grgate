//go:build unit

package utils

import (
	"testing"

	"github.com/fikaworks/grgate/pkg/platforms"
	"github.com/kylelemons/godebug/pretty"
)

func TestRenderReleaseNoteAppendStatuses(t *testing.T) {
	data := &ReleaseNoteData{
		ReleaseNote: "This is a release note",
		Statuses: []*platforms.Status{
			{
				Name:   "e2e A",
				State: "success",
			},
			{
				Name:   "e2e B",
				State: "failed",
			},
			{
				Name:   "e2e C",
				State: "success",
			},
		},
	}

	template := `{{ .ReleaseNote }}
<!-- GRGate start -->
<details><summary>Status check</summary>
{{ range .Statuses }}
- [{{ if eq .State "success" }}x{{ else }} {{ end }}] {{ .Name }}
{{- end }}

</details>
<!-- GRGate end -->`

	expected := `This is a release note
<!-- GRGate start -->
<details><summary>Status check</summary>

- [x] e2e A
- [ ] e2e B
- [x] e2e C

</details>
<!-- GRGate end -->`

	result, err := RenderReleaseNote(template, data)
	if err != nil {
		t.Error("Error rendering release note", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("diff: (-got +want)\n%s", diff)
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
				State: "success",
			},
			{
				Name:   "e2e B",
				State: "success",
			},
			{
				Name:   "e2e C",
				State: "failed",
			},
		},
	}

	template := `{{- .ReleaseNote -}}
<!-- GRGate start -->
<details><summary>Status check</summary>
{{ range .Statuses }}
- [{{ if eq .State "success" }}x{{ else }} {{ end }}] {{ .Name }}
{{- end }}

</details>
<!-- GRGate end -->`

	expected := `This is a release note
<!-- GRGate start -->
<details><summary>Status check</summary>

- [x] e2e A
- [x] e2e B
- [ ] e2e C

</details>
<!-- GRGate end -->`

	result, err := RenderReleaseNote(template, data)
	if err != nil {
		t.Error("Error rendering release note", err)
	}
	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}

func TestMergeStatuses(t *testing.T) {
	expected := []*platforms.Status{
		{
			Name: "feature-flow",
		},
		{
			Name: "happy-flow",
		},
	}
	result := MergeStatuses([]*platforms.Status{
		{
			Name: "feature-flow",
		},
	}, []string{
		"happy-flow",
		"feature-flow",
	})

	if diff := pretty.Compare(result, expected); diff != "" {
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}
