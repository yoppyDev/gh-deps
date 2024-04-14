package client

import (
	"sort"

	"github.com/google/go-github/v32/github"
)

func SortByStar(repos []*github.Repository) {
	if len(repos) == 0 {
		return
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].GetStargazersCount() > repos[j].GetStargazersCount()
	})
}