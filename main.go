package main

import (
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"log"
)

var session     *mgo.Session
var collection  *mgo.Collection

func init() {
	var err error
	session, err = mgo.Dial("mongodb://127.0.0.1")

	if err != nil {
		panic(err)
	}

	collection = session.DB("github").C("repos")
}

func main() {
	client := github.NewClient(nil)

	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}

	repos, _, err := client.Repositories.List("fengyfei", opt)

	if err != nil {
		panic(err)
	}

	for _, repo := range repos {
		err = collection.Insert(repo)

		if err != nil {
			log.Print(err)
		}
	}
}


