package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/yoppyDev/github-dependency-searcher/client"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	library  string
	language string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:           "depSearch",
		Short:         "Search for repositories that use a specific library",
		Run:           run,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.PersistentFlags().StringVarP(&library, "library", "l", "", "Library to search for (e.g. spf13/cobra)")
	rootCmd.PersistentFlags().StringVarP(&language, "language", "L", "", "Language to search for (e.g. go)")

	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	token := viper.GetString("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set in the environment variables")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	c := github.NewClient(tc)

	languageQuery := ""
	if language != "" {
		languageQuery = " language:" + language
	}

	query := fmt.Sprintf("\"%s\" %s", library, languageQuery)
	result, _, err := c.Search.Code(ctx, query, &github.SearchOptions{})
	if err != nil {
		log.Fatalf("Error searching code: %v", err)
	}

	repos := client.FetchRepositoriesDetails(c, ctx, result.CodeResults)

	// スター数に基づいてリポジトリをソート
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].GetStargazersCount() > repos[j].GetStargazersCount()
	})

	// ソートされたリポジトリを出力
	for _, repo := range repos {
		fmt.Printf("- %s (%s) - %d stars\n", repo.GetFullName(), repo.GetHTMLURL(), repo.GetStargazersCount())
	}
}