package changelog_types

type ChangelogSectionName string

type ChangelogSectionType string

const (
	SemVerMajor  ChangelogSectionType = "semver_major"
	SemVerMinor  ChangelogSectionType = "semver_minor"
	SemVerPatch  ChangelogSectionType = "semver_patch"
	MergeCommits ChangelogSectionType = "merge_commits"
)
