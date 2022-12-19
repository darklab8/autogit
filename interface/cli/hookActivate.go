/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autogit/semanticgit/git"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Shortcut activating hookPath from autogit.yml",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OK activate called")
		hook_folder := ".git-hooks"
		_ = os.Mkdir(hook_folder, 0777)
		commit_msg_path := filepath.Join(hook_folder, "commit-msg")
		ioutil.WriteFile(commit_msg_path, []byte("#!/bin/sh\n\nautogit hook commitMsg \"$1\"\n"), 0777)

		git := (&git.Repository{}).NewRepoInWorkDir()
		git.HookEnabled(true)
	},
}

func init() {
	hookCmd.AddCommand(activateCmd)
}
