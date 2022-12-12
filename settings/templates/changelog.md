{{.Header}}

{{ if .Features }}### Features{{ end }}

{{range $val := .Features}}{{$val}}{{end}}
{{ if .Fixes }}### Bug Fixes{{ end }}

{{range $val := .Fixes}}{{$val}}{{end}}