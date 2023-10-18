/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"autogit/settings"
	"autogit/utils"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

const (
	initAdvice string = "activate hook with `autogit hook activate`"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init repository settings. " + initAdvice,
	Run: func(cmd *cobra.Command, args []string) {
		config_path := "autogit.yml"
		if *initGLobally {
			config_path = settings.GlobSettingPath
		}

		if utils.FileExists(config_path) {
			log.Fatalln("file with settings=", config_path, " already exists")
			return
		}

		file := utils.NewFile(config_path).CreateToWriteF()
		defer file.Close()
		file.WritelnF(settings.ConfigExample)

		fmt.Println("Succesfully created autogit.yml in location", config_path)
		fmt.Println("Try to " + initAdvice + ". It will automatically verify committs for you!")
	},
}

var initGLobally *bool

func init() {
	rootCmd.AddCommand(initCmd)
	initGLobally = initCmd.Flags().BoolP("global", "g", false, "Init settings file globally")
}
