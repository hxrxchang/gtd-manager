package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	repo := os.Getenv("GITHUB_REPO")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	issues, _, err :=client.Issues.ListByRepo(ctx, username, repo, nil)
	if err != nil {
		fmt.Printf("Issues.Get returned error: %v\n", err)
		os.Exit(1)
	}
	issue := issues[0]
	fmt.Printf("Body: \n%v\n", *issue.Body)
}
