package client

import (
	"log"
	"context"

	"github.com/google/go-github/v32/github"
)

func getRepository(client *github.Client, ctx context.Context, owner string, name string) (*github.Repository, error) {
	repo, _, err := client.Repositories.Get(ctx, owner, name)
	if err != nil {
		return nil, err
	}

	return repo, nil
}


// FetchRepositoriesDetails コード検索結果からユニークなリポジトリの詳細情報を取得します。
func FetchRepositoriesDetails(client *github.Client, ctx context.Context, codeResults []*github.CodeResult) ([]*github.Repository) {
	uniqueRepos := make(map[string]*github.Repository)
	var repos []*github.Repository

	for _, codeResult := range codeResults {
		repo := codeResult.GetRepository()
		if repo == nil {
			continue
		}

		// リポジトリが既にマップに存在するか確認
		if _, exists := uniqueRepos[repo.GetFullName()]; !exists {
			detailedRepo, err := getRepository(client, ctx, repo.GetOwner().GetLogin(), repo.GetName())
			if err != nil {
				log.Printf("Error getting repository details: %v", err)
				continue
			}

			uniqueRepos[repo.GetFullName()] = detailedRepo
			repos = append(repos, detailedRepo)
		}
	}

	return repos
}
