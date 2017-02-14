package ini

import (
	"Andariel/utility"
	"Andariel/spider"
)

func InitMongoCollections() {
	util.PrepareCsvParser()
	spider.PrepareRepos()
}
