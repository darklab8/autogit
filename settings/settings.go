package settings

import (
	"autogit/semanticgit/git/gitraw"
	"autogit/settings/envs"
	"autogit/settings/logus"
	"autogit/settings/types"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

const ToolName = "autogit"

var HookFolderName = fmt.Sprintf(".%s-hooks", ToolName)

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

var GlobalConfigPath types.ConfigPath

var ProjectConfigPath types.ConfigPath

var ProjectPath types.FilePath

const ConfigFileName types.ConfigPath = "autogit.yml"

func init() {
	GlobalConfigPath = types.ConfigPath(filepath.Join(string(envs.PathUserHome), string(ConfigFileName)))

	g := gitraw.NewGitRepo()
	w, err := g.Worktree()
	logus.CheckFatal(err, "are we not in git repo folder?")
	ProjectPath = types.FilePath(w.Filesystem.Root())

	ProjectConfigPath = types.ConfigPath(filepath.Join(string(ProjectPath), string(ConfigFileName)))
}

var cachedConfigFile []byte = []byte{}

func is_file_missing(err error) bool {
	_, ok := err.(*fs.PathError)
	return ok
}

func readSettingsfile(configPath types.ConfigPath) []byte {
	// TODO You could have written less config readings across your code.
	// Caching for the purpose of rendering logging records only once
	if len(cachedConfigFile) != 0 {
		return cachedConfigFile
	}

	file, err := os.ReadFile(string(configPath))
	local_file_is_not_found := false
	if err != nil {
		if is_file_missing(err) {
			logus.Debug("not found at path repository local file with config. Fallback to global config", logus.FilePath(configPath.ToFilePath()), logus.OptError(err))
			local_file_is_not_found = true
		} else {
			logus.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath), logus.OptError(err))
		}
	} else {
		logus.Debug("succesfuly read config from local repository project path", logus.ConfigPath(configPath))
	}

	global_file_is_not_found := false
	if local_file_is_not_found {
		file, err = os.ReadFile(string(GlobalConfigPath))
		if err != nil {
			if is_file_missing(err) {
				logus.Debug("not found at path repository global file with config. Fallback to other in memory config", logus.FilePath(configPath.ToFilePath()), logus.OptError(err))
				global_file_is_not_found = true
			} else {
				logus.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath), logus.OptError(err))
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

	// Config overrides for dev env purposes
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_SSH_PATH"); ok {
		result.Git.SSHPath = value
	}
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_CHANGELOG_COMMIT_URL"); ok {
		result.Changelog.CommitURL = value
	}
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_CHANGELOG_COMMIT_RANGE_URL"); ok {
		result.Changelog.CommitRangeURL = value
	}
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_CHANGELOG_ISSUE_URL"); ok {
		result.Changelog.IssueURL = value
	}

	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_VALIDATION_RULES_HEADER_SUBJECT_MIN_WORDS"); ok {
		res, err := strconv.Atoi(value)
		logus.CheckFatal(err, "crashed when trying to atoi min words env value")
		result.Validation.Rules.Header.Subject.MinWords = res
	}

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
	validateSettingsScheme(configPath)

	return config
}

var config *ConfigScheme

func GetConfig() ConfigScheme {
	if config == nil {
		config = LoadSettings(ProjectConfigPath)
	}
	return *config
}

//go:embed autogit.example.yml
var ConfigExample string
