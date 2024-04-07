package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	library  string
	language string
	stars    string
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
	rootCmd.PersistentFlags().StringVarP(&stars, "stars", "s", "", "Stars condition (e.g. >=500, 10..20)")

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

	client := github.NewClient(tc)

	starsQuery := ""
	if stars != "" {
		starsQuery = " stars:" + stars
	}

	languageQuery := ""
	if language != "" {
		languageQuery = " language:" + language
	}

	query := fmt.Sprintf("\"%s\" in:file%s%s", library, starsQuery, languageQuery)
	fmt.Println(query);
	opts := &github.SearchOptions{Sort: "stars", Order: "desc"}

	allRepos, err := all(client, ctx, query, opts)
	if err != nil {
		log.Fatalf("Error fetching repositories: %v", err)
	}

	fmt.Printf("Found %v repositories using %s with language %s and stars condition %s:\n", len(allRepos), library, language, stars)
	for _, repo := range allRepos {
		fmt.Printf("- %s (%s) - %d stars\n", *repo.FullName, *repo.HTMLURL, *repo.StargazersCount)
	}
}



func all(client *github.Client, ctx context.Context, query string, opts *github.SearchOptions) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Search.Repositories(ctx, query, opts)
		if err != nil {
			if _, ok := err.(*github.RateLimitError); ok {
				log.Println("hit rate limit")
				return nil, err
			}
			return nil, err
		}

		allRepos = append(allRepos, repos.Repositories...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}