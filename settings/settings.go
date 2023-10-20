package settings

import (
	"autogit/settings/envs"
	"autogit/settings/logus"
	"autogit/settings/types"
	_ "embed"
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
var ProjectConfigPath types.ConfigPath

var UserHomeDir string

func init() {
	dirname, err := os.UserHomeDir()
	logus.CheckFatal(err, "failed obtaining user home dir")
	UserHomeDir = dirname
	GlobSettingPath = types.ConfigPath(filepath.Join(dirname, "autogit.yml"))
	ProjectConfigPath = types.ConfigPath("autogit.yml")
}

var cachedConfigFile []byte = []byte{}

func readSettingsfile(configPath types.ConfigPath) []byte {
	// TODO You could have written less config readings across your code.
	// Caching for the purpose of rendering logging records only once
	if len(cachedConfigFile) != 0 {
		return cachedConfigFile
	}

	file, err := os.ReadFile(string(configPath))
	local_file_is_not_found := false
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			logus.Debug("not found at path repository local file with config. Fallback to global config", logus.FilePath(configPath.ToFilePath()))
			local_file_is_not_found = true
		} else {
			logus.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath))
		}
	}

	global_file_is_not_found := false
	if local_file_is_not_found {
		// TODO replace with structured logging
		// fmt.Println("fallback to global settings file")
		file, err = os.ReadFile(string(GlobSettingPath))
		if err != nil {
			if strings.Contains(err.Error(), "no such file") {
				logus.Debug("not found at path repository global file with config. Fallback to other in memory config", logus.FilePath(configPath.ToFilePath()))
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

	cachedConfigFile = file
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
	var settingsPath types.ConfigPath
	workdir, _ := os.Getwd()
	project_folder := envs.TestProjectFolder
	if project_folder != "" {
		logus.Debug("OK TestProjectFolder is not empty, changing search settings to ", logus.ProjectFolder(project_folder))
		project_folder = types.ProjectFolder(workdir)
		settingsPath = types.ConfigPath(filepath.Join(string(project_folder), string(ProjectConfigPath)))
	} else {
		settingsPath = types.ConfigPath(string(ProjectConfigPath))
	}
	return settingsPath
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
