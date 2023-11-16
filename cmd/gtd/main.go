package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	issue, err := getIssueData(token, username, repo)

	if err != nil {
		log.Fatal(err)
	}

	// step3: issueのBodyとコメントのmarkdownから未完了タスクだけを抽出する
	var filtered string
	fmt.Println(filtered)
	filterNotChecked(issue.Body, &filtered)

	for _, comment := range issue.Comments {
		filterNotChecked(comment, &filtered)
	}
	fmt.Println(filtered)
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

type Issue struct {
	Body     string
	Comments []string
}

func getIssueData(token, username, repo string) (*Issue, error) {
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
						Title    githubv4.String
						Body     githubv4.String
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
		return nil, fmt.Errorf("failed to get issue data: %w", err)
	}

	if len(query.Repository.Issues.Edges) == 0 {
		return nil, fmt.Errorf("issue not found")
	}

	body := string(query.Repository.Issues.Edges[0].Node.Body)
	var comments []string
	for _, c := range query.Repository.Issues.Edges[0].Node.Comments.Edges {
		comments = append(comments, string(c.Node.Body))
	}
	return &Issue{
		Body:     body,
		Comments: comments,
	}, nil
}

func splitByLine(s string) []string {
	return strings.Split(s, "\n")
}

func filterNotChecked(body string, res *string) {
	for _, line := range splitByLine(body) {
		speceTrimmed := strings.TrimLeft(line, " ")
		if strings.HasPrefix(speceTrimmed, "- [ ]") {
			*res += line + "\n"
		}
	}
}
