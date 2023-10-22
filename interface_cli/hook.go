/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package interface_cli

import (
	"autogit/interface_cli/actions"
	"fmt"

	"github.com/spf13/cobra"
)

// hookCmd represents the hook command
var hookCmd = &cobra.Command{
	Use: "hook",
	Short: `git hooks entry points.
			Contains nested sub commands!
			` + actions.InitAdvice,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OK autogit hook called")
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
