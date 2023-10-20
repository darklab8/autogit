package actions

import (
	"autogit/semanticgit"
	"autogit/semanticgit/git"
	"autogit/semanticgit/semver"
	"autogit/settings"
	"autogit/settings/types"

	"github.com/spf13/cobra"
)

type GitActions struct {
	Tag  bool
	Push bool
}

type VersionParams struct {
	semver.OptionsSemVer
}

func (v *VersionParams) Bind(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&v.DisableVFlag, "no-v", false, "Disable v flag")
	cmd.PersistentFlags().StringVar(&v.Build, "build", "", "Build metadata, not affecting semantic versioning. Added as semver+build")
	cmd.PersistentFlags().BoolVar(&v.Alpha, "alpha", false, "Enable next version as alpha")
	cmd.PersistentFlags().BoolVar(&v.Beta, "beta", false, "Enable next version as beta")
	cmd.PersistentFlags().BoolVar(&v.Rc, "rc", false, "Enable next version as prerelease")
	cmd.PersistentFlags().BoolVar(&v.Publish, "publish", false, "Breaking from 0.x.x to 1+.x.x versions")
}

type ActionVersionParams struct {
	VersionParams
	GitActions
}

func (v *ActionVersionParams) Bind(cmd *cobra.Command) {
	v.VersionParams.Bind(cmd)
	cmd.PersistentFlags().BoolVar(&v.EnableNewline, "newline", true, "Newline pressence, disable with --newline=false")
	cmd.PersistentFlags().BoolVar(&v.Tag, "tag", false, "shortcut to 'git -a tag -m $(autogit changelog)', not requiring installed git")
	cmd.PersistentFlags().BoolVar(&v.Push, "push", false, "shortcut to 'git push', not requiring installed git")
}

// gitw - (&git.Repository{}).NewRepoInWorkDir() for cmd
// we can overrise with git in memory
func Version(params ActionVersionParams, gitw *git.Repository) types.TagName {
	settings.LoadSettings(settings.GetSettingsPath())

	g := (&semanticgit.SemanticGit{}).NewRepo(gitw)
	vers := g.GetNextVersion(params.OptionsSemVer)
	renderedVers := vers.ToString()

	vers.Options.EnableNewline = false
	changelog := Changelog(ChangelogParams{VersionParams: params.VersionParams}, gitw)
	if params.Tag {
		gitw.CreateTag(vers.ToString(), changelog)
	}

	if params.Push {
		gitw.PushTag(vers.ToString())
	}

	return renderedVers
}
