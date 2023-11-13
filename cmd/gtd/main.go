package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	// step1: 環境変数から必要な情報を取得する
	token, username, repo, err := getGitHubInfo()

	if err != nil {
		log.Fatal(err)
	}

	// step2: GitHub APIを叩いてissueを取得する
	getIssueData(token, username, repo)
}

func getGitHubInfo() (string, string, string, error) {
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

func getIssueData(token, username, repo string) {
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
		Repository struct {
			Issues struct {
				Edges []struct {
					Node struct {
						Title githubv4.String
						Body  githubv4.String
						Comments struct {
							Edges []struct {
								Node struct {
									Body githubv4.String
								}
							}
						} `graphql:"comments(first: 100)"`
					}
				}
			} `graphql:"issues(first: 1, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Title", query.Repository.Issues.Edges[0].Node.Title)
	fmt.Println("Body", query.Repository.Issues.Edges[0].Node.Body)
	fmt.Println("Comment", query.Repository.Issues.Edges[0].Node.Comments)
}
