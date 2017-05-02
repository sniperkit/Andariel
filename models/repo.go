/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/04/17        Yusan Kurban
 	   Update:  2017/04/26        Jia Chenhui
*/

package models

import (
	"errors"

	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"Andariel/mongo"
)

// 对外服务接口
type GitReposServiceProvider struct {
}

var GitReposService *GitReposServiceProvider
var GitReposCollection *mgo.Collection

// 连接、设置索引
func PrepareGitRepos() {
	GitReposCollection = mongo.GithubSession.DB(mongo.MDGitName).C("gitRepos")
	repoIndex := mgo.Index{
		Key:        []string{"FullName", "StarCount", "ForkCount", "Language"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}

	if err := GitReposCollection.EnsureIndex(repoIndex); err != nil {
		panic(err)
	}

	GitReposService = &GitReposServiceProvider{}
}

type Repos struct {
	ID              bson.ObjectId     `bson:"_id,omitempty",json:"id"`
	RepoID          uint64            `bson:"RepoID,omitempty" json:"repoid,omitempty"`
	Owner           bson.ObjectId     `bson:"Owner,omitempty" json:"-"`
	Name            string            `bson:"Name,omitempty" json:"name"`
	FullName        string            `bson:"FullName,omitempty" json:"fullname"`
	Description     string            `bson:"Description,omitempty" json:"description"`
	DefaultBranch   string            `bson:"DefaultBranch,omitempty" json:"defaultbranch"`
	Language        string            `bson:"Language,omitempty" json:"language"`
	Created         *github.Timestamp `bson:"Created,omitempty" json:"created"`
	Updated         *github.Timestamp `bson:"Updated,omitempty" json:"updated"`
	Pushed          *github.Timestamp `bson:"Pushed,omitempty" json:"pushed"`
	HasWiki         bool              `bson:"HasWiki,omitempty" json:"haswiki"`
	HasIssues       bool              `bson:"HasIssues,omitempty" json:"hasissues"`
	HasDownloads    bool              `bson:"HasDownloads,omitempty" json:"hasdownloads"`
	ForkCount       uint64            `bson:"ForkCount,omitempty" json:"forkcount"`
	StarCount       uint64            `bson:"StarCount,omitempty" json:"starcount"`
	WatchersCounts  uint64            `bson:"WatchersCounts,omitempty" json:"watcherscounts"`
	OpenIssuesCount uint64            `bson:"OpenIssuesCount,omitempty" json:"openissuescount"`
	Size            uint64            `bson:"Size,omitempty" json:"size"`
}

// 存储库信息及作者在 User 数据库中的 ID
func (rsp *GitReposServiceProvider) Create(repos *github.Repository, owner *string) error {
	r := Repos{
		ID:              bson.NewObjectId(),
		RepoID:          uint64(repos.ID),
		Owner:           bson.ObjectIdHex(string(owner)),
		Name:            string(repos.Name),
		FullName:        string(repos.FullName),
		Description:     string(repos.Description),
		DefaultBranch:   string(repos.DefaultBranch),
		Language:        string(repos.Language),
		Created:         repos.CreatedAt,
		Updated:         repos.UpdatedAt,
		Pushed:          repos.PushedAt,
		HasWiki:         bool(repos.HasWiki),
		HasIssues:       bool(repos.HasIssues),
		HasDownloads:    bool(repos.HasDownloads),
		ForkCount:       uint64(repos.ForksCount),
		StarCount:       uint64(repos.StargazersCount),
		WatchersCounts:  uint64(repos.WatchersCount),
		OpenIssuesCount: uint64(repos.OpenIssuesCount),
		Size:            uint64(repos.Size),
	}

	err := GitReposCollection.Insert(&r)
	if err != nil {
		return err
	}

	return nil
}

// 逻辑判断后，存储库信息到数据库
func StoreRepo(repo *github.Repository) (err error) {
	// fork 的库不保存
	// TODO: 此处的逻辑判断拿到外面，数据库中不保存 fork 的库，考虑如何处理
	if repo.Fork {
		err = errors.New("this repos is forked from others")
		return err
	}

	// 判断数据库中是否有此作者信息
	oldUserID, err := GitUserService.GetUserID(string(repo.Owner.Name))
	if err != nil {
		if err != mgo.ErrNotFound {
			return err
		}

		// User 数据库中无此作者信息
		newUserID, err := GitUserService.Create(repo.Owner)
		if err != nil {
			return err
		}

		err = GitReposService.Create(repo, &newUserID)
		if err != nil {
			return err
		}
	}

	// User 数据库中有此作者信息
	err = GitReposService.Create(repo, &oldUserID)
	if err != nil {
		return err
	}

	return nil
}
