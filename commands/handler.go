package commands

import (
	"github.com/pustserg/lab/gitlab"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const basicOsCommand = "git"
var remoteServer string
var remoteGitlabClient gitlab.GitlabClient

// Handle command from user. Do own command or provide it to system git
func Handle(commandArgs []string) {
	if len(commandArgs) > 0 {
		detectRemote()
		if remoteGitlabClient.Remote != ""{
			if gitlab.CanHandleCommand(commandArgs[0]) {
				fmt.Println("Client can handle this command")
				gitlab.HandleCommand(commandArgs, &remoteGitlabClient)
			} else {
				fmt.Println("client can't handle this command")
				handleWithSystemGit(commandArgs)
			}
		}
	} else {
		handleWithSystemGit(commandArgs)
	}
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

	execErr := syscall.Exec(cmd, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func detectRemote() {
	cmd:= exec.Command("git", "remote", "-v")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	response := strings.Fields(string(out))

	for _, item:= range response {
		if strings.HasPrefix(item, "git@github") {
			remoteServer = item
		} else if strings.HasPrefix(item, "git@gitlab") {
			remoteServer = item
			remoteGitlabClient = gitlab.GitlabClient{Remote: remoteServer}
		}
	}
	if remoteServer == "" {
		fmt.Println("No remote server found")
	}
}