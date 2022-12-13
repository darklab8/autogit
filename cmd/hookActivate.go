/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"autogit/utils"
	"fmt"
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
		_ = os.Mkdir(hook_folder, os.ModePerm)
		commit_msg_path := filepath.Join(hook_folder, "commit-msg")
		file := utils.File{Filepath: commit_msg_path}
		file.CreateToWriteF()
		file.WritelnF(`#!/bin/sh`)
		file.WritelnF(``)
		file.WritelnF(`autogit hook commitMsg "$1"`)
		file.Close()
		utils.ShellRunArgs("git", "config", "core.hooksPath", hook_folder)
	},
}

func init() {
	hookCmd.AddCommand(activateCmd)
}
