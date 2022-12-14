/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/actions"
	"fmt"

	"github.com/spf13/cobra"
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Auto generated changelog according to git conventional commits",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", actions.Changelog())
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)

	actions.ChangelogTag = changelogCmd.PersistentFlags().String("tag", "", "Select from which tag")
	actions.ChangelogValidate = changelogCmd.PersistentFlags().Bool("validate", false, "Validate to rules")
}
