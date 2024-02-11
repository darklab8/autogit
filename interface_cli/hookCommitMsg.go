/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package interface_cli

import (
	"github.com/darklab8/autogit/interface_cli/actions"

	"github.com/spf13/cobra"
)

// commitMsgCmd represents the commitMsg command
var commitMsgCmd = &cobra.Command{
	Use:   "commitMsg",
	Short: "FOR MACHINE USAGE ONLY: git hook for commit-msg. Not for human usage.",
	Run: func(cmd *cobra.Command, args []string) {
		shared.hook_commit_msg.Run()
		actions.CommmitMsg(args)
	},
}

func init() {
	hookCmd.AddCommand(commitMsgCmd)
}
