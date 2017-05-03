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

package github

import (
	"Andariel/log"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"time"
)

var logger *log.AndarielLogger = log.AndarielCreateLogger(
	&log.AndarielLogTag{
		log.LogTagService: "github",
		log.LogTagType: "client",
	},
	log.AndarielLogLevelDefault)


const (
	empty              		= 0
	authNonSearchLimit 		= 5000	//per hour
	nonAuthNonSearchLimit 	= 60
	authSearchLimit	   		= 30 	//per minutes
	nonAuthSearchLimit 		= 10

	authClient				= 0x01
	nonAuthClient			= 0x10

	core					= "core"
	search 					= "search"
)

type GithubClient struct {
	Client      *github.Client
	StartAt     time.Time
	LimitAt     time.Time
	RequestTime time.Duration
	ClientType  int
	Core 		Rate
	Search 		Rate
}

type Rate struct {
	Times		int
	Limit 		int
	Left 		int
	Limited 	bool
	Reset 		time.Time
	ResetIn 	time.Duration
}

var GitClient *GithubClient = newClient("")

func newClient(token string) (client *GithubClient) {
	if token == "" {
		client = new(GithubClient)
		tokenSource := *new(oauth2.TokenSource)
		client.ClientType = nonAuthClient
		client.init(tokenSource)
	} else {
		client = new(GithubClient)
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client.ClientType = authClient
		client.init(tokenSource)
	}

	return client
}

func (this *GithubClient) init(tokenSource oauth2.TokenSource) {
	httpClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(httpClient)
	this.Client = client
	err, ok := this.requestTimes()

	if err != nil {
		logger.Error("Get limits crash with error:", err)
		return
	}

	if ok {
		this.StartAt = time.Now()
	}
}

func (this *GithubClient) checkLimit(limit string) bool {
	switch {
	case limit == search:
		return this.Search.Limited
	default:
		return this.Core.Limited
	}

}

func (this *GithubClient) onErr() error {
	var err error
	if s, ok := err.(*github.RateLimitError); ok {
		reset := this.StartAt.Add(time.Hour * 1)
		this.LimitAt = time.Now()
		this.Core.Limited = true
		this.RequestTime = this.LimitAt.Sub(this.StartAt)
		this.Core.ResetIn = reset.Sub(this.LimitAt)
		e := (github.RateLimitError)(*s)
		return errors.New(e.Message)
	}

	return nil
}

func (this *GithubClient) reset(limit string) {
	switch {
	case limit == search:
		if this.ClientType == authClient {
			this.Search.Left = authSearchLimit
		} else {
			this.Search.Left = nonAuthSearchLimit
		}
		this.Search.Times = empty
		this.Search.Limited = false
	default:
		if this.ClientType == authClient {
			this.Core.Left = authNonSearchLimit
		} else {
			this.Core.Left = nonAuthNonSearchLimit
		}
		this.Core.Times = empty
		this.Core.Limited = false
	}
}

func (this *GithubClient) requestTimes() (error, bool) {
	rate, _, err := this.Client.RateLimits(oauth2.NoContext)
	if err != nil {
		return err, false
	}
	this.Core.Times = rate.Core.Limit - rate.Core.Remaining
	this.Core.Left = rate.Core.Remaining
	this.Core.Reset = rate.Core.Reset.Time
	this.Core.ResetIn = rate.Core.Reset.Sub(time.Now())


	if this.Core.Left != authNonSearchLimit- 1 {
		return nil, false
	}

	return nil, true
}

func (this *GithubClient) monitor() {
	for {
		this.onErr()
		if this.Core.Left == empty {
			if this.Core.ResetIn == empty {
				this.reset(core)
			}
		}
		if this.Search.Left == empty {
			if this.Search.ResetIn == empty{
				this.reset(search)
			}
		}
	}
}
