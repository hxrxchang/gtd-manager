package env

import (
	"fmt"
	"os"
	"strings"
)

func GetGitHubInfo() (string, string, string, error) {
	var err error
	token := os.Getenv("ACCESS_TOKEN")
	ghrepo := os.Getenv("REPOSITORY")

	if token == "" {
		err = fmt.Errorf("ACCESS_TOKEN is required")
		return "", "", "", err
	}

	if ghrepo == "" {
		err = fmt.Errorf("REPOSITORY is required")
		return "", "", "", err
	}

	username := strings.Split(ghrepo, "/")[0]
	repo := strings.Split(ghrepo, "/")[1]

	return token, username, repo, nil
}
