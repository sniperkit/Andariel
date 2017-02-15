package spider

import (
	"errors"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"Andariel/mongo"
	"Andariel/utility"
)

// 对外服务接口
type SpiderServiceProvider struct {
}

var SpiderService *SpiderServiceProvider

var ReposCollection *mgo.Collection

// 连接、设置索引
func PrepareRepos() {
	ReposCollection = mongo.GithubSession.DB(mongo.MDGitName).C("repos")
	idIndex := mgo.Index{
		Key:        []string{"fullname"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := ReposCollection.EnsureIndex(idIndex); err != nil {
		panic(err)
	}

	SpiderService = &SpiderServiceProvider{}
}

// 切换 token
func changeToken() func() *github.Client {
	var client *github.Client

	i := 0
	tokens := []string{
		"54f7488c8f72d3e63692b2bf04167d97e7a29e1d",
		"5511599ff7aebf94476ce3eda7741ab7ae797ef9",
		"78d1dcb42b8c4368884603cfcd4f3a1581d771d2",
		"5df193b89001e9fabfdb947a88cdd8b6e45378f5",
	}

	return func() *github.Client {
		if i < len(tokens) {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: tokens[i]},
			)
			tc := oauth2.NewClient(oauth2.NoContext, ts)
			client = github.NewClient(tc)

			i += 1

			return client
		} else {

			return nil
		}
	}
}

// 调用 GitHub 官方 API 获取库信息
func GetReposByAPI() {
	var client *github.Client

	// 循环设置 token, 提高请求速率
	nextClient := changeToken()
	client = nextClient()

	if client == nil {
		errors.New("Tokens has run out.")
	}

	// 获取所有解析 csv 文件后所得的 id 记录
	results, err := util.CsvService.GetAllRecords()

	if err != nil {
		panic(err)
	}

	for _, result := range results {
		repo, _, err := client.Repositories.GetByID(result.ReposID)

		// 遇到速率限制错误后, 切换 token
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("Hit rate limit.")
			client = nextClient()

			if client == nil {
				errors.New("Tokens has run out.")
			}

			log.Println("Change token, get this repository again.")
			repo, _, err = client.Repositories.GetByID(result.ReposID)
		}

		if err != nil {
			log.Print(err)
			continue
		}

		if *repo.Fork == true {
			log.Println("This repository is forked.")

			continue
		} else {
			_, err = ReposCollection.Upsert(bson.M{"id": result.ReposID}, repo)

			if err != nil {
				log.Print(err)
			}
		}
	}
}
