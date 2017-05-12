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
func InitGithub(url string) {
	var err error
	GithubSession, err = mgo.Dial(url)

	if err != nil {
		panic(err)
	}

	// 利用 MongoDB 分布性特性
	GithubSession.SetMode(mgo.Monotonic, true)
}
