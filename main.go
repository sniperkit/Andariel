package main

import (
	"fmt"

	"Andariel/init"
	"Andariel/mongo"
	"Andariel/spider"
	"Andariel/utility"
)

func initBeforeRun() {
	mongo.InitGithub()
	ini.InitMongoCollections()

	util.CsvService.ParseCsv()
	fmt.Print("Csv parsing is over.\n")
}

func main() {
	initBeforeRun()

	spider.GetReposByAPI()
}
