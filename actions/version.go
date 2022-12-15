package actions

import (
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"strings"

	"github.com/spf13/cobra"
)

type VersionParams struct {
	DisableVFlag *bool
	BuildMeta    *string
	Alpha        *bool
	Beta         *bool
	Prerelease   *bool
	Publish      *bool
}

func (v *VersionParams) Init(cmd *cobra.Command) {
	v.DisableVFlag = cmd.PersistentFlags().Bool("no-v", false, "Disable v flag")
	v.BuildMeta = cmd.PersistentFlags().String("build", "", "Build metadata, not affecting semantic versioning. Added as semver+build")
	v.Alpha = cmd.PersistentFlags().Bool("alpha", false, "Enable next version as alpha")
	v.Beta = cmd.PersistentFlags().Bool("beta", false, "Enable next version as beta")
	v.Prerelease = cmd.PersistentFlags().Bool("rc", false, "Enable next version as prerelease")
	v.Publish = cmd.PersistentFlags().Bool("publish", false, "Breaking from 0.x.x to 1+.x.x versions")
}

var CMDversion struct {
	VersionParams
	DisableNewLine *bool
}

func Version() string {
	g := (&semanticgit.SemanticGit{}).NewRepo((&git.Repository{}).NewRepoInWorkDir())

	semver_options := semver.OptionsSemVer{
		DisableVFlag:  *CMDversion.DisableVFlag,
		EnableNewline: !(*CMDversion.DisableNewLine),
		Build:         *CMDversion.BuildMeta,
		Alpha:         *CMDversion.Alpha,
		Beta:          *CMDversion.Beta,
		Rc:            *CMDversion.Prerelease,
		Publish:       *CMDversion.Publish,
	}
	vers := g.GetNextVersion(semver_options)

	var sb strings.Builder
	sb.WriteString(vers.ToString())
	return sb.String()
}
