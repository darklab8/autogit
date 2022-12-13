/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/settings"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Shortcut activating hookPath from autogit.yml",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("activate called")
		fmt.Printf("hookPath=%s\n", settings.Config.HookPath)
		executable, _ := exec.LookPath("git")
		command := exec.Cmd{
			Path:   executable,
			Args:   []string{"", "config", "core.hooksPath", settings.Config.HookPath},
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
	rootCmd.AddCommand(activateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// activateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// activateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
