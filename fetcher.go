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

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not do request for fetching repository stat: %w", err)
	}

	req.Header.Set("Authorization", "token "+g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do request for fetching repository stat: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected status code %d", resp.StatusCode)
	}

	var repos []repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("could not decode body payload from fetched repository stat: %w", err)
	}
	return repos, nil
}

func (g *githubFetcher) fetchAndFormat(username string, formatter outputFormatter) error {
	repoChan := make(chan []repository)
	errChan := make(chan error)

	go func() {
		g.fetchStat(repoChan, errChan, username)
		close(repoChan)
		close(errChan)
	}()

	for {
		select {
		case repos, ok := <-repoChan:
			if !ok {
				repoChan = nil
			} else {
				formatter.format(repos)
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				log.Println("Error received on stat fetching repository:", err)
			}
		}

		// Break out of the loop when both channels are closed
		if repoChan == nil && errChan == nil {
			break
		}
	}
	return nil
}

func (g *githubFetcher) fetchStat(repoCh chan []repository, errCh chan error, username string) {
	page, perPage := 1, 100

	for {
		repos, err := g.fetch(username, page, perPage)
		if err != nil {
			errCh <- fmt.Errorf("could not fetch data for username: '%s', page: '%d': %w", username, page, err)
			return
		}

		if len(repos) == 0 {
			return // Stop when no more repositories are returned
		}
		repoCh <- repos
		page++
	}
}
