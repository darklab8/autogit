{{.Header}}

{{ if .AreThereFeatures }}### Features{{ end }}

{{range $val := .Features}}{{$val}}{{end}}
{{range $key, $list := .FeaturesScoped}}#### {{$key}}
{{range $val := $list}}{{$val}}{{end}}{{end}}

{{ if .AreThereFixes }}### Bug Fixes{{ end }}

{{range $val := .Fixes}}{{$val}}{{end}}
{{range $key, $list := .FixesScoped}}#### {{$key}}
{{range $val := $list}}{{$val}}{{end}}{{end}}