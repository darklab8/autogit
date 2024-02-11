/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package interface_cli

import (
	"fmt"

	"github.com/darklab8/autogit/interface_cli/actions"
	"github.com/darklab8/autogit/semanticgit/git"
	"github.com/darklab8/autogit/settings"

	"github.com/spf13/cobra"
)

// semverCmd represents the version command
var semverCmd = &cobra.Command{
	Use:   "semver",
	Short: "next semantic version offered for your product",
	Run: func(cmd *cobra.Command, args []string) {
		shared.version.Run()
		fmt.Printf("%s", actions.Version(versionParams, git.NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))))
	},
}

var versionParams actions.ActionVersionParams

func init() {
	rootCmd.AddCommand(semverCmd)
	versionParams.Bind(semverCmd)
}
