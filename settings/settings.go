package settings

import (
	"autogit/settings/logus"
	"autogit/settings/types"
	_ "embed"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

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

var GlobSettingPath types.ConfigPath
var RelativeConfigPath types.ConfigPath

var UserHomeDir string

func init() {
	dirname, err := os.UserHomeDir()
	logus.CheckFatal(err, "failed obtaining user home dir")
	UserHomeDir = dirname
	GlobSettingPath = types.ConfigPath(filepath.Join(dirname, "autogit.yml"))
	RelativeConfigPath = types.ConfigPath("autogit.yml")
}

func readSettingsfile(configPath types.ConfigPath) []byte {
	file, err := ioutil.ReadFile(string(configPath))
	local_file_is_not_found := false
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			local_file_is_not_found = true
		} else {
			logus.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath))
		}
	}

	global_file_is_not_found := false
	if local_file_is_not_found {
		// TODO replace with structured logging
		// fmt.Println("fallback to global settings file")
		file, err = ioutil.ReadFile(string(GlobSettingPath))
		if err != nil {
			if strings.Contains(err.Error(), "no such file") {
				global_file_is_not_found = true
			} else {
				logus.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath))
			}
		}
	}

	if local_file_is_not_found && global_file_is_not_found {
		logus.Debug("fallback to memory settings file")
		file = []byte(ConfigExample)
	}

	return file
}

func ConfigRead(configPath types.ConfigPath) *ConfigScheme {
	file := readSettingsfile(configPath)

	result := ConfigScheme{}

	err := yaml.Unmarshal(file, &result)
	logus.CheckFatal(err, "unable to unmarshal settings")
	return &result
}

// yml package has no way to validate that there is no unknown undeclared fields
func validateSettingsScheme(configPath types.ConfigPath) {
	var config ConfigScheme
	var err error

	file := readSettingsfile(configPath)
	// Marshal file to struct
	err = yaml.Unmarshal(file, &config)
	logus.CheckFatal(err, "failed to unmarshal config")

	// Unmarshal struct to bytes
	m, err := yaml.Marshal(&config)
	logus.CheckFatal(err, "unable to marshal settings")

	// Marshal bytes to map
	a := make(map[interface{}]interface{})
	err = yaml.Unmarshal(m, &a)
	logus.CheckFatal(err, "failed unmarshaling to yaml")

	// compare with file marshaled to map
	b := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &b)
	logus.CheckFatal(err, "failed unmarshaling to yaml again")

	if !reflect.DeepEqual(a, b) {
		logus.Fatal(`
		setting file contains not registered keys.
		Check your version of autogit, and documentation related to settings
		`, logus.Expected(a), logus.Actual(b))
	}
}

func LoadSettings(configPath types.ConfigPath) *ConfigScheme {
	config := ConfigRead(configPath)
	ChangelogInit(*config)
	RegexInit(config)
	ValidationInit(config)
	validateSettingsScheme(configPath)

	return config
}

func GetSettingsPath() types.ConfigPath {
	workdir, _ := os.Getwd()
	project_folder := os.Getenv("AUTOGIT_PROJECT_FOLDER")
	if project_folder != "" {
		log.Println("OK AUTOGIT_PROJECT_FOLDER is not empty, changing search settings to ", project_folder)
		workdir = project_folder
	}
	settingsPath := filepath.Join(workdir, "autogit.yml")
	return types.ConfigPath(settingsPath)
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
