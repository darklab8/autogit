/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autogit/settings"
	"autogit/settings/logus"
	"autogit/utils"

	"github.com/spf13/cobra"
)

const (
	initAdvice string = "activate hook with `autogit hook activate [--global]`"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init repository settings. " + initAdvice,
	Run: func(cmd *cobra.Command, args []string) {
		config_path := settings.ProjectConfigPath
		if *initGLobally {
			config_path = settings.GlobSettingPath
		}

		if utils.FileExists(string(config_path)) {
			logus.Fatal("file with settings already exists", logus.ConfigPath(config_path))
			return
		}

		file := utils.NewFile(string(config_path)).CreateToWriteF()
		defer file.Close()
		file.WritelnF(settings.ConfigExample)

		logus.Info("Succesfully created settings in location", logus.ConfigPath(config_path))
		logus.Info("Try to " + initAdvice + ". It will automatically verify committs for you!")
	},
}

var initGLobally *bool

func init() {
	rootCmd.AddCommand(initCmd)
	initGLobally = initCmd.Flags().BoolP("global", "g", false, "Init settings file globally")
}
