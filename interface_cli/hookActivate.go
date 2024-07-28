/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package interface_cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/darklab8/autogit/v2/semanticgit/git"
	"github.com/darklab8/autogit/v2/settings"
	"github.com/darklab8/autogit/v2/settings/envs"
	"github.com/darklab8/autogit/v2/settings/logus"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

const NoSubsection = ""

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Shortcut activating hookPath from autogit.yml",
	Run: func(cmd *cobra.Command, args []string) {
		shared.hook_activate.Run()
		fmt.Println("OK activate called")
		hook_folder := settings.HookFolderName
		if *activateHookGLobally {
			hook_folder = filepath.Join(string(envs.PathUserHome), hook_folder)
		}
		_ = os.Mkdir(hook_folder, 0777)
		commit_msg_path := filepath.Join(hook_folder, "commit-msg")

		verbose_propagating_cmd := ""
		if *(shared.hook_activate.verboseLogging) {
			verbose_propagating_cmd = "-v "
		}
		os.WriteFile(commit_msg_path, []byte(fmt.Sprintf("#!/bin/sh\n\nautogit hook commitMsg %s\"$1\"\n", verbose_propagating_cmd)), 0777)

		if !*activateHookGLobally {
			git := git.NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))
			git.HookEnabled(true)
		} else {
			cfg, err := config.LoadConfig(config.GlobalScope)
			logus.Log.CheckFatal(err, "failed to read global scoped config")
			cfg.Raw.SetOption("core", NoSubsection, "hooksPath", hook_folder)
			logus.Log.CheckFatal(cfg.Validate(), "failed to validate global config")
			file, err := cfg.Marshal()
			logus.Log.CheckFatal(err, "failed to marshal global settings")
			fmt.Println("file", string(file))

			err = os.WriteFile(string(envs.PathGitConfig), file, 0777)
			logus.Log.CheckFatal(err, "failed to write global settings")
		}
	},
}

var activateHookGLobally *bool

func init() {
	hookCmd.AddCommand(activateCmd)

	activateHookGLobally = activateCmd.Flags().BoolP("global", "g", false, "Init hook globally")
}
