package ini

import (
	"Andariel/models/utility"
	"Andariel/models/spider"
)

func InitMongoCollections() {
	util.PrepareCsvParser()
	spider.PrepareRepos()
	spider.PrepareStar()
}
