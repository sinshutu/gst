package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/daviddengcn/go-colortext"
	"strings"
)

var flagsOfUpdate = []cli.Flag{
	cli.BoolFlag{
		Name:  "short, s",
		Usage: "shorten result for pipeline processing",
	},
}

var commandUpdate = cli.Command{
	Name:   "update",
	Action: doUpdate,
	Flags:  flagsOfUpdate,
}

type GitPullResponse struct {
	Path  string
	Body  string
	Error error
}

func doUpdate(c *cli.Context) error {
	ghqPath := verifyGhqPath()
	repos := searchForRepos(ghqPath)

	responses := []GitPullResponse{}

	for repo := range repos {
		printlnWithColor(repo.Path, ct.Cyan)
		out, err := GitPull(repo.Path)
		if err != nil {
			// something wrong
			response := GitPullResponse{
				Path:  repo.Path,
				Body:  "",
				Error: err,
			}
			responses = append(responses, response)
		} else {
			if strings.Contains(out, "Already up-to-date") != true {
				// something new arrived
				response := GitPullResponse{
					Path:  repo.Path,
					Body:  out,
					Error: nil,
				}
				responses = append(responses, response)
			}
		}
	}

	for _, res := range responses {
		fmt.Println(res.Path)
		fmt.Println(res.Body)
		fmt.Println(res.Error)
	}

	return nil
}
