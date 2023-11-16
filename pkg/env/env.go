package env

import (
	"fmt"
	"os"
)

func GetGitHubInfo() (string, string, string, error) {
	var err error
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	repo := os.Getenv("GITHUB_REPO")

	if token == "" {
		err = fmt.Errorf("GITHUB_ACCESS_TOKEN is required")
	}

	if username == "" {
		err = fmt.Errorf("GITHUB_USERNAME is required")
	}

	if repo == "" {
		err = fmt.Errorf("GITHUB_REPO is required")
	}
	return token, username, repo, err
}
