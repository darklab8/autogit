/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/actions"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "next semantic version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s", actions.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	actions.CMDversion.Init(versionCmd)
	actions.CMDversion.DisableNewLine = versionCmd.PersistentFlags().Bool("no-newline", false, "Disable newline")
}
