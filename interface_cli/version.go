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

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "next semantic version",
	Run: func(cmd *cobra.Command, args []string) {
		shared.version.Run()
		fmt.Printf("%s", actions.Version(versionParams, git.NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))))
	},
}

var versionParams actions.ActionVersionParams

func init() {
	rootCmd.AddCommand(versionCmd)
	versionParams.Bind(versionCmd)
}
