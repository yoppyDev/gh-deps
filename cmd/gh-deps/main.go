package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yoppyDev/gh-deps/client"
	"golang.org/x/oauth2"
)

var (
	library  string
	path string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:           "depSearch",
		Short:         "Search for repositories that use a specific library",
		Run:           run,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.Flags().StringVarP(&library, "library", "l", "", "Library to search for (e.g. spf13/cobra)")
	rootCmd.Flags().StringVarP(&path, "path", "p", "", "path to search for (e.g. **/**go.mod)")

	if err := rootCmd.MarkFlagRequired("library"); err != nil {
		fmt.Printf("Error marking library as required: %v\n", err)
		os.Exit(1)
	}
	if err := rootCmd.MarkFlagRequired("path"); err != nil {
		fmt.Printf("Error marking path as required: %v\n", err)
		os.Exit(1)
	}

	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	token := viper.GetString("GITHUB_TOKEN")
	if token == "" {
		fmt.Printf("GITHUB_TOKEN is not set in the environment variables")
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	c := github.NewClient(tc)

	query := fmt.Sprintf("\"%s\" %s", library, path)
	result, _, err := c.Search.Code(ctx, query, &github.SearchOptions{})
	if err != nil {
		fmt.Printf("Error searching code: %v\n", err)
		return
	}

	repos := client.FetchRepositoriesDetails(c, ctx, result.CodeResults)

	// スター数に基づいてリポジトリをソート
	client.SortByStar(repos)

	// ソートされたリポジトリを出力
	for _, repo := range repos {
		fmt.Printf("- %s (%s) - %d stars\n", *repo.FullName, *repo.HTMLURL, *repo.StargazersCount)
	}
}
