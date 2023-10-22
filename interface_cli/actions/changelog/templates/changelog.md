{{.Header}}
{{range $semver_group := .OrderedSemverGroups }}
## {{ $semver_group.Name }} Changes
{{range $commit_type, $type_group := $semver_group.CommitTypeGroups -}}
### {{ $commit_type }}
{{range $commit := $type_group.NoScopeCommits}}{{$commit}}{{end}}
{{- range $scope, $commits := $type_group.ScopedCommits -}}* {{$scope}}
{{range $commit := $commits}}   {{$commit}}{{end -}}
{{ end -}}
{{ end -}}
{{ end -}}
