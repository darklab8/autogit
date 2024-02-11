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
	"strings"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/goutils/utils/utils_logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"

	"gopkg.in/yaml.v3"
)

const ToolName = "autogit"

var HookFolderName = fmt.Sprintf(".%s-hooks", ToolName)

//go:embed version.txt
var autogit_version string

func GetAutogitVersion() string {
	// cleaning up version from... debugging logs used during dev env
	lines := strings.Split(autogit_version, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "v") {
			return line
		}
	}
	return autogit_version
}

type ConfigScheme struct {
	Changelog  ChangelogScheme  `yaml:"changelog"`
	Regex      RegexScheme      `yaml:"regex"`
	Validation ValidationScheme `yaml:"validation"`
	Git        struct {
		SSHPath string `yaml:"sshPath"`
	} `yaml:"git"`
}

// var GlobalConfigPath types.ConfigPath

// var ProjectConfigPath types.ConfigPath

const ConfigFileName types.ConfigPath = "autogit.yml"

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
			logus.Log.Debug("not found at path repository local file with config. Fallback to global config", utils_logus.FilePath(configPath.ToFilePath()), typelog.OptError(err))
			local_file_is_not_found = true
		} else {
			logus.Log.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath), typelog.OptError(err))
		}
	} else {
		logus.Log.Debug("succesfuly read config from local repository project path", logus.ConfigPath(configPath))
	}

	global_file_is_not_found := false
	if local_file_is_not_found {
		file, err = os.ReadFile(string(GetGlobalConfigPath()))
		if err != nil {
			if is_file_missing(err) {
				logus.Log.Debug("not found at path repository global file with config. Fallback to other in memory config", utils_logus.FilePath(configPath.ToFilePath()), typelog.OptError(err))
				global_file_is_not_found = true
			} else {
				logus.Log.CheckFatal(err, "Could not read the file due to error", logus.ConfigPath(configPath), typelog.OptError(err))
			}
		}
	}

	if local_file_is_not_found && global_file_is_not_found {
		logus.Log.Debug("fallback to memory settings file")
		file = []byte(ConfigExample)
	}

	cachedConfigFile = file
	return file
}

func merge_maps[T comparable](result map[T]interface{}, additions map[T]interface{}) map[T]interface{} {
	/*
		for key value in result {
			if key is present in additions {

				if value is not map[string]interface{} {
					result[key] = additions[key]
				} else {
					merge_maps(result[key], additions[key])
				}

			}
		}
	*/

	// For all keys and values in resulting hashmap
	for result_key, result_value := range result {

		// if key is present in additions hasmap
		if additions_value, is_present := additions[result_key]; is_present {

			// if value is not map[string]{interface}
			if asserted_addition_value, ok := additions_value.(map[string]interface{}); !ok {
				// just override with value from additions to resulting hashmap
				result[result_key] = additions_value
			} else {

				// otherwise try to merge recursively values from additions nested hashmaps to resulting key
				if asserted_result_value, ok := result_value.(map[string]interface{}); ok {
					merge_maps(asserted_result_value, asserted_addition_value)
				} else {
					logus.Log.Fatal(`
						failed to assert value of config in memory being of same type as value in input config
						potentially broken config
					`)
				}
			}
		}
	}

	return result
}

func configRead(file []byte) *ConfigScheme {
	config := ConfigScheme{}

	file_config := make(map[interface{}]interface{})
	err := yaml.Unmarshal(file, &file_config)
	logus.Log.CheckFatal(err, "unable to unmarshal input config")

	memory_config := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(ConfigExample), &memory_config)
	logus.Log.CheckFatal(err, "unable to unmrashal memory config")

	// merged file content onto memory config
	merged_config := merge_maps(memory_config, file_config)

	merged_config_as_bytes, err := yaml.Marshal(&merged_config)
	logus.Log.CheckFatal(err, "unable to marshal merged config")

	err = yaml.Unmarshal(merged_config_as_bytes, &config)
	logus.Log.CheckError(err, `unable to unmarshal merged config.
	Your autogit.yml settings file is having invalue key: value pairs.
	Try to remove previous autogit.yml settings files and generate a new one with "autogit init [--global] command"`)

	return &config
}

func (config *ConfigScheme) configLoadEnvOverrides() {
	// Config overrides for dev env purposes
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_SSH_PATH"); ok {
		config.Git.SSHPath = value
	}
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_CHANGELOG_REPOSITORY_OWNER"); ok {
		config.Changelog.REPOSITORY_OWNER = value
	}
	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_CHANGELOG_REPOSITORY_NAME"); ok {
		config.Changelog.REPOSITORY_NAME = value
	}

	if value, ok := os.LookupEnv("AUTOGIT_CONFIG_VALIDATION_RULES_HEADER_SUBJECT_MIN_WORDS"); ok {
		res, err := strconv.Atoi(value)
		logus.Log.CheckFatal(err, "crashed when trying to atoi min words env value")
		config.Validation.Rules.Header.Subject.MinWords = res
	}
}

func check_file_is_not_having_invalid[T comparable](example map[T]interface{}, checkable map[T]interface{}) {
	// For all keys and values in resulting hashmap
	for checkable_key, checkable_value := range checkable {

		// if key is present in additions hasmap
		if example_value, is_present := example[checkable_key]; !is_present {
			logus.Log.Error(fmt.Sprintf(`autogit.yml file is having not allowed user settings key
			Please, remove previous autogit.yml local in your repository and user global one and generate new one
			with "autogit init [--global]" command :). Or remove the specified "settings_key"`), logus.SettingsKey(checkable_key))
		} else {
			if reflect.TypeOf(example_value) != reflect.TypeOf(checkable_value) {
				logus.Log.Fatal(fmt.Sprintf(
					"wrong value type in config. Expected: %v, Received:%v",
					reflect.TypeOf(example_value),
					reflect.TypeOf(checkable_value),
				))
			}
		}

		if nested_checkable_map, ok := checkable_value.(map[string]interface{}); ok {
			if nested_example_map, ok := example[checkable_key].(map[string]interface{}); ok {
				check_file_is_not_having_invalid(nested_example_map, nested_checkable_map)
			}
		}
	}
}

func validate_file_config(file []byte) {

	file_config := make(map[interface{}]interface{})
	err := yaml.Unmarshal(file, &file_config)
	logus.Log.CheckFatal(err, "unable to unmarshal input config")

	memory_config := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(ConfigExample), &memory_config)
	logus.Log.CheckFatal(err, "unable to unmrashal memory config")

	check_file_is_not_having_invalid(memory_config, file_config)
}

func NewConfig(configPath types.ConfigPath) *ConfigScheme {

	file := readSettingsfile(configPath)

	config := configRead(file)
	config.configLoadEnvOverrides()
	config.changelogValidate()
	config.regexCompile()

	validate_file_config(file)

	return config
}

var config *ConfigScheme

func GetGlobalConfigPath() types.ConfigPath {
	return types.ConfigPath(filepath.Join(string(envs.PathUserHome), string(ConfigFileName)))
}

func GetProjectConfigPath() types.ConfigPath {
	g := gitraw.NewGitRepo()
	w, err := g.Worktree()
	logus.Log.CheckFatal(err, "are we not in git repo folder?")
	project_path := utils_types.FilePath(w.Filesystem.Root())
	return types.ConfigPath(utils_types.FilePath(filepath.Join(string(project_path), string(ConfigFileName))))
}

func GetConfig() ConfigScheme {
	if config == nil {
		config = NewConfig(GetProjectConfigPath())
	}
	return *config
}

//go:embed autogit.example.yml
var ConfigExample string
