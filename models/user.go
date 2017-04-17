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
	"time"

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
	UserID            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	HTMLURL           string        `bson:"HTMLURL,omitempty" json:"htmlurl"`
	Name              string        `bson:"Name,omitempty" json:"name"`
	Email             string        `bson:"Email,omitempty" json:"email"`
	PublicRepos       uint64        `bson:"PublicRepos,omitempty" json:"publicrepos"`
	PublicGists       uint64        `bson:"PublicGists,omitempty" json:"publicgists"`
	Followers         uint64        `bson:"Followers,omitempty" json:"followers"`
	Following         uint64        `bson:"Following,omitempty" json:"following"`
	CreatedAt         time.Time     `bson:"CreatedAt,omitempty" json:"created"`
	UpdatedAt         time.Time     `bson:"UpdatedAt,omitempty" json:"updated"`
	SuspendedAt       time.Time     `bson:"SuspendedAt,omitempty" json:"suspended"`
	Type              string        `bson:"Type,omitempty" json:"type"`
	TotalPrivateRepos uint64        `bson:"TotalPrivateRepos,omitempty" json:"totalprivaterepos"`
	OwnedPrivateRepos uint64        `bson:"OwnedPrivateRepos,omitempty" json:"ownedprivaterepos"`
	PrivateGists      uint64        `bson:"PrivateGists,omitempty" json:"privategists"`

	// API URLs
	URL               string `bson:"URL,omitempty" json:"url"`
	EventsURL         string `bson:"EventsURL,omitempty" json:"eventsurl"`
	FollowingURL      string `bson:"FollowingURL,omitempty" json:"followingurl"`
	FollowersURL      string `bson:"FollowersURL,omitempty" json:"followersurl"`
	GistsURL          string `bson:"GistsURL,omitempty" json:"gistsurl"`
	OrganizationsURL  string `bson:"OrganizationsURL,omitempty" json:"organizationsurl"`
	ReceivedEventsURL string `bson:"ReceivedEventsURL,omitempty" json:"receivedeventsurl"`
	ReposURL          string `bson:"ReposURL,omitempty" json:"reposurl"`
	StarredURL        string `bson:"StarredURL,omitempty" json:"starredurl"`
	SubscriptionsURL  string `bson:"SubscriptionsURL,omitempty" json:"subscriptionsurl"`

	// 用户在指定库上拥有的权限
	Permissions map[string]bool `bson:"Permissions,omitempty" json:"permissions"`
}
