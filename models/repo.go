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
 */

package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Repository struct {
	ID 					bson.ObjectId	`bson:"_id,omitempty",json:"_id"`
	RepoID 				uint64 			`json:"id,omitempty"`
	Owner				OwnerInfo		`json:"owner,omitempty"`
	Name				string			`json:"name,omitempty"`
	FullName			string			`json:"fullname,omitempty"`
	Description 		string			`json:"description,omitempty"`
	DefaultBranch 		string			`json:"defaultbranch,omitempty"`
	Language 			string			`json:"language,omitempty"`
	Created     		TimeStamp		`json:"created,omitempty"`
	Updated 			TimeStamp		`json:"updated,omitempty"`
	Pushed 				TimeStamp		`json:"pushed,omitempty"`
	
	// Urls
	ApiUrl 				string			`json:"url,omitempty"`
	HtmlUrl 			string			`json:"htmlurl,omitempty"`
	ArchiveUrl 			string			`json:"archiveurl,omitempty"`
	BranchesUrl			string 			`json:"branchesurl,omitempty"`
	CloneUrl 			string			`json:"cloneurl,omitempty"`
	CommentsUrl 		string			`json:"commentsurl,omitempty"`
	CommitsUrl			string			`json:"commitsurl,omitempty"`
	ContributorsUrl 	string			`json:"contributorsurl,omitempty"`
	DownloadUrl 		string			`json:"downloadurl,omitempty"`
	EventsUrl 			string			`json:"eventsurl,omitempty"`
	ForksUrl			string			`json:"forksurl,omitempty"`
	GitUrl 				string			`json:"giturl,omitempty"`
	IssuesUrl 			string 			`json:"issuesurl,omitempty"`
	LanguageUrl 		string 			`json:"languageurl,omitempty"`
	MergesUrl 			string 			`json:"mergesurl,omitempty"`
	MilestoneUrl		string 			`json:"milestoneurl,omitempty"`
	NotificationsUrl	string 			`json:"notificationsurl,omitempty"`
	SshUrl 				string			`json:"sshurl,omitempty"`
	SubscribersUrl	 	string			`json:"subscribersurl,omitempty"`
	SubscriptionUrl		string			`json:"subscriptionurl,omitempty"`
	
	Private 			bool            `json:"private,omitempty"`
	Fork 				bool            `json:"fork,omitempty"`
	HasWiki 			bool 			`json:"haswiki,omitempty"`
	HasIssues			bool 			`json:"hasissues,omitempty"`
	HasDownloads 		bool            `json:"hasdownloads,omitempty"`
	
	// Counts
	ForkCount 			uint64			`json:"forkcount,omitempty"`
	StarCount 			uint64			`json:"starcount,omitempty"`
	WatchersCounts		uint64 			`json:"watcherscounts,omitempty"`
	OpenIssuesCount 	uint64			`json:"openissuescount,omitempty"`
	Size 				uint64			`json:"size,omitempty"`
}

type TimeStamp struct {
	Time 				time.Time		`json:"time,omitempty"`
}

type OwnerInfo struct {
	OwnerID 			uint64			`json:"id,omitempty"`
	AvatarUrl 			string			`json:"avatarurl,omitempty"`
	OwnerUrl 			string			`json:"url,omitempty"`
	HtmlUrl 			string			`json:"htmlurl,omitempty"`
	ReposUrl			string			`json:"reposurl,omitempty"`
} 