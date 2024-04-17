# gh-deps

A command-line tool to search for repositories on GitHub based on their dependencies.


# Installation

```bash
brew tap yoppyDev/gh-deps
brew install gh-deps
```

# Usage

```bash
export GITHUB_TOKEN=<your_github_token>

gh-deps -l <library> -p <path>

# Example
gh-deps -l "spf13/cobra" -p "**/**go.mod"
```

# Parameters
| Parameter | Description | Required |
| --- | --- | --- |
| -l, --library | The library to search for | Yes |
| -p, --path | The path to search for | Yes |