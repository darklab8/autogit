/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/semanticgit/git"
	"fmt"

	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Shortcut to run git config --unset core.hooksPath",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OK deactivate called")
		git := (&git.Repository{}).NewRepoInWorkDir()
		git.HookEnabled(false)
	},
}

func init() {
	hookCmd.AddCommand(deactivateCmd)
}
