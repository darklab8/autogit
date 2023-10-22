/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package interface_cli

import (
	"autogit/interface_cli/actions"
	"autogit/semanticgit/git"
	"autogit/settings"
	"fmt"

	"github.com/spf13/cobra"
)

// semverCmd represents the version command
var semverCmd = &cobra.Command{
	Use:     "semver",
	Short:   "next semantic version offered for your product",
	Aliases: []string{"version"},
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
