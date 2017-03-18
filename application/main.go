package main

import (
	//"fmt"

	"Andariel/application/init"
	"Andariel/mongo"
	"Andariel/models/spider"
	//"Andariel/models/utility"
)

func initBeforeRun() {
	mongo.InitGithub()
	ini.InitMongoCollections()

    /*
	util.CsvService.ParseCsv()
	fmt.Print("Csv parsing is over.\n")*/
}

func main() {
	initBeforeRun()
	spider.RequestService.Chanel()
}
