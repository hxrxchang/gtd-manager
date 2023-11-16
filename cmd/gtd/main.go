package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hxrxchang/gtd-manager/pkg/env"
	"github.com/hxrxchang/gtd-manager/pkg/github"
	"github.com/hxrxchang/gtd-manager/pkg/issue"
)

func main() {
	// step1: 環境変数から必要な情報を取得する
	token, username, repo, err := env.GetGitHubInfo()

	if err != nil {
		log.Fatal(err)
	}

	// step2: GitHub APIを叩いてissueを取得する
	gh := github.New(token)
	iss, err := gh.GetIssueData(username, repo)

	if err != nil {
		log.Fatal(err)
	}

	// step3: issueのBodyとコメントのmarkdownから未完了タスクだけを抽出する
	filtered := issue.FilterNotChecked(*iss)

	// step4: 抽出したタスクをGitHubのIssueに登録する
	currentTime := time.Now()
	today := currentTime.Format("2006-01-02")
	title, err := gh.CreateIssue(iss.RepoID, today, filtered);
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("issue title: %s is created", title)
}
