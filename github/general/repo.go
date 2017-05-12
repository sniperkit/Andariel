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
 *     Initial: 28/04/2017        Jia Chenhui
 */

package general

import (
	"context"
	"math"
	"time"

	"github.com/google/go-github/github"
)

// GetRepoByID 根据库 ID 调用 github API 获取库信息
func GetRepoByID(repoID int) (*github.Repository, *github.Response, error) {
	repo, resp, err := GitClient.Client.Repositories.GetByID(context.Background(), repoID)
	if err != nil {
		if resp == nil {
			return nil, nil, err
		}

		return nil, resp, err
	}

	return repo, resp, nil
}

// GetAllRepos 调用  github API 获取所有库信息
// 参数 opt 指定本次要获取多少个库，将获取的库分多少页，每页包含多少个库（最高 100 个）
// For example:
//
//	opt := &github.RepositoryListAllOptions{
//		ListOptions: github.ListOptions{PerPage: 10},
//	}
func GetAllRepos(opt *github.RepositoryListAllOptions) ([]*github.Repository, *github.Response, error) {
	var (
		allRepos []*github.Repository
		resp     *github.Response
	)

	for {
		repos, resp, err := GitClient.Client.Repositories.ListAll(context.Background(), opt)
		if err != nil {
			return allRepos, resp, err
		}

		allRepos = append(allRepos, repos...)

		if len(repos) == 0 {
			break
		}

		opt.Since = *repos[len(repos)-1].ID
	}

	return allRepos, resp, nil
}

// 按条件从 github 搜索库，受 github API 限制，一次请求只能获取 1000 条记录
func SearchRepos(query string, opt *github.SearchOptions) ([]github.Repository, *github.Response, error) {
	var (
		result []github.Repository
		resp   *github.Response
	)

	page := 1
	maxPage := math.MaxInt32

	for page <= maxPage {
		opt.Page = page
		repos, resp, err := GitClient.Client.Search.Repositories(context.Background(), query, opt)
		Wait(resp)

		if err != nil {
			if resp == nil {
				return nil, nil, err
			}
			return nil, resp, err
		}

		maxPage = resp.LastPage
		result = append(result, repos.Repositories...)

		page++
	}

	return result, resp, nil
}

// SearchReposByCreated 按创建时间及其它指定条件搜索库
// queries: 指定库的创建时间
// For example:
//     queries := []string{"\"2008-06-01 .. 2012-09-01\"", "\"2012-09-02 .. 2013-03-01\"", "\"2013-03-02 .. 2013-09-03\"", "\"2013-09-04 .. 2014-03-05\"", "\"2014-03-06 .. 2014-09-07\"", "\"2014-09-08 .. 2015-03-09\"", "\"2015-03-10 .. 2015-09-11\"", "\"2015-09-12 .. 2016-03-13\"", "\"2016-03-14 .. 2016-09-15\"", "\"2016-09-16 .. 2017-03-17\""}
//
// querySeg: 指定除创建时间之外的其它条件
// For example:
//     queryPart := common.QueryLanguage + ":" + common.LangLua + " " + common.QueryCreated + ":"
//
// opt: 为搜索方法指定可选参数
// For example:
//     opt := &github.SearchOptions{
//         Sort:        common.SortByStars,
//         Order:       common.OrderByDesc,
//         ListOptions: github.ListOptions{PerPage: 100},
//     }
func SearchReposByCreated(queries []string, querySeg string, opt *github.SearchOptions) ([]github.Repository, *github.Response, error) {
	var (
		result []github.Repository
		resp   *github.Response
	)

	for _, q := range queries {
		query := querySeg + q

		repos, resp, err := SearchRepos(query, opt)
		if err != nil {
			return nil, resp, err
		}

		result = append(result, repos...)
	}

	return result, resp, nil
}

// 根据 *github.Response 等待相应时间
func Wait(resp *github.Response) {
	if resp != nil && resp.Remaining <= 1 {
		gap := time.Duration(resp.Reset.Local().Unix() - time.Now().Unix())
		sleep := gap * time.Second

		if sleep < 0 {
			sleep = -sleep
		}

		time.Sleep(sleep)
	}
}
