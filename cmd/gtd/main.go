package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v53/github"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")
	repo := os.Getenv("GITHUB_REPO")

	if token == "" {
		fmt.Println("GITHUB_ACCESS_TOKEN is required")
		os.Exit(1)
	}

	if username == "" {
		fmt.Println("GITHUB_USERNAME is required")
		os.Exit(1)
	}

	if repo == "" {
		fmt.Println("GITHUB_REPO is required")
		os.Exit(1)
	}

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

	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	var buf bytes.Buffer
	if err := md.Convert([]byte(*issue.Body), &buf); err != nil {
		fmt.Printf("md.Convert returned error: %v\n", err)
		os.Exit(1)
	  }
	fmt.Println(buf.String())
}
