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
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"log"
	"time"
)

const (
	empty = 0
	limit = 5000
)

type GithubClient struct {
	Client      *github.Client
	StartAt     time.Time
	LimitAt     time.Time
	RequestTime time.Duration
	Timer       time.Duration
	Limited     bool
	Times       int
	Left        int
	Count       func()
}

var GitClient *GithubClient = new(GithubClient)

func newClient(token string) (client *GithubClient) {
	client = new(GithubClient)
	client.init(token)

	return client
}

func (this *GithubClient) init(token string) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(httpClient)
	this.Client = client
	_, ok := this.requestTimes()
	if ok {
		this.StartAt = time.Now()
	}
}

func (this *GithubClient) checkLimit() bool {
	return this.Limited
}

func (this *GithubClient) onErr() error {
	var err error
	if s, ok := err.(*github.RateLimitError); ok {
		reset := this.StartAt.Add(time.Hour * 1)
		this.LimitAt = time.Now()
		this.Limited = true
		this.RequestTime = this.LimitAt.Sub(this.StartAt)
		this.Timer = reset.Sub(this.LimitAt)
		e := github.RateLimitError.Error(*s)
		return errors.New(e)
	}

	return nil
}

func (this *GithubClient) reset() {
	this.Times = empty
	this.Left = limit
	this.Limited = false
}

func (this *GithubClient) requestTimes() (error, bool) {
	rate, _, err := this.Client.RateLimits(oauth2.NoContext)
	if err != nil {
		log.Println("Get limits crash with error:", err)
		return err, false
	}
	this.Times = rate.Core.Limit - rate.Core.Remaining
	this.Left = rate.Core.Remaining

	if this.Left != limit-1 {
		return nil, false
	}

	return nil, true
}

func (this *GithubClient) monitor() {
	for {
		if this.Left == empty {
			if this.Timer == empty {
				this.reset()
			}
		}
	}
}
