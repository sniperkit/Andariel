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

	"Andariel/common/container"
	"Andariel/log"
)

const (
	empty   = 0
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
	rateLimits [categories]Rate
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
		client.init(tokenSource)
		return
	}

	client = new(GHClient)
	tokenSource := oauth2.StaticTokenSource(token)
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
		logger.Debug("Hit rate limit.")
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
	req, err := c.Client.NewRequest("GET", c.Client.BaseURL, nil)
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
			c.rateLimits[coreCategory] = *rate.Core
			if rate.Core.Remaining <= unLimit {
				c.rateLimits[coreCategory].Limited = true
				return true
			}
			return false
		}
		if rate.Search != nil {
			c.rateLimits[searchCategory] = *rate.Search
			if rate.Search.Remaining <= unLimit {
				c.rateLimits[searchCategory].Limited = true
				return true
			}
			return false
		}
	}

	return true
}

// 重置 client
func (c *GHClient) reset() {
	if c.rateLimits[coreCategory] != nil {
		coreTimer := time.NewTimer(c.rateLimits[coreCategory].Reset.Sub(time.Now()))
		<-coreTimer.C
		c.rateLimits[coreCategory].Remaining = unLimit
		c.rateLimits[coreCategory].Limited = false
	}
	if c.rateLimits[searchCategory] != nil {
		searchTimer := time.NewTimer(c.rateLimits[searchCategory].Reset.Sub(time.Now()))
		<-searchTimer.C
		c.rateLimits[searchTimer].Remaining = unLimit
		c.rateLimits[searchTimer].Limited = false
	}
}

// 新建多个 client
func NewClients(tokens []string) []*GHClient {
	var clients []*GHClient

	for _, t := range tokens {
		client := newClient(t)
		clients = append(clients, client)
	}

	return clients
}

// 监控 client，将可用的 client 放入 ring
func (c *GHClient) Monitor() {
	for {
		switch {
		case !c.isLimited():
			addClient(c)
		case c.isLimited():
			popClient()
		}
	}
}

var (
	ringCap = 50
	cRing   *container.Ring
)

// 存储 client
func addClient(client interface{}) {
	if cRing == nil {
		cRing = container.NewRing(ringCap)
	}

	client = client.(*GHClient)

	if !cRing.Full() {
		if err := cRing.Push(client); err != nil {
			logger.Error("Ring.Push returned error:", err)
		}
	}
}

// 获取 client
func getClient() interface{} {
	return cRing.Get()
}

// 删除 client
func popClient() (interface{}, error) {
	return cRing.Pop()
}

// 初始化 clientRing
func InitCRing(tokens []string) {
	clients := NewClients(tokens)
	for _, c := range clients {
		addClient(c)
	}
}
