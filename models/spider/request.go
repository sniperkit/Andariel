package spider

import (
	"Andariel/mongo"
	"gopkg.in/mgo.v2"
	"github.com/andygrunwald/go-trending"
	"gopkg.in/mgo.v2/bson"
	"github.com/google/go-github/github"
	"github.com/nats-io/go-nats"

	"log"
	"time"
	"strconv"
	"fmt"
	"encoding/json"
)

type RequestServiceProvider struct {
}

var RequestService *RequestServiceProvider

var TrendingCollection *mgo.Collection
var PopularCollection *mgo.Collection


func PrepareTren() {
	// 全局变量
	RequestService = &RequestServiceProvider{}

	TrendingCollection = mongo.GithubSession.DB(mongo.MDGitName).C("trending")
	idIndex := mgo.Index{
		Key: 		[]string{"name"},
		Unique: 	true,
		DropDups: 	true,
		Background: 	true,
		Sparse: 	true,
	}

	if err := TrendingCollection.EnsureIndex(idIndex); err != nil {
		panic(err)
	}
}

func PreparePop() {
	PopularCollection = mongo.GithubSession.DB(mongo.MDGitName).C("popular")
	idIndex := mgo.Index{
		Key: 		[]string{"name"},
		Unique: 	true,
		DropDups: 	true,
		Background: true,
		Sparse: 	true,
	}

	if err:= PopularCollection.EnsureIndex(idIndex); err != nil {
		panic(err)
	}
}
type Trending struct {
	CreateTime	string	 		`json:"createtime"`
	Language 	string 			`json:"language"`
	Repos 		[]trending.Project	`json:"repos"`
}

type Popular struct {
	ParseTime 	time.Time 			`json:"parsetime"`
	Language 	string 				`json:"language"`
	Repos 		*github.RepositoriesSearchResult	`json:"repos"`
}

type Gitspider struct {
	Type 	string
}

type Gitoffice struct {
	Type 		string
	Time 		string
	Language 	string
}
// 获取一天的 trending
func (this *RequestServiceProvider) GetTrendingToday(l string) {

	trend := trending.NewTrending()
	// trending.TimeTodya 可以换成TimeWeek or TimeMonth 来获取本周或本月的 trending .
	result, err := trend.GetProjects(trending.TimeToday, l)

	if err != nil {
		log.Print("first error")
		log.Print(err)
	}
	t := time.Now().Format("20060102")
	i := Trending{
		CreateTime: 	t,
		Language: 		l,
		Repos: 		result,
	}

	err = TrendingCollection.Insert(i)

	if err != nil {
		log.Print("third error")
		log.Print(err)
	}
	log.Print("complate.")
}

func (this *RequestServiceProvider) GetTrendingWeek(l string) {

	trend := trending.NewTrending()
	result, err := trend.GetProjects(trending.TimeWeek, l)

	if err != nil {
		log.Print(err)
	}
	_, w := time.Now().ISOWeek()
	t := strconv.Itoa(w)
	i := Trending{
		CreateTime: 	t,
		Language: 		l,
		Repos: 		result,
	}

	err = TrendingCollection.Insert(i)

	if err != nil {
		log.Print(err)
	}
}

func (this *RequestServiceProvider) GetTrendingMonth(l string) {

	trend := trending.NewTrending()
	result, err := trend.GetProjects(trending.TimeMonth, l)

	if err != nil {
		log.Print(err)
	}
	m := time.Now().Month().String()
	i := Trending{
		CreateTime: 	m,
		Language: 		l,
		Repos: 		result,
	}

	err = TrendingCollection.Insert(i)

	if err != nil {
		log.Print(err)
	}
}

//从数据库获取 trending
func (this *RequestServiceProvider) GetTrendingFromMD(t string, l string) (Trending, error) {
	var m Trending

	err := TrendingCollection.Find(bson.M{"createtime":t, "language": l}).One(&m)

	if err != nil {
		log.Print(err)
	}
	return m, err
}

// 从数据库获取 popular
func (this *RequestServiceProvider) GetPopularFromMD(l string) (Popular,error) {
		var r Popular

		err := PopularCollection.Find(bson.M{"language": l}).Sort("-parsetime").One(&r)
		if err != nil {
		log.Print(err)
	}
	return r, err
}

//获取 popular 库
func (this *RequestServiceProvider) GetPopular(l string) {
	var client *github.Client

	nextClient := ChangeToken()
	client = nextClient()

	if client == nil {
		log.Print("Token has run out.")
	}

	opt := &github.SearchOptions{Sort: "stars"}
	query := fmt.Sprintf("language:%s",l)
	log.Print(query)
	result, _, err := client.Search.Repositories(query, opt)
	if err != nil {
		log.Print("4th error")
		log.Print(err)
	} else {
		r := Popular{
			ParseTime: 	time.Now(),
			Language: 	l,
			Repos: 		result,
		}
		err = PopularCollection.Insert(r)
		if err != nil {
			log.Print("second error.")
			log.Print(err)
		}
	}
	log.Print("complate too.")
}

// 指定语言获取
func (this *RequestServiceProvider) GetTrendingByLanguage() {
	arr := [6]string{"", "go", "python", "js", "swift", "html"}

	for i := 0; i < 6; i++ {
		RequestService.GetTrendingToday(arr[i])
		RequestService.GetPopular(arr[i])
	}
}

// 指定语言一周获取一次trending
func (this *RequestServiceProvider) GetTrendingByLanguageWeek() {
	arr := [6]string{"","go", "python", "js", "swift", "html"}

	for i := 0; i < 6; i++ {
		RequestService.GetTrendingWeek(arr[i])
	}
}

// 指定语言一个月获取一次trending
func (this *RequestServiceProvider) GetTrendingByLanguageMonth() {
	arr := [6]string{"","go", "python", "js", "swift", "html"}

	for i := 0; i < 6; i++{
		RequestService.GetTrendingMonth(arr[i])
	}
}

// sub 的实现
func (this *RequestServiceProvider) Chanel() {
	var s Gitspider
	var o Gitoffice
	nc ,err := nats.Connect("nats://10.0.0.254:4223")
	if err != nil {
		log.Print("Cann`t connect:",err)
	}
	nc.Subscribe("git/spider", func(msg *nats.Msg) {

		err := json.Unmarshal(msg.Data, &s)
		if err != nil {
			log.Print("unmarshall error")
			log.Print(err)
		}
		switch {
		case s.Type == "gettrendingbylanguage":
			RequestService.GetTrendingByLanguage()
		case s.Type == "gettrendingbylanguageweek":
			RequestService.GetTrendingByLanguageWeek()
		case s.Type == "gettrendingbylanguagemonth":
			RequestService.GetTrendingByLanguageMonth()
		}
	})
	nc.Subscribe("git/office", func(msg *nats.Msg){
			err := json.Unmarshal(msg.Data, &o)
			if err != nil {
				log.Print(err)
			}
			log.Print("git/office...", o)
			switch {
			case o.Type == "gettrendingfrommd":
				log.Print("gettrendingfrommd start.")
				trend, err := RequestService.GetTrendingFromMD(o.Time, o.Language)
				if err != nil {
					log.Print(err)
				}
				s, er := json.Marshal(trend)
				if er != nil {
					log.Print(er)
				}
				nc.Publish(msg.Reply, s)
			case o.Type == "getpopularfrommd":
				pop, err := RequestService.GetPopularFromMD(o.Language)
				if err != nil {
					log.Print(err)
				}
				p, er := json.Marshal(pop.Repos)
				if er != nil {
					log.Print(er)
				}
				nc.Publish(msg.Reply, p)
			default:
				log.Print("Not Found Func.")
			}

	})
	nc.Flush()

	select {}
}