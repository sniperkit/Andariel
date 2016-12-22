package main

import (
	"github.com/google/go-github/github"
	"log"
)

func main() {
	client := github.NewClient(nil)

	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}

	repos, _, err := client.Repositories.List("fengyfei", opt)

	if err != nil {
		log.Print(err)
	}

	log.Print(repos)
}


