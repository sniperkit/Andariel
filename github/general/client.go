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

package general

import (
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	ctn "Andariel/common/container"
	"Andariel/log"
)

const (
	empty              = 0
	emptyDuration      = time.Duration(0)
	authCoreLimit      = 5000 //per hour
	nonAuthCoreLimit   = 60
	authSearchLimit    = 30 //per minutes
	nonAuthSearchLimit = 10

	authClient    = 0x01
	nonAuthClient = 0x10

	coreType   = "core"
	searchType = "search"

	invalidTokenErr = "401 Unauthorized"
)

var logger *log.AndarielLogger = log.AndarielCreateLogger(
	&log.AndarielLogTag{
		log.LogTagService: "github",
		log.LogTagType:    "client",
	},
	log.AndarielLogLevelDefault)

type GithubClient struct {
	Client  *github.Client
	LimitAt time.Time
	Type    int
	UseType string
	Rate
}

type Rate struct {
	Times     int
	Limit     int
	Remaining int
	Limited   bool
	Reset     time.Time
	ResetIn   time.Duration
}

func newClient(token string) (client *GithubClient) {
	if token == "" {
		client = new(GithubClient)
		tokenSource := new(oauth2.TokenSource)
		client.Type = nonAuthClient
		client.init(tokenSource)
	} else {
		client = new(GithubClient)
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client.Type = authClient
		client.init(tokenSource)
	}

	return client
}

func (gc *GithubClient) init(tokenSource oauth2.TokenSource) {
	httpClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(httpClient)
	gc.Client = client
	if !gc.isValidToken(httpClient) {
		logger.Debug("Invalid token.")
		return
	}

	err, _ := gc.requestTimes()
	if err != nil {
		logger.Error("Get limits crash with error:", err)
		return
	}
}

func (gc *GithubClient) isValidToken(c *http.Client) bool {
	req, err := gc.Client.NewRequest("GET", "", nil)
	if err != nil {
		logger.Error("crash with:", err)
		return false
	}

	resp, err := c.Do(req)
	if err != nil {
		logger.Error("crash with:", err)
		return false
	}

	if resp.Header.Get("Status") == invalidTokenErr {
		return false
	}

	return true
}

func (gc *GithubClient) onLimit(resp *github.Response) error {
	if resp != nil && resp.Remaining <= 1 {
		gc.LimitAt = time.Now()
		gc.Times = resp.Limit - resp.Remaining
		gc.Limit = resp.Limit
		gc.Remaining = empty
		gc.Limited = true
		gc.ResetIn = resp.Reset.Time.Sub(gc.LimitAt)
		gc.Reset = resp.Reset.Time
	}

	return nil
}

func (gc *GithubClient) reset(useType string) {
	if gc.Type == authClient {
		switch useType {
		case coreType:
			gc.Remaining = authCoreLimit
		case searchType:
			gc.Remaining = authSearchLimit
		}
	}

	if gc.Type == nonAuthClient {
		switch useType {
		case coreType:
			gc.Remaining = nonAuthCoreLimit
		case searchType:
			gc.Remaining = nonAuthSearchLimit
		}
	}

	gc.Times = empty
	gc.Limited = false
}

func (gc *GithubClient) requestTimes() (error, bool) {
	rate, _, err := gc.Client.RateLimits(oauth2.NoContext)
	if err != nil {
		return err, false
	}

	// TODO: 此处未分类, 需修改
	gc.Times = rate.Core.Limit - rate.Core.Remaining
	gc.Remaining = rate.Core.Remaining
	gc.Reset = rate.Core.Reset.Time
	gc.ResetIn = rate.Core.Reset.Sub(time.Now())

	if gc.Remaining != authCoreLimit-1 {
		return nil, false
	}

	return nil, true
}

func (gc *GithubClient) monitor(useType string, resp *github.Response) {
	gc.onLimit(resp)

	if gc.Remaining == empty {
		if gc.Reset.Sub(time.Now()) == emptyDuration {
			gc.reset(useType)
		}
	}
}

func (gc *GithubClient) checkLimit() bool {
	return gc.Limited
}

func (gc *GithubClient) isValidClient() bool {
	if gc.Remaining == empty {
		return false
	}

	return true
}

type ClientRing struct {
	*ctn.Ring
}

// 创建 client 并存储到 ring
func (cr *ClientRing) newClientRing(tokens []string) error {
	var clients []*GithubClient

	for _, t := range tokens {
		client := newClient(t)
		clients = append(clients, client)
	}

	clientRing := ctn.NewRing(len(tokens))

	err := clientRing.MPush(clients)
	if err != nil {
		return err
	}

	cr.Ring = clientRing

	return nil
}

// 存储 client 到 ring
func (cr *ClientRing) push(client *GithubClient) error {
	err := cr.Push(client)
	if err != nil {
		return err
	}

	return nil
}

// 从 ring 内读取 client
func (cr *ClientRing) get() *GithubClient {
	client := cr.Get()

	c, ok := client.(*GithubClient)
	if !ok {
		return nil
	}

	return c
}

//  从 ring 删除 client
func (cr *ClientRing) pop() (*GithubClient, error) {
	client, err := cr.Pop()
	if err != nil {
		return nil, err
	}

	c, ok := client.(*GithubClient)
	if !ok {
		return nil, err
	}

	return c, nil
}

// 判断 ring 中的 client 是否有效，ring 中只保留有效的 client
func (cr *ClientRing) clean() {
	var (
		c   *GithubClient
		err error
	)

	for {
		client := cr.get()

		if client.Remaining == empty {
			c, err = cr.pop()
			if err != nil {
				logger.Error("An error occurred while popping the client: ", err)
			}
		}

		if !c.Limited {
			err = cr.push(c)
			if err != nil {
				logger.Error("An error occurred while pushing the client: ", err)
			}
		}
	}
}

// 生成多个 client 存入 ring 后返回一个 client
func createClient(tokens []string) *GithubClient {
	err := ClientRingService.newClientRing(tokens)
	if err != nil {
		logger.Error("An error occurred while creating the client ring: ", err)
		return nil
	}

	client := ClientRingService.get()

	// TODO： 清理无效的 client
	return client
}

var tokens []string = []string{
	"54f7488c8f72d3e63692b2bf04167d97e7a29e1d",
	"5511599ff7aebf94476ce3eda7741ab7ae797ef9",
	"78d1dcb42b8c4368884603cfcd4f3a1581d771d2",
	"5df193b89001e9fabfdb947a88cdd8b6e45378f5",
}

var GitClient *GithubClient = createClient(tokens)

var ClientRingService *ClientRing = new(ClientRing)

func GenerateClient(useType string, resp *github.Response) {
	for {
		GitClient.monitor(useType, resp)
		ClientRingService.clean()
	}
}
