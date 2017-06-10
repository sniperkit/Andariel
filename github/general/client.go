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
 *     Update: 2017/05/30         Jia Chenhui
 */

package general

import (
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"Andariel/log"
)

const (
	unLimit = 1
)

type rateLimitCategory uint8

const (
	coreCategory rateLimitCategory = iota
	searchCategory
	categories
)

var logger *log.AndarielLogger = log.AndarielCreateLogger(
	&log.AndarielLogTag{
		log.LogTagService: "github",
		log.LogTagType:    "client",
	},
	log.AndarielLogLevelDefault)

type GHClient struct {
	Client     *github.Client
	rateLimits [categories]*Rate
	manager    *ClientManager
	timer      *time.Timer
}

type Rate struct {
	Limit     int
	Remaining int
	Reset     time.Time
	Limited   bool
}

// 新建 client
func newClient(token string) (client *GHClient) {
	if token == "" {
		client = new(GHClient)
		tokenSource := new(oauth2.TokenSource)
		client.init(*tokenSource)
		return
	}

	client = new(GHClient)
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client.init(tokenSource)

	return
}

// 初始化 client
func (c *GHClient) init(tokenSource oauth2.TokenSource) {
	httpClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	// 检查 token 是否有效
	if !c.isValidToken(httpClient) {
		logger.Debug("Invalid token.")
		return
	}

	ghClient := github.NewClient(httpClient)
	c.Client = ghClient

	if c.isLimited() {
		logger.Debug("Hit rate limit while initializing the client.")
	}
}

// 检查 token 是否有效
func (c *GHClient) isValidToken(httpClient *http.Client) bool {
	resp, err := c.makeRequest(httpClient)
	if err != nil {
		logger.Error("makeRequest returned error:", err)
		return false
	}

	err = github.CheckResponse(resp)
	if e, ok := err.(*github.TwoFactorAuthError); ok {
		logger.Error("401 Unauthorized:", e)
		return false
	}

	return true
}

// 发起一次 API 请求
func (c *GHClient) makeRequest(httpClient *http.Client) (*http.Response, error) {
	req, err := c.Client.NewRequest("GET", c.Client.BaseURL.Path, nil)
	if err != nil {
		logger.Error("Client.NewRequest returned error:", err)
		return nil, err
	}

	// 发起请求
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("httpClient.Do returned error:", err)
		return nil, err
	}

	return resp, nil
}

// 检查 rateLimit
func (c *GHClient) isLimited() bool {
	rate, _, err := c.Client.RateLimits(oauth2.NoContext)
	if err != nil {
		logger.Error("Client.RateLimits returned error:", err)
		return true
	}

	if rate != nil {
		if rate.Core != nil {
			c.rateLimits[coreCategory].Limit = rate.Core.Limit
			c.rateLimits[coreCategory].Remaining = rate.Core.Remaining
			c.rateLimits[coreCategory].Reset = rate.Core.Reset.Time
			if rate.Core.Remaining <= unLimit {
				c.rateLimits[coreCategory].Limited = true
				return true
			}
			return false
		}
		if rate.Search != nil {
			c.rateLimits[searchCategory].Limit = rate.Search.Limit
			c.rateLimits[searchCategory].Remaining = rate.Search.Remaining
			c.rateLimits[searchCategory].Reset = rate.Search.Reset.Time
			if rate.Search.Remaining <= unLimit {
				c.rateLimits[searchCategory].Limited = true
				return true
			}
			return false
		}
	}

	return true
}

// 初始化 client 的 timer
func (c *GHClient) initTimer() {
	if c.rateLimits[coreCategory] != nil {
		coreTimer := time.NewTimer(c.rateLimits[coreCategory].Reset.Sub(time.Now()))
		c.timer = coreTimer
		return
	}
	if c.rateLimits[searchCategory] != nil {
		searchTimer := time.NewTimer(c.rateLimits[searchCategory].Reset.Sub(time.Now()))
		c.timer = searchTimer
		return
	}
}

// 新建多个 client
func newClients(tokens []string) []*GHClient {
	var clients []*GHClient

	for _, t := range tokens {
		client := newClient(t)
		clients = append(clients, client)
	}

	return clients
}

var tokens []string = []string{
	"54f7488c8f72d3e63692b2bf04167d97e7a29e1d",
	"5511599ff7aebf94476ce3eda7741ab7ae797ef9",
	"78d1dcb42b8c4368884603cfcd4f3a1581d771d2",
	"5df193b89001e9fabfdb947a88cdd8b6e45378f5",
}

type ClientManager struct {
	inputChan  chan *GHClient
	OutputChan chan *GHClient
}

func (r *ClientManager) Run(done chan bool) {
	for {
		select {
		case v := <-r.inputChan:
			r.OutputChan <- v
		case <-done:
			break
		}
	}
	close(r.OutputChan)
}

// NewClientManager 创建新的 ClientManager
func NewClientManager() *ClientManager {
	var rb *ClientManager = &ClientManager{
		inputChan:  make(chan *GHClient),
		OutputChan: make(chan *GHClient, len(tokens)),
	}

	clients := newClients(tokens)
	done := make(chan bool)
	defer func() {
		done <- true
	}()

	go rb.Run(done)
	go func() {
		for _, c := range clients {
			if !c.isLimited() {
				rb.inputChan <- c
			}
		}
	}()

	return rb
}

// GetClient 读取 client
func (m *ClientManager) GetClient() *GHClient {
	select {
	case c := <-m.OutputChan:
		return c
	default:
		return nil
	}
}

// PutClient 将 client 放回 manager
func PutClient(client *GHClient) {
	client.initTimer()
	<-client.timer.C

	select {
	case client.manager.inputChan <- client:
	}
}
