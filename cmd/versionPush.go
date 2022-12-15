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

// tagPushCmd represents the tagPush command
var tagPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Shortcut to 'git push origin $(autogit version)'",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OK version push is called")
		RemoveNewLIne := true
		actions.CMDversion.DisableNewLine = &RemoveNewLIne
		vers := actions.Version()
		utils.ShellRunArgs("git", "push", "origin", vers)
		fmt.Printf("OK tag=%s is pushed\n", vers)
	},
}

func init() {
	versionCmd.AddCommand(tagPushCmd)
}
