package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"log"
	"strconv"

	. "Andariel/utility"
)

var session *mgo.Session
var collection *mgo.Collection

func init() {
	var err error
	session, err = mgo.Dial("mongodb://10.0.0.254:27017")

	if err != nil {
		panic(err)
	}

	collection = session.DB("github").C("repos")

	CsvService.ParseCsv()
	fmt.Print("Csv parse is over.")
}

func main() {
	client := github.NewClient(nil)

	// 获取数据库中所有 csv 记录
	results, err := CsvService.GetAllRecords()

	if err != nil {
		panic(err)
	}

	for _, result := range results {
		reposID, err := strconv.Atoi(result.ReposID)

		if err != nil {
			log.Fatal(err)
			continue
		}

		repo, _, err := client.Repositories.GetByID(int(reposID))

		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
		}

		if err != nil {
			panic(err)
		}

		err = collection.Insert(repo)

		if err != nil {
			log.Print(err)
		}
	}
}
