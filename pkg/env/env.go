package env

import (
	"fmt"
	"os"
	"strings"
)

func GetGitHubInfo() (string, string, string, error) {
	var err error
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	ghrepo := os.Getenv("GITHUB_REPOSITORY")

	if token == "" {
		err = fmt.Errorf("GITHUB_ACCESS_TOKEN is required")
	}

	username := strings.Split(ghrepo, "/")[0]
	repo := strings.Split(ghrepo, "/")[1]

	return token, username, repo, err
}
