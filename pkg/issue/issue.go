package issue

import (
	"strings"

	"github.com/hxrxchang/gtd-manager/pkg/github"
)

func FilterNotChecked(issue github.Issue) string {
	var filtered string
	filterBody(issue.Body, &filtered)

	for _, comment := range issue.Comments {
		filterBody(comment, &filtered)
	}
	return filtered
}

func filterBody(body string, res *string) {
	for _, line := range splitByLine(body) {
		speceTrimmed := strings.TrimLeft(line, " ")
		if strings.HasPrefix(speceTrimmed, "- [ ]") {
			*res += line + "\n"
		}
	}
}

func splitByLine(s string) []string {
	return strings.Split(s, "\n")
}
