package main

import (
	"fmt"

	"Andariel/init"
	"Andariel/mongo"
	"Andariel/spider"
	. "Andariel/utility"
)

func initBeforeRun() {
	mongo.InitGithub()
	ini.InitMongoCollections()

	CsvService.ParseCsv()
	fmt.Print("Csv parse is over.\n")
}

func main() {
	initBeforeRun()

	spider.GetReposByAPI()
}
