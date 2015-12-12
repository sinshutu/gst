package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/daviddengcn/go-colortext"
	"github.com/dustin/go-humanize"
	"os"
	"sort"
)

var flagsOfList = []cli.Flag{
	cli.BoolFlag{
		Name:  "short, s",
		Usage: "shorten result for pipeline processing",
	},
}

var commandList = cli.Command{
	Name:   "list",
	Action: doList,
	Flags:  flagsOfList,
}

func doList(c *cli.Context) {
	ghqPath, err := getGhqPath()
	if err != nil {
		fmt.Println("You must setup ghq first")
		os.Exit(1)
	}

	shortExpression := c.Bool("short")

	reposChannel := searchForRepos(ghqPath)

	// Sort by time
	repos := []Repository{}
	for repo := range reposChannel {
		repos = append(repos, repo)
	}
	sort.Sort(RepositoriesByModTime{repos})

	// Listing repos
	for _, repo := range repos {
		changes, err := gitStatus(repo.Path)
		if err != nil {
			continue
		}

		if shortExpression {
			fmt.Println(repo.Path)
			continue
		}

		printlnWithColor(repo.Path, ct.Cyan)
		printlnWithColor("-- "+humanize.Time(repo.ModTime), ct.Blue)

		for _, change := range changes[:len(changes)-1] {
			staged := change[:1]
			unstaged := change[1:2]
			filename := change[3:]

			if staged == "?" {
				printWithColor(staged, ct.Red)
			} else {
				printWithColor(staged, ct.Green)
			}
			printWithColor(unstaged, ct.Red)
			fmt.Println("", filename)
		}

		fmt.Println()
	}
}