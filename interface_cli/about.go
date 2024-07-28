/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package interface_cli

import (
	"fmt"

	"github.com/darklab8/autogit/v2/interface_cli/actions"

	"github.com/spf13/cobra"
)

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:   "version",
	Short: "Info about autogit version",
	Run: func(cmd *cobra.Command, args []string) {
		shared.about.Run()
		fmt.Printf("%s", actions.About())
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aboutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// aboutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
