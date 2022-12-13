/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/settings"
	"autogit/utils"
	"fmt"

	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Shortcut activating hookPath from autogit.yml",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("activate called")
		fmt.Printf("hookPath=%s\n", settings.Config.HookPath)
		utils.ShellRunArgs("git", "config", "core.hooksPath", settings.Config.HookPath)
	},
}

func init() {
	hookCmd.AddCommand(activateCmd)
}
