package actions

import (
	"autogit/parser/conventionalcommits"
	"autogit/utils"
	"fmt"
	"io/ioutil"
)

func CommmitMsg(args []string) {
	inputFile := args[0]

	file, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	// convert the file binary into a string using string
	fileContent := string(file)

	commit, err := conventionalcommits.NewCommit(fileContent)

	if err != nil {
		utils.CheckFatal(err, "unable to parse commit to conventional commits standard")
	}

	fmt.Println("parsed commit")
	fmt.Printf("type=%s\n", commit.Type)
	if commit.Scope != "" {
		fmt.Printf("scope=%s\n", commit.Scope)
	}
	fmt.Printf("subject=%s\n", commit.Subject)

	if commit.Body != "" {
		fmt.Printf("body=%s\n", commit.Body)
	}

	for index, footer := range commit.Footers {
		fmt.Printf("footer #%d - token: %s, content: %s\n", index, footer.Token, footer.Content)
	}

}
