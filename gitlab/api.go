package gitlab

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Namespace struct {
	FullPath string `json:"full_path"`
	ID       int    `json:"id"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	Path     string `json:"path"`
}

type RemoteProject struct {
	AvatarUrl         string `json:"avatar_url"`
	CreatedAt         string `json:"created_at"`
	Description       string `json:"description"`
	ForksCount        int    `json:"forks_count"`
	HTTPUrlToRepo     string `json:"http_url_to_repo"`
	ID                int    `json:"id"`
	LastActivityAt    string `json:"last_activity_at"`
	Name              string `json:"name"`
	NameWithNamespace string `json:"name_with_namespace"`
	Namespace         *Namespace
	Path              string   `json:"path"`
	PathWithNamespace string   `json:"path_with_namespace"`
	ReadmeUrl         string   `json:"readme_url"`
	SSHUrlToRepo      string   `json:"ssh_url_to_repo"`
	StarCount         int      `json:"star_count"`
	TagList           []string `json:"tag_list"`
	WebUrl            string   `json:"web_url"`
}

func searchProjectByName(serverUrl, projectName string) (*RemoteProject, error) {
	req, err := http.NewRequest("GET", buildUrl(serverUrl)+"search", nil)
	if err != nil {
		panic(err)
	}
	apiToken := os.Getenv("GITLAB_API_TOKEN")
	req.Header.Add("PRIVATE-TOKEN", apiToken)
	q := req.URL.Query()
	q.Add("scope", "projects")
	q.Add("search", projectName)
	req.URL.RawQuery = q.Encode()

	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var searchRes []RemoteProject
		jsonErr := json.NewDecoder(resp.Body).Decode(&searchRes)
		if jsonErr != nil {
			panic(jsonErr)
		}

		project := searchRes[0]
		return &project, nil
	} else {
		var searchRes []map[string]interface{}
		jsonErr := json.NewDecoder(resp.Body).Decode(&searchRes)
		fmt.Println("server returned", searchRes)
		if jsonErr != nil {
			panic(jsonErr)
		}
		return nil, errors.New("project not found")
	}
}

func createPullRequest(projectId int, args []string, serverUrl string) {
	fmt.Println("creating pull request for", projectId, "with", args)
	fullUrl := buildUrl(serverUrl) + "projects/" + strconv.Itoa(projectId) + "/merge_requests"
	// curl --request POST --header "PRIVATE-TOKEN: _____" -F "id=60" -F
	// "source_branch=pustserg/test-mr-branch" -F "target_branch=master"
	// -F "title=test" https://gitlab.rocketbank.sexy/api/v4/projects/60/merge_requests
	// TODO: Make good struct for it
	// {
	//   "id":14276,
	//   "iid":4,
	//   "project_id":60,
	//   "title":"test",
	//   "description":null,
	//   "state":"opened",
	//   "created_at":"2019-04-08T07:40:18.027Z",
	//   "updated_at":"2019-04-08T07:40:18.027Z",
	//   "merged_by":null,
	//   "merged_at":null,
	//   "closed_by":null,
	//   "closed_at":null,
	//   "target_branch":"master",
	//   "source_branch":"pustserg/test-mr-branch",
	//   "upvotes":0,
	//   "downvotes":0,
	//   "author":{
	//     "id":43,
	//     "name":"Sergey Pustovalov",
	//     "username":"s.pustovalov",
	//     "state":"active",
	//     "avatar_url":"https://gitlab.rocketbank.sexy/uploads/-/system/user/avatar/43/avatar.png",
	//     "web_url":"https://gitlab.rocketbank.sexy/s.pustovalov"
	//    },
	//    "assignee":null,
	//    "source_project_id":60,
	//    "target_project_id":60,
	//    "labels":[],
	//    "work_in_progress":false,
	//    "milestone":null,
	//    "merge_when_pipeline_succeeds":false,
	//    "merge_status":"can_be_merged",
	//    "sha":"b09cc9d277b25bfc219bfdf8a56b35bb00bb75ee",
	//    "merge_commit_sha":null,
	//    "user_notes_count":0,
	//    "discussion_locked":null,
	//    "should_remove_source_branch":null,
	//    "force_remove_source_branch":null,
	//    "web_url":"https://gitlab.rocketbank.sexy/backend/shared/rocket-protocat/merge_requests/4",
	//    "time_stats":{"time_estimate":0,
	//    "total_time_spent":0,
	//    "human_time_estimate":null,
	//    "human_total_time_spent":null},
	//    "squash":false,
	//    "subscribed":true,
	//    "changes_count":"1",
	//    "latest_build_started_at":null,
	//    "latest_build_finished_at":null,
	//    "first_deployed_to_production_at":null,
	//    "pipeline":null,
	//    "diff_refs":{"base_sha":"1ad0a0faa4095ad943129b126073f745ec743c02",
	//    "head_sha":"b09cc9d277b25bfc219bfdf8a56b35bb00bb75ee",
	//    "start_sha":"1ad0a0faa4095ad943129b126073f745ec743c02"},
	//    "merge_error":null,
	//    "approvals_before_merge":null
	//  }

	var requestTitle string
	if len(args) > 0 {
		requestTitle = args[0]
	} else {
		fmt.Print("Enter merge request title: ")
		reader := bufio.NewReader(os.Stdin)
		requestTitle, _ = reader.ReadString('\n')
	}
	fmt.Println("requestTitle", requestTitle)
	fmt.Println(fullUrl)

	cmd:= exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	currentBranch := string(out)

	formData := url.Values{
		"id":            {string(projectId)},
		"source_branch": {currentBranch},
		"target_branch": {"master"},
		"title":         {requestTitle},
	}
	fmt.Println("formData", formData)
	fmt.Println("fullUrl", fullUrl)
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Add("PRIVATE-TOKEN", os.Getenv("GITLAB_API_TOKEN"))

	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	var serverResponse map[string]interface{}
	jsonErr := json.NewDecoder(resp.Body).Decode(&serverResponse)
	if jsonErr != nil {
		panic(jsonErr)
	}

	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 201 {
		fmt.Println("Pull request created", serverResponse["web_url"])
	} else {
		fmt.Println("server returned error", serverResponse)
	}

}

func buildUrl(serverUrl string) string {
	return "https://" + serverUrl + "/api/v4/"
}
