package env

import (
	"fmt"
	"os"
	"strings"
)

func GetGitHubInfo() (string, string, string, error) {
	token := os.Getenv("ACCESS_TOKEN")
	ghrepo := os.Getenv("REPOSITORY")

	if token == "" {
		return "", "", "", fmt.Errorf("ACCESS_TOKEN is required")
	}

	if ghrepo == "" {
		return "", "", "", fmt.Errorf("REPOSITORY is required")
	}

	username := strings.Split(ghrepo, "/")[0]
	repo := strings.Split(ghrepo, "/")[1]

	return token, username, repo, nil
}
