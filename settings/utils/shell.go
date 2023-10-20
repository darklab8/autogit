package utils

import (
	"autogit/settings/logus"
	"fmt"
	"os"
	"os/exec"
)

func ShellRunArgs(program string, args ...string) {
	// we are allowed breaking logging rules for code not related to application.
	logus.Debug(fmt.Sprintf("attempting to run: %s %s\n", program, args))
	executable, _ := exec.LookPath(program)

	args = append([]string{""}, args...)
	command := exec.Cmd{
		Path:   executable,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	err := command.Run()

	logus.CheckFatal(err, "unable to run cobra command")
}