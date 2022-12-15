/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/actions"
	"autogit/utils"
	"fmt"

	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Shortcut to 'git tag $(autogit version)'",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OK version tag is called")
		RemoveNewLIne := true
		actions.CMDversion.DisableNewLine = &RemoveNewLIne
		vers := actions.Version()
		utils.ShellRunArgs("git", "tag", vers)
		fmt.Printf("OK tag=%s is created\n", vers)
	},
}

func init() {
	versionCmd.AddCommand(tagCmd)
}
