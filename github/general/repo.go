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

	"Andariel/common"
	"Andariel/utility"
)

// GetRepoByID 根据库 ID 调用 github API 获取库信息
func GetRepoByID(client *GHClient, repoID int) (*github.Repository, *github.Response, error) {
	repo, resp, err := client.Client.Repositories.GetByID(context.Background(), repoID)
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
func GetAllRepos(client *GHClient, opt *github.RepositoryListAllOptions) ([]*github.Repository, *github.Response, error) {
	var (
		allRepos []*github.Repository
		resp     *github.Response
	)

	for {
		repos, resp, err := client.Client.Repositories.ListAll(context.Background(), opt)
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

// SearchRepos 按条件从 github 搜索库，受 github API 限制，一次请求只能获取 1000 条记录
// GitHub API docs: https://developer.github.com/v3/search/#search-repositories
func searchRepos(client *GHClient, query string, opt *github.SearchOptions) ([]github.Repository, *github.Response, string, error) {
	var (
		result []github.Repository
		repos  *github.RepositoriesSearchResult
		resp   *github.Response
		stopAt string
		err    error
	)

	page := 1
	maxPage := math.MaxInt32

	for page <= maxPage {
		opt.Page = page

		repos, resp, err = client.Client.Search.Repositories(context.Background(), query, opt)
		if err != nil {
			goto finish
		}

		maxPage = resp.LastPage
		result = append(result, repos.Repositories...)

		page++
	}

finish:
	if len(result) != 0 {
		stopAt = util.SplitQuery(query)
	} else {
		stopAt = ""
	}

	return result, resp, stopAt, err
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
// GitHub API docs: https://developer.github.com/v3/search/#search-repositories
func SearchReposByCreated(client *GHClient, queries []string, querySeg string, opt *github.SearchOptions) ([]github.Repository, *github.Response, string, error) {
	var (
		result, repos []github.Repository
		resp          *github.Response
		stopAt        string
		err           error
	)

	for _, q := range queries {
		query := querySeg + q

		repos, resp, stopAt, err = searchRepos(client, query, opt)
		if err != nil {
			goto finish
		}

		result = append(result, repos...)
	}

finish:
	return result, resp, stopAt, nil
}

// SearchReposByStartTime 按指定创建时间、时间间隔及其它条件搜索库
// year、month、day: 从此创建时间开始搜索
// For example：
//     year = 2016 month = time.January day = 1
//     时间格式化只能使用 "2006-01-02 15:04:05" 进行，可将年月日和 时分秒拆开使用
//
// incremental: 以此时间增量搜索，如第一次搜索 1 月份的库，第二次搜索 2 月份的库
// For example:
//     interval = "month"
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
// GitHub API docs: https://developer.github.com/v3/search/#search-repositories
func SearchReposByStartTime(client *GHClient, year int, month time.Month, day int, incremental, querySeg string, opt *github.SearchOptions) ([]github.Repository, *github.Response, string, error) {
	var (
		result, repos []github.Repository
		resp          *github.Response
		stopAt        string
		err           error
	)

	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	for date.Unix() < time.Now().Unix() {
		var dateFormat string

		switch incremental {
		case common.OneQuarter:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 3, 0).Format("2006-01-02")
		case common.OneMonth:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 1, 0).Format("2006-01-02")
		case common.OneWeek:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 0, 7).Format("2006-01-02")
		case common.OneDay:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 0, 1).Format("2006-01-02")
		default:
			dateFormat = date.Format("2006-01-02") + " .. " + date.AddDate(0, 1, 0).Format("2006-01-02")
		}

		query := querySeg + dateFormat

		repos, resp, stopAt, err = searchRepos(client, query, opt)
		if err != nil {
			goto finish
		}

		result = append(result, repos...)

		switch incremental {
		case common.OneQuarter:
			date = date.AddDate(0, 3, 1)
		case common.OneMonth:
			date = date.AddDate(0, 1, 1)
		case common.OneWeek:
			date = date.AddDate(0, 0, 8)
		case common.OneDay:
			date = date.AddDate(0, 0, 2)
		default:
			date = date.AddDate(0, 1, 1)
		}
	}

finish:
	return result, resp, stopAt, err
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
