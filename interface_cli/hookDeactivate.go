/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package interface_cli

import (
	"autogit/semanticgit/git"
	"autogit/settings"
	"autogit/settings/envs"
	"autogit/settings/logus"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Shortcut to run git config --unset core.hooksPath",
	Run: func(cmd *cobra.Command, args []string) {
		shared.hook_deactivate.Run()
		fmt.Println("OK deactivate called")
		if !*deactivateHookGlobally {
			git := git.NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))
			git.HookEnabled(false)
		} else {
			cfg, err := config.LoadConfig(config.GlobalScope)
			logus.Log.CheckFatal(err, "failed to read global scoped config")
			for section_number, section := range cfg.Raw.Sections {
				if section.Name == "core" {
					cfg.Raw.Sections[section_number] = section.RemoveOption("hooksPath")
				}
			}
			logus.Log.CheckFatal(cfg.Validate(), "failed to validate global config")
			file, err := cfg.Marshal()
			logus.Log.CheckFatal(err, "failed to marshal global settings")
			fmt.Println("file", string(file))

			err = os.WriteFile(string(envs.PathGitConfig), file, 0777)
			logus.Log.CheckFatal(err, "failed to write global settings")
		}
	},
}

var deactivateHookGlobally *bool

func init() {
	hookCmd.AddCommand(deactivateCmd)

	deactivateHookGlobally = deactivateCmd.Flags().BoolP("global", "g", false, "deactivate hook globally")
}
