package mongo

import (
	"log"

	"gopkg.in/mgo.v2"
)

var (
	GithubSession *mgo.Session
)

const (
	MDGitName = "github"
)

// 初始化 MongoDB 连接
func InitGithub() {
	url := "mongodb://10.0.0.254:27017"

	var err error
	GithubSession, err = mgo.Dial(url)

	if err != nil {
		panic(err)
	}

	log.Print("The MongoDB of GitHub connected.")

	// 利用 MongoDB 分布性特性
	GithubSession.SetMode(mgo.Monotonic, true)
}
