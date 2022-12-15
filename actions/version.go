package actions

import (
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"strings"
)

var VersionDisableNewLine *bool
var VersionDisableVFlag *bool
var VersionBuildMeta *string
var VersionAlpha *bool
var VersionBeta *bool
var VersionPrerelease *bool

// if one is enabled, only this one is incremented.
// if all enabled? we increment only most volatile one only
// when we read latest tag, ->
// we read history up until latest non-prerelease, while collecting alpha/beta/prerelease latest versions
// if we enabled only beta? it will add version +1 based on latest beta.
// if we enabled beta+prerelease, it will not change beta, beta will be rendered, prerelease version will be increased

func Version() string {
	g := (&semanticgit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())

	semver_options := semver.OptionsSemVer{
		DisableVFlag:  *VersionDisableVFlag,
		EnableNewline: !(*VersionDisableNewLine),
		Build:         *VersionBuildMeta,
		Alpha:         *VersionAlpha,
		Beta:          *VersionBeta,
		Rc:            *VersionPrerelease,
	}
	vers := g.GetNextVersion(semver_options)

	var sb strings.Builder
	sb.WriteString(vers.ToString())
	return sb.String()
}
