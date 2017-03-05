package spider

import (
	"Andariel/mongo"
	"gopkg.in/mgo.v2"
	"github.com/andygrunwald/go-trending"

	"log"
	"time"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"github.com/google/go-github/github"
	"fmt"
)

type RequestServiceProvider struct {
}

var RequestService *RequestServiceProvider

var TrendingCollection *mgo.Collection
var PopularCollection *mgo.Collection


func PrepareTren() {
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
		Background: 	true,
		Sparse: 	true,
	}

	if err:= PopularCollection.EnsureIndex(idIndex); err != nil {
		panic(err)
	}
}
type Trending struct {
	CreateTime	string	 		`json:"create_time"`
	Repos 		[]trending.Project	`json:"repos"`
}

type Popular struct {
	ParseTime 	time.Time 			`json:"parse_time"`
	Repos 		*github.RepositoriesSearchResult	`json:"repos"`
}

// 获取一天的 trending
func (this *RequestServiceProvider) GetTrendingToday(l string) {

	trend := trending.NewTrending()
	// trending.TimeTodya 可以换成TimeWeek or TimeMonth 来获取本周或本月的 trending .
	result, err := trend.GetProjects(trending.TimeToday, l)

	if err != nil {
		log.Print(err)
	}
	t := time.Now().Format("20060102")
	i := Trending{
		CreateTime: 	t,
		Repos: 		result,
	}

	err = TrendingCollection.Insert(i)

	if err != nil {
		log.Print(err)
	}
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
		Repos: 		result,
	}

	err = TrendingCollection.Insert(i)

	if err != nil {
		log.Print(err)
	}
}

//从数据库获取 trending
func (this *RequestServiceProvider) GetTrendingFromMD(t string, l string) ([]trending.Project, error) {
	var m Trending

	err := TrendingCollection.Find(bson.M{"createtime":t, "language": l}).One(&m)

	if err != nil {
		log.Print(err)
	}
	return m.Repos, err
}

//获取 popular 库
func (this *RequestServiceProvider) GetPopular(l string) {
	var client *github.Client

	opt := &github.SearchOptions{Sort: "stars"}
	query := fmt.Sprintf("tetris+language:%s",l)
	result, _, err := client.Search.Repositories(query, opt)
	if err != nil {
		log.Print(err)
	} else {
		r := Popular{
			ParseTime: 	time.Now(),
			Repos: 		result,
		}
		err = PopularCollection.Insert(r)
		if err != nil {
			log.Print(err)
		}
	}
}

// 从数据库获取 popular
func (this *RequestServiceProvider) GetPopularFromDB() (Popular,error) {
	var r Popular

	err := PopularCollection.Find(bson.M{}).Sort("-parse_time").One(&r)
	if err != nil {
		log.Print(err)
	}
	return r, err
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