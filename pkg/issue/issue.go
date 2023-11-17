package issue

import (
	"strings"
)

type Issue struct {
	RepoID   string
	Body     string
	Comments []string
}

func (i *Issue) FilterNotChecked() string {
	var filtered string
	filterBody(i.Body, &filtered)

	for _, comment := range i.Comments {
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
