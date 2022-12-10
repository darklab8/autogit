package hook

import (
	"darkvers/utils"
	"fmt"
	"os"
)

func Run(args []string) {
	fmt.Printf("commitMsg called, with args=%v\n", args)
	fmt.Println(len(os.Args), os.Args)
	commit_msg_file := args[0]
	fmt.Printf("commit_msg_file=%s\n", commit_msg_file)

	file := utils.File{Filepath: commit_msg_file}
	lines := file.FileReadLines()

	for _, line := range lines {
		fmt.Printf(";%s\n", line)
	}
}
