/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autogit/actions"
	"autogit/semanticgit/git"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "next semantic version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", actions.Version(versionParams, (&git.Repository{}).NewRepoInWorkDir()))
	},
}

var versionParams actions.ActionVersionParams

func init() {
	rootCmd.AddCommand(versionCmd)
	versionParams.Bind(versionCmd)
}
