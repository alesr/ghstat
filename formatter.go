package main

import (
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// outputFormatter is an interface for formatting output.
type outputFormatter interface {
	format(repos []repository)
}

// tableFormatter formats the output as a table.
type tableFormatter struct{}

// format formats the repositories into a table.
func (tf *tableFormatter) format(repos []repository) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repository", "Forks", "Stars", "Watchers"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})

	for _, repo := range repos {
		if repo.Fork {
			continue // Skip forked repositories
		}

		data := []string{
			repo.Name,
			highlightZero(repo.ForksCount),
			highlightZero(repo.StargazersCount),
			highlightZero(repo.WatchersCount),
		}

		table.Append(data)
	}
	table.Render()
}

// highlightZero highlights values greater than 0 in green.
func highlightZero(value int) string {
	if value > 0 {
		return color.New(color.FgGreen).SprintFunc()(strconv.Itoa(value))
	}
	return strconv.Itoa(value)
}
