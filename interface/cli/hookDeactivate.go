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
	"path/filepath"

	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Shortcut to run git config --unset core.hooksPath",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("OK deactivate called")
		if !*deactivateHookGlobally {
			git := (&git.Repository{}).NewRepoInWorkDir(git.SshPath(settings.GetConfig().Git.SSHPath))
			git.HookEnabled(false)
		} else {
			cfg, err := config.LoadConfig(config.GlobalScope)
			utils.CheckFatal(err, "failed to read global scoped config")
			for section_number, section := range cfg.Raw.Sections {
				if section.Name == "core" {
					cfg.Raw.Sections[section_number] = section.RemoveOption("hooksPath")
				}
			}
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

var deactivateHookGlobally *bool

func init() {
	hookCmd.AddCommand(deactivateCmd)

	deactivateHookGlobally = deactivateCmd.Flags().BoolP("global", "g", false, "deactivate hook globally")
}
