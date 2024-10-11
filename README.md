# GHSTAT

**GHSTAT** is a simple command-line tool that collects metrics (forks, stars and watchers) of each repository of a given GitHub user account.

## Usage:

1. Create a [GitHub Personal Access Token](https://github.com/settings/tokens?type=beta).

2. Set the environment variable `GITHUB_TOKEN` with the token value:

```bash
export GITHUB_TOKEN=your_github_token
```

3. Run the application:

```bash
go run . <GITHUB_USERNAME>
```
