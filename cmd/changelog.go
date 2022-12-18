/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/actions"
	"autogit/semanticgit/git"
	"fmt"

	"github.com/spf13/cobra"
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Auto generated changelog according to git conventional commits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", actions.Changelog(paramsChangelog, (&git.Repository{}).NewRepoInWorkDir()))
	},
}

var paramsChangelog actions.ChangelogParams

func init() {
	rootCmd.AddCommand(changelogCmd)
	paramsChangelog.Bind(changelogCmd)
}
