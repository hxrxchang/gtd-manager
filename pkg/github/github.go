package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHub struct {
	client *githubv4.Client
}

type Issue struct {
	RepoID   string
	Body     string
	Comments []string
}

func New(token string) *GitHub {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)
	return &GitHub{client: client}
}

func (g *GitHub) GetIssueData(username, repo string) (*Issue, error) {
	variables := map[string]interface{}{
		"name":  githubv4.String(repo),
		"owner": githubv4.String(username),
	}
	var query struct {
		Repository struct {
			ID   githubv4.String
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

	err := g.client.Query(context.Background(), &query, variables)

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
		RepoID:   string(query.Repository.ID),
		Body:     body,
		Comments: comments,
	}, nil
}

func (g *GitHub) CreateIssue(repoID, title, body string) (string, error) {
	// https://docs.github.com/en/graphql/reference/mutations#createissue
	var m struct {
		CreateIssue struct {
			Issue struct {
				Title githubv4.String
			}
		} `graphql:"createIssue(input: $input)"`
	}
	input := githubv4.CreateIssueInput{
		RepositoryID: repoID,
		Title: githubv4.String(title),
		Body: githubv4.NewString(githubv4.String(body)),
	}
	err := g.client.Mutate(context.Background(), &m, input, nil)
	if err != nil {
		return "", err
	}
	return string(m.CreateIssue.Issue.Title), nil
}
