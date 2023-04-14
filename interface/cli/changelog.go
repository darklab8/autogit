/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autogit/actions"
	"autogit/semanticgit/git"
	"autogit/settings"
	"fmt"

	"github.com/spf13/cobra"
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Auto generated changelog according to git conventional commits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", actions.Changelog(paramsChangelog, (&git.Repository{}).NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))))
	},
}

var paramsChangelog actions.ChangelogParams

func init() {
	rootCmd.AddCommand(changelogCmd)
	paramsChangelog.Bind(changelogCmd)
}
