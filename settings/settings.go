package settings

import (
	"autogit/utils"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v3"
)

const ToolName = "autogit"

//go:embed version.txt
var Version string

type ConfigScheme struct {
	Changelog  ChangelogScheme  `yaml:"changelog"`
	Regex      RegexScheme      `yaml:"regex"`
	Validation ValidationScheme `yaml:"validation"`
	Git        struct {
		SSHPath string `yaml:"sshPath"`
	} `yaml:"git"`
}

type SettingPath string

func ConfigRead(settingsPath SettingPath) *ConfigScheme {

	file, err := ioutil.ReadFile(string(settingsPath))
	utils.CheckFatal(err, "Could not read the file due to error, autogit_path=%s\n", string(settingsPath))

	result := ConfigScheme{}

	err = yaml.Unmarshal(file, &result)
	utils.CheckFatal(err, "unable to unmarshal settings")
	return &result
}

// yml package has no way to validate that there is no unknown undeclared fields
func validateSettingsScheme(settingsPath SettingPath) {
	var config ConfigScheme
	var err error

	file, err := ioutil.ReadFile(string(settingsPath))
	utils.CheckFatal(err, "Could not read the file due to error, autogit_path=%s\n", string(settingsPath))

	// Marshal file to struct
	err = yaml.Unmarshal(file, &config)
	utils.CheckFatal(err)

	// Unmarshal struct to bytes
	m, err := yaml.Marshal(&config)
	utils.CheckFatal(err, "unable to unmarshal settings")

	// Marshal bytes to map
	utils.CheckFatal(err)
	a := make(map[interface{}]interface{})
	err = yaml.Unmarshal(m, &a)
	utils.CheckFatal(err)

	// compare with file marshaled to map
	b := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &b)
	utils.CheckFatal(err)

	if !reflect.DeepEqual(a, b) {
		fmt.Printf("ERR autogit.yml contains not registered keys. Check your version of autogit, and documentation related to settings\n")
		fmt.Printf("--- expected:\n%v\n\n", a)
		fmt.Printf("--- actual:\n%v\n\n", b)
		os.Exit(1)
	}
}

func LoadSettings(settingsPath SettingPath) *ConfigScheme {
	config := ConfigRead(settingsPath)
	ChangelogInit(*config)
	RegexInit(config)
	ValidationInit(config)
	validateSettingsScheme(settingsPath)

	return config
}

func GetSettingsPath() SettingPath {
	workdir, _ := os.Getwd()
	project_folder := os.Getenv("AUTOGIT_PROJECT_FOLDER")
	if project_folder != "" {
		log.Println("OK AUTOGIT_PROJECT_FOLDER is not empty, changing search settings to ", project_folder)
		workdir = project_folder
	}
	settingsPath := filepath.Join(workdir, "autogit.yml")
	return SettingPath(settingsPath)
}

var config *ConfigScheme

func GetConfig() ConfigScheme {
	if config == nil {
		settingPath := GetSettingsPath()
		config = LoadSettings(settingPath)
	}
	return *config
}

//go:embed autogit.example.yml
var ConfigExample string
