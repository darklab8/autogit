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
		config_name := "autogit.yml"
		if utils.FileExists(config_name) {
			log.Fatalln("file with settings=", config_name, " already exists")
			return
		}

		file := utils.NewFile(config_name).CreateToWriteF()
		defer file.Close()
		file.WritelnF(settings.ConfigExample)

		fmt.Println("Succesfully created autogit.yml in current location")
		fmt.Println("Try to " + initAdvice + ". It will automatically verify committs for you!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
