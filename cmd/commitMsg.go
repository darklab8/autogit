/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/actions"

	"github.com/spf13/cobra"
)

// commitMsgCmd represents the commitMsg command
var commitMsgCmd = &cobra.Command{
	Use:   "commitMsg",
	Short: "git hook for commit-msg. Not for human usage.",
	Run: func(cmd *cobra.Command, args []string) {
		actions.CommmitMsg(args)
	},
}

func init() {
	rootCmd.AddCommand(commitMsgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitMsgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitMsgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
