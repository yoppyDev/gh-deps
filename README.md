# gh-deps

A command-line tool to search for repositories on GitHub based on their dependencies.

## Getting Started

Before using the GitHub Dependency Searcher CLI, you need to set up your GitHub personal access token. This token allows the tool to access GitHub's API and perform searches on your behalf. Follow the steps below to get started:

1. Generate a GitHub Personal Access Token:

- Go to your GitHub settings.
- Navigate to the "Developer settings" section.
- Click on "Personal access tokens".
- Generate a new token with the permissions required to access repositories and read their dependencies.

2. Set up your environment:
Export your GitHub token as an environment variable in your terminal. Replace ghp_XXXXXXXX with your actual token.

```
export GITHUB_TOKEN=ghp_XXXXXXXX
```

This command sets the GITHUB_TOKEN environment variable, which the CLI tool uses to authenticate with the GitHub API.

Now you are ready to use the GitHub Dependency Searcher CLI to find repositories based on specific dependencies.🎉
