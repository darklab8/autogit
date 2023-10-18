/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autogit/semanticgit/git"
	"autogit/settings"
	"autogit/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

const NoSubsection = ""

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Shortcut activating hookPath from autogit.yml",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OK activate called")
		hook_folder := ".git-hooks"
		if *activateHookGLobally {
			hook_folder = filepath.Join(settings.UserHomeDir, hook_folder)
		}
		_ = os.Mkdir(hook_folder, 0777)
		commit_msg_path := filepath.Join(hook_folder, "commit-msg")
		ioutil.WriteFile(commit_msg_path, []byte("#!/bin/sh\n\nautogit hook commitMsg \"$1\"\n"), 0777)

		if !*activateHookGLobally {
			git := (&git.Repository{}).NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))
			git.HookEnabled(true)
		} else {
			cfg, err := config.LoadConfig(config.GlobalScope)
			utils.CheckFatal(err, "failed to read global scoped config")
			cfg.Raw.SetOption("core", NoSubsection, "hooksPath", hook_folder)
			utils.CheckFatal(cfg.Validate(), "failed to validate global config")
			file, err := cfg.Marshal()
			utils.CheckFatal(err, "failed to marshal global settings")
			fmt.Println("file", string(file))

			git_config_path := filepath.Join(settings.UserHomeDir, ".gitconfig")
			err = ioutil.WriteFile(git_config_path, file, 0777)
			utils.CheckFatal(err, "failed to write global settings")
		}
	},
}

var activateHookGLobally *bool

func init() {
	hookCmd.AddCommand(activateCmd)

	activateHookGLobally = activateCmd.Flags().BoolP("global", "g", false, "Init hook globally")
}
