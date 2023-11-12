package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	// step1: 環境変数から必要な情報を取得する
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

	// step2: GitHub APIを叩いてissueを取得する
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)

	variables := map[string]interface{}{
		"name":  githubv4.String(repo),
		"owner": githubv4.String(username),
	}

	var query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
		Repository struct {
			Issues struct {
				Edges []struct {
					Node struct {
						Title githubv4.String
						Body  githubv4.String
					}
				}
			} `graphql:"issues(first: 1, states: OPEN, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err)
		// Handle error.
	}
	fmt.Println("    Login:", query.Viewer.Login)
	fmt.Println("CreatedAt:", query.Viewer.CreatedAt)
	fmt.Println("Title", query.Repository.Issues.Edges[0].Node.Title)
	fmt.Println("Body", query.Repository.Issues.Edges[0].Node.Body)
}
