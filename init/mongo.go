package ini

import (
	"Andariel/utility"
	"Andariel/spider"
)

func InitMongoCollections() {
	utility.PrepareCsvParser()
	spider.PrepareRepos()
}
