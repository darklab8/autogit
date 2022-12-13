/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Shortcut to run git config --unset core.hooksPath",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deactivate called")
		executable, _ := exec.LookPath("git")
		command := exec.Cmd{
			Path:   executable,
			Args:   []string{"", "config", "--unset", "core.hooksPath"},
			Stdout: os.Stdout,
			Stderr: os.Stdout,
		}
		err := command.Run()

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deactivateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deactivateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
