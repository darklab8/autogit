package settings

import (
	"autogit/utils"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"gopkg.in/yaml.v3"
)

//go:embed version.txt
var Version string

type ConfigScheme struct {
	Changelog  ChangelogScheme  `yaml:"changelog"`
	Regex      RegexScheme      `yaml:"regex"`
	Validation ValidationScheme `yaml:"validation"`
}

var Config ConfigScheme

func ConfigRead() {
	file, err := ioutil.ReadFile("autogit.yml")
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}

	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

// yml package has no way to validate that there is no unknown undeclared fields
func validateSettingsScheme() {
	var config ConfigScheme
	var err error

	file, _ := ioutil.ReadFile("autogit.yml")

	// Marshal file to struct
	err = yaml.Unmarshal(file, &config)
	utils.CheckFatal(err)

	// Unmarshal struct to bytes
	m, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

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
	}
}

func init() {
	ConfigRead()
	ChangelogInit()
	RegexInit()
	validateSettingsScheme()
}
