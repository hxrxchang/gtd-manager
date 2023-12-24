package github

import (
	"context"
	"fmt"

	"github.com/hxrxchang/gtd-manager/pkg/issue"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHub struct {
	client  *githubv4.Client
	options *Options
}

type Options struct {
	LabelID string
}

type OptionsInput struct {
	Label string
}

func New(token string) *GitHub {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)
	return &GitHub{client: client, options: &Options{}}
}

func (g *GitHub) GetIssue(username, repo string, options *OptionsInput) (*issue.Issue, error) {
	variables := map[string]interface{}{
		"name":   githubv4.String(repo),
		"owner":  githubv4.String(username),
		"labels": (*[]githubv4.String)(nil),
	}
	if options.Label != "" {
		variables["labels"] = []githubv4.String{githubv4.String(options.Label)}
	}
	var query struct {
		Repository struct {
			ID     githubv4.String
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
			} `graphql:"issues(first: 1, orderBy: {field: CREATED_AT, direction: DESC}, labels: $labels)"`
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

	if options.Label != "" {
		var variables2 = map[string]interface{}{
			"name":      githubv4.String(repo),
			"owner":     githubv4.String(username),
			"labelName": githubv4.String(options.Label),
		}
		var query2 struct {
			Repository struct {
				Labels struct {
					Edges []struct {
						Node struct {
							ID   githubv4.String
							Name githubv4.String
						}
					}
				} `graphql:"labels(first: 1, query: $labelName)"`
			} `graphql:"repository(owner: $owner, name: $name)"`
		}

		err = g.client.Query(context.Background(), &query2, variables2)

		if err != nil {
			return nil, fmt.Errorf("failed to get label data: %w", err)
		}

		if len(query2.Repository.Labels.Edges) > 0 {
			g.options = &Options{LabelID: string(query2.Repository.Labels.Edges[0].Node.ID)}
		}
	}

	return &issue.Issue{
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
		Title:        githubv4.String(title),
		Body:         githubv4.NewString(githubv4.String(body)),
		LabelIDs:     (*[]githubv4.ID)(nil),
	}
	if g.options.LabelID != "" {
		input.LabelIDs = &[]githubv4.ID{githubv4.ID(g.options.LabelID)}
	}
	err := g.client.Mutate(context.Background(), &m, input, nil)
	if err != nil {
		return "", err
	}
	return string(m.CreateIssue.Issue.Title), nil
}
