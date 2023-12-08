package env

import (
	"fmt"
	"os"
	"strings"
)

func GetGitHubInfo() (string, string, string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	ghrepo := os.Getenv("GITHUB_REPOSITORY")

	if token == "" {
		return "", "", "", fmt.Errorf("GITHUB_TOKEN is required")
	}

	if ghrepo == "" {
		return "", "", "", fmt.Errorf("GITHUB_REPOSITORY is required")
	}

	splited := strings.Split(ghrepo, "/")
	if len(splited) != 2 {
		return "", "", "", fmt.Errorf("GITHUB_REPOSITORY is invalid")
	}

	username := splited[0]
	repo := splited[1]

	return token, username, repo, nil
}
