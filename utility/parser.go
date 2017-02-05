package utility

import (
	"encoding/csv"
	"io"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// 对外服务接口
type CsvServiceProvider struct {
}

var CsvService *CsvServiceProvider

var session *mgo.Session
var collection *mgo.Collection

type CsvParser struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	ReposID   string        `bson:"Id"`
	FullName  string        `bson:"FullName"`
	StarCount string        `bson:"StarCount"`
	Language  string        `bson:"Language"`
	Domain    string        `bson:"Domain"`
	Growth    string        `bson:"Growth"`
	Tf        string        `bson:"Tf"`
}

// 连接数据库
func init() {
	var err error
	session, err = mgo.Dial("mongodb://10.0.0.254:27017")

	if err != nil {
		panic(err)
	}

	collection = session.DB("github").C("csv")
}

// 解析 csv 文件并存入数据库
func (this *CsvServiceProvider) ParseCsv() {
	file, err := os.Open("/Users/LLLeon/Desktop/repositories.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	reader.LazyQuotes = true
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		r := CsvParser{
			Id:        bson.NewObjectId(),
			ReposID:   record[0],
			FullName:  record[1],
			StarCount: record[2],
			Language:  record[3],
			Domain:    record[4],
			Growth:    record[5],
			Tf:        record[6],
		}
		err = collection.Insert(&r)

		if err != nil {
			log.Println(err)
		}
	}
}

// 获取数据库中所有记录的 id
func (this *CsvServiceProvider) GetAllRecords() ([]CsvParser, error) {
	var results []CsvParser

	err := collection.Find(nil).All(&results)

	if err != nil {
		panic(err)
	}

	return results, err
}

// 根据 ID 获取库的数字编号
func (this *CsvServiceProvider) GetReposIDByID(id string) (string, error) {
	var csv CsvParser

	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&csv)

	return csv.ReposID, err
}
