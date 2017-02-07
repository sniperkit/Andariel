package main

import (
	. "Andariel/utility"
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"log"

	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
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
	fmt.Print("Csv parse is over.\n")
}

func main() {

	// 添加身份验证, 提高请求速率
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "54f7488c8f72d3e63692b2bf04167d97e7a29e1d"},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	// 获取所有解析 csv 文件后所得的 id 记录
	results, err := CsvService.GetAllRecords()

	if err != nil {
		panic(err)
	}

	for _, result := range results {
		repo, _, err := client.Repositories.GetByID(result.ReposID)

		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("Hit rate limit.")
		}

		// 如果获取库时报错, 则跳过该库
		if err != nil {
			log.Print(err)
			continue
		}

		// 如果是 fork 的库则跳过
		if *repo.Fork == true {
			log.Println("This repository is forked.")

			continue
		} else {
			_, err = collection.Upsert(bson.M{"id": result.ReposID}, repo)

			if err != nil {
				log.Print(err)
			}
		}
	}
}
