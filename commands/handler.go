package commands

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const basicOsCommand = "git"

// Handle command from user. Do own command or provide it to system git
func Handle(commandArgs []string) {
	handleWithSystemGit(commandArgs)
}

func handleWithSystemGit(commandArgs []string) {
	cmd, lookErr := exec.LookPath("git", )

	if lookErr != nil {
		panic(lookErr)
	}

	env := os.Environ()

	args := make([]string, 1)
	args[0] = basicOsCommand

	args = append(args, commandArgs...)
	fmt.Println(args)
	execErr := syscall.Exec(cmd, args, env)
	if execErr != nil {
		panic(execErr)
	}
}