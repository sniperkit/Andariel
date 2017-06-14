package mongo

import (
	"gopkg.in/mgo.v2"
)

type Session struct {
	MongoSession *mgo.Session
}

// 初始化 MongoDB 连接
func InitMongoSes(url string) *Session {
	ses, err := mgo.Dial(url)

	if err != nil {
		panic(err)
	}

	// 利用 MongoDB 分布性特性
	ses.SetMode(mgo.Monotonic, true)

	return &Session{
		MongoSession: ses,
	}
}
