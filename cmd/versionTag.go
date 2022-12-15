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
		actions.VersionDisableNewLine = &RemoveNewLIne
		vers := actions.Version()
		utils.ShellRunArgs("git", "tag", vers)
		fmt.Printf("OK tag=%s is created\n", vers)
	},
}

func init() {
	versionCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
