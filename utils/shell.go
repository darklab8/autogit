package utils

import (
	"log"
	"os"
	"os/exec"
)

func ShellRunArgs(program string, args ...string) {
	executable, _ := exec.LookPath(program)

	args = append([]string{""}, args...)
	command := exec.Cmd{
		Path:   executable,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	err := command.Run()

	if err != nil {
		log.Fatal(err)
	}
}
