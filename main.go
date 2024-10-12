package main

import (
	"log"
	"os"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GitHub token is not set. Please set the GITHUB_TOKEN environment variable.")
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <GITHUB_USERNAME>")
	}

	username := os.Args[1]

	fetcher := newGitHubFetcher(token)

	var formatter tableFormatter

	if err := fetcher.fetchAndFormat(username, &formatter); err != nil {
		log.Fatalf("Error fetching and formatting: %s", err)
	}
}
