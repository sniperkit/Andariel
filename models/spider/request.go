package spider

import (
	"Andariel/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/robfig/cron"
	"github.com/google/go-github/github"

	"log"
)

type RequestServiceProvider struct {
}

var RequestService *RequestServiceProvider

var RepoStarCollection *mgo.Collection

type StarCount struct {
	ID 		int 	`json:"id"`
	FullName	*string	`json:"full_name"`
	StarOld		*int	`json:"star_old"`
	StarNew 	*int 	`json:"star_new"`
	StarChange	int 	`json:"star_change"`
	ForkOld 	*int 	`json:"fork_old"`
	ForkNew		*int	`json:"fork_new"`
	ForkChange 	int	`json:"fork_change"`
}

type Count struct {
	StarNew		int
	StarOld		int
	ForkNew		int
	ForkOld		int
}
type RepoId struct {
	Id		int 	`json:"id"`
}
func PrepareStar() {
	RepoStarCollection = mongo.GithubSession.DB(mongo.MDGitName).C("stars")
	idIndex := mgo.Index{
		Key:        []string{"fullname"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := ReposCollection.EnsureIndex(idIndex); err != nil {
		panic(err)
	}

	SpiderService = &SpiderServiceProvider{}
}

// 获取repos表里的库id
func GetReposId()([]RepoId, error)  {
	var i []RepoId

	err := ReposCollection.Find(bson.M{}).All(&i)

	return i, err
}

// 一周一次爬取之前把star数放到star表里
func GetReposPerWeekBefore() {
	var re  github.Repository

	results, err := GetReposId()
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		ok := ReposCollection.Find(bson.M{"id":result.Id }).One(&re)
		if ok != nil {
			log.Print(err)
		}

		selector := bson.M{"id": result.Id}
		data := bson.M{"$set": bson.M{"starold": re.StargazersCount,
			"forkold": re.ForksCount}}

		_, err = RepoStarCollection.Upsert(selector, data)
		/*repo := StarCount{
			ID:		result.Id,
			FullName: 	re.FullName,
			StarOld: 	re.StargazersCount,
			ForkOld: 	re.ForksCount,
		}
		_, err = RepoStarCollection.Upsert(bson.M{"id": result.Id},repo)*/
		if err != nil {
			log.Print(err)
		}

	}
}

func GetReposPerWeekAfter() {
	var re  github.Repository

	results, err := GetReposId()
	if err != nil {
		panic(err)
	}

	for _, result := range results {

		ok := ReposCollection.Find(bson.M{"id":result.Id }).One(&re)
		if ok != nil {
			log.Print(err)
		}

		selector := bson.M{"id": result.Id}
		data := bson.M{"$set": bson.M{"fullname":re.FullName,
			"starnew":re.StargazersCount,
			"forknew": re.ForksCount}}

		_, err = RepoStarCollection.Upsert(selector, data)
		if err != nil {
			log.Print(err)
		}

	}
}

func UpdateChange() {
	var s Count
	results, err := GetReposId()
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		ok := RepoStarCollection.Find(bson.M{"id": result.Id}).One(&s)
		if ok != nil {
			log.Print(ok)
		}

		selector := bson.M{"id": result.Id}
		data := bson.M{"$set": bson.M{"starchange": s.StarNew-s.StarOld,
					"forkchange": s.ForkNew-s.ForkOld}}
		_, err = RepoStarCollection.Upsert(selector, data)
	}
}

//包装
func (this *RequestServiceProvider)GetRepo() {
	GetReposPerWeekBefore()
	GetReposByAPI()
	GetReposPerWeekAfter()
	UpdateChange()
}

func (this *RequestServiceProvider)CronJob() {
	c := cron.New()
	c.AddFunc("@daily", RequestService.GetRepo)
	c.Start()
	select {}
}
