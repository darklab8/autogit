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

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Auto generated changelog according to git conventional commits",
	Run: func(cmd *cobra.Command, args []string) {
		shared.changelog.Run()
		fmt.Printf("%s", actions.Changelog(paramsChangelog, git.NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))))
	},
}

var paramsChangelog actions.ChangelogParams

func init() {
	rootCmd.AddCommand(changelogCmd)
	paramsChangelog.Bind(changelogCmd)
}
