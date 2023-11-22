package issue

import (
	"slices"
	"strings"
)

type Issue struct {
	RepoID   string
	Body     string
	Comments []string
}

type DividedTasks map[string][]string

func (i *Issue) Process() string {
	tasks := make(DividedTasks)
	keysOrder := []string{""}

	mergeAndFilter(i.Body, &tasks, &keysOrder)
	for _, comment := range i.Comments {
		mergeAndFilter(comment, &tasks, &keysOrder)
	}

	var res string
	for _, key := range keysOrder {
		res += key + "\n"
		for _, task := range tasks[key] {
			res += task + "\n"
		}
	}

	return res
}

func mergeAndFilter(body string, tasks *DividedTasks, keyOrder *[]string) {
	headingPrefixes := []string{"#", "##", "###", "####", "#####", "######"}
	tmpKey := ""
	for _, line := range strings.Split(body, "\n") {
		spaceTrimed := strings.TrimLeft(line, " ")
		prefix := strings.Split(spaceTrimed, " ")[0]

		if slices.Contains(headingPrefixes, prefix) {
			tmpKey = spaceTrimed
			if !slices.Contains(*keyOrder, tmpKey) {
				*keyOrder = append(*keyOrder, tmpKey)
			}
		}

		if strings.HasPrefix(spaceTrimed, "- [ ]") {
			(*tasks)[tmpKey] = append((*tasks)[tmpKey], line)
		}
	}
}


