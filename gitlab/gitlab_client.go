package gitlab

import (
	"fmt"
	"strings"
)

var extendedCommands = [...]string{"pull-request"}

type GitlabClient struct {
	Remote string
}

func CanHandleCommand(command string) bool {
	for _, extendedCommand := range extendedCommands {
		if command == extendedCommand {
			return true
		}
	}
	return false
}

func HandleCommand(commandArgs []string, client *GitlabClient) {
	fmt.Println("handling command", commandArgs)
	serverUrl, projectName := parseServerUrlAndProjectName(client.Remote)
	projectId := getProjectId(serverUrl, projectName)
	command := commandArgs[0]
	switch command {
	case "pull-request":
		createPullRequest(projectId, commandArgs[1:], serverUrl)
	case "merge-request":
		createPullRequest(projectId, commandArgs[1:], serverUrl)
	}
	fmt.Println(projectId)
}

func getProjectId(serverUrl, projectName string) (id int) {
	project, err := searchProjectByName(serverUrl, projectName)
	if err != nil {
		panic(err)
	}
	return project.ID
}

func parseServerUrlAndProjectName(remote string) (string, string) {
	parts := strings.Split(remote,":")
	url := strings.Split(parts[0], "@")[1]
	projectParts := strings.Split(parts[1], "/")
	projectName := projectParts[len(projectParts)-1]
	projectName = strings.Split(projectName, ".")[0]

	return url, projectName
}