package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// repository represents a GitHub repository.
type repository struct {
	Name            string `json:"name"`
	ForksCount      int    `json:"forks_count"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	Fork            bool   `json:"fork"`
}

// githubFetcher is a struct that implements repositoryFetcher for GitHub.
type githubFetcher struct {
	token  string
	client *http.Client
}

// newGitHubFetcher creates a new instance of githubFetcher.
func newGitHubFetcher(token string) *githubFetcher {
	return &githubFetcher{
		token: token,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// fetch fetches repositories for the given username and pagination settings.
func (g *githubFetcher) fetch(username string, page int, perPage int) ([]repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=%d&page=%d", username, perPage, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	var repos []repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// fetchAndFormat fetches repositories and formats the output.
func (g *githubFetcher) fetchAndFormat(username string, formatter outputFormatter) error {
	repoChan := make(chan []repository)

	go func() {
		defer close(repoChan)
		page := 1
		perPage := 100
		for {
			repos, err := g.fetch(username, page, perPage)
			if err != nil {
				log.Println(err)
				break
			}
			if len(repos) == 0 {
				break
			}
			repoChan <- repos
			page++
		}
	}()

	// Collect and format repositories
	for repos := range repoChan {
		formatter.format(repos)
	}
	return nil
}
