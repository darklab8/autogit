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
	Use:   "tagPush",
	Short: "Shortcut to 'git push origin $(autogit version)'",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OK version push is called")
		vers := actions.Version()
		utils.ShellRunArgs("git", "push", "origin", vers)
		fmt.Printf("OK tag=%s is pushed\n", vers)
	},
}

func init() {
	versionCmd.AddCommand(tagPushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagPushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagPushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
