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
 *     Initial: 17/04/2017        Jia Chenhui
 */

package models

import (
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"Andariel/mongo"
)

// 对外服务接口
type GitUserServiceProvider struct {
}

var GitUserService *GitUserServiceProvider
var GitUserCollection *mgo.Collection

// 连接、设置索引
func PrepareGitUser() {
	GitUserCollection = mongo.GithubSession.DB(mongo.MDGitName).C("gitUser")
	userIndex := mgo.Index{
		Key:        []string{"Name"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}

	if err := GitUserCollection.EnsureIndex(userIndex); err != nil {
		panic(err)
	}

	GitUserService = &GitUserServiceProvider{}
}

// GitHub 用户数据结构
type User struct {
	UserID            bson.ObjectId     `bson:"_id,omitempty" json:"id"`
	ID                uint64            `bson:"ID,omitempty" json:"userid"`
	HTMLURL           string            `bson:"HTMLURL,omitempty" json:"htmlurl"`
	Name              string            `bson:"Name,omitempty" json:"name"`
	Email             string            `bson:"Email,omitempty" json:"email"`
	PublicRepos       uint64            `bson:"PublicRepos,omitempty" json:"publicrepos"`
	PublicGists       uint64            `bson:"PublicGists,omitempty" json:"publicgists"`
	Followers         uint64            `bson:"Followers,omitempty" json:"followers"`
	Following         uint64            `bson:"Following,omitempty" json:"following"`
	CreatedAt         *github.Timestamp `bson:"CreatedAt,omitempty" json:"created"`
	UpdatedAt         *github.Timestamp `bson:"UpdatedAt,omitempty" json:"updated"`
	SuspendedAt       *github.Timestamp `bson:"SuspendedAt,omitempty" json:"suspended"`
	Type              string            `bson:"Type,omitempty" json:"type"`
	TotalPrivateRepos uint64            `bson:"TotalPrivateRepos,omitempty" json:"totalprivaterepos"`
	OwnedPrivateRepos uint64            `bson:"OwnedPrivateRepos,omitempty" json:"ownedprivaterepos"`
	PrivateGists      uint64            `bson:"PrivateGists,omitempty" json:"privategists"`
}

// 查询作者信息
func (usp *GitUserServiceProvider) GetUserByID(userID uint64) (User, error) {
	var u User

	err := GitUserCollection.Find(bson.M{"ID": userID}).One(&u)

	if err != nil {
		return nil, err
	}

	return u, nil
}

// 通过 name 获取作者在数据库中的 _id
func (usp *GitUserServiceProvider) GetUserID(name string) (string, error) {
	var u User

	err := GitUserCollection.Find(bson.M{"Name": name}).One(&u)

	if err != nil {
		return "", err
	}

	return u.UserID.Hex(), nil
}

// 通过 name 判断作者是否存在数据库中
func (usp *GitUserServiceProvider) IsUserExists(name string) bool {
	uID, _ := usp.GetUserID(name)

	return uID != ""
}

// 存储作者信息
func (usp *GitUserServiceProvider) Create(user *github.User) (string, error) {
	u := User{
		UserID:            bson.NewObjectId(),
		ID:                uint64(user.ID),
		HTMLURL:           string(user.HTMLURL),
		Name:              string(user.Name),
		Email:             string(user.Email),
		PublicRepos:       uint64(user.PublicRepos),
		PublicGists:       uint64(user.PublicGists),
		Followers:         uint64(user.Followers),
		Following:         uint64(user.Following),
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		SuspendedAt:       user.SuspendedAt,
		Type:              string(user.Type),
		TotalPrivateRepos: uint64(user.TotalPrivateRepos),
		OwnedPrivateRepos: uint64(user.OwnedPrivateRepos),
		PrivateGists:      uint64(user.PrivateGists),
	}

	err := GitUserCollection.Insert(&u)

	if err != nil {
		return nil, err
	}

	return u.UserID.Hex(), nil
}
