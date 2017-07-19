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
 *     Initial: 08/05/2017        Jia Chenhui
 */

package process

import (
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"

	"Andariel/models"
	git "Andariel/pkg/github"
	"Andariel/pkg/log"
	"Andariel/pkg/utility"
)

var clientManager *git.ClientManager = git.NewClientManager()

// 逻辑判断后，存储库信息到数据库
func StoreRepo(repo *github.Repository, client *git.GHClient) error {
	// 判断数据库中是否有此作者信息
	oldUserID, err := models.GitUserService.GetUserID(repo.Owner.Login)
	if err != nil {
		if err != mgo.ErrNotFound {
			return err
		}

		// MDUser 数据库中无此作者信息
		newOwner, _, err := git.GetOwnerByID(*repo.Owner.ID, client)
		if err != nil {
			return err
		}

		newUserID, err := models.GitUserService.Create(newOwner)
		if err != nil {
			return err
		}

		err = models.GitReposService.Create(repo, &newUserID)
		if err != nil {
			return err
		}
	} else {
		// MDUser 数据库中有此作者信息
		err = models.GitReposService.Create(repo, &oldUserID)
		if err != nil {
			return err
		}
	}

	return nil
}

// SearchRepos 从指定时间（库的创建时间）开始搜索，并将结果保存到数据库
func SearchRepos(year int, month time.Month, day int, incremental, querySeg string, opt *github.SearchOptions) {
	var (
		client  *git.GHClient
		ok      bool
		wg      sync.WaitGroup
		e       *github.AbuseRateLimitError
		newDate []int
		result  []github.Repository
	)

	client = clientManager.GetClient()
	client.Manager = clientManager

search:
	repos, resp, stopAt, err := git.SearchReposByStartTime(client, year, month, day, incremental, querySeg, opt)
	result = append(result, repos...)

	if err != nil {
		if _, ok = err.(*github.RateLimitError); ok {
			log.Logger.Error("SearchReposByStartTime hit limit error, it's time to change client.", zap.Error(err))

			goto changeClient
		} else if e, ok = err.(*github.AbuseRateLimitError); ok {
			log.Logger.Error("SearchReposByStartTime have triggered an abuse detection mechanism.", zap.Error(err))

			time.Sleep(*e.RetryAfter)
			goto search
		} else if strings.Contains(err.Error(), "timeout") {
			log.Logger.Info("SearchReposByStartTime has encountered a timeout error. Sleep for five minutes.")
			time.Sleep(5 * time.Minute)

			goto search
		} else {
			log.Logger.Error("SearchRepos terminated because of this error.", zap.Error(err))

			return
		}
	} else {

		goto store
	}

changeClient:
	{
		go func() {
			wg.Add(1)
			defer wg.Done()

			git.PutClient(client, resp)
		}()

		client = clientManager.GetClient()
		client.Manager = clientManager

		if stopAt != "" {
			newDate, err = utility.SplitDate(stopAt)
			if err != nil {
				log.Logger.Error("SplitDate returned error.", zap.Error(err))

				return
			}

			year = newDate[0]
			monthInt := newDate[1]
			switch monthInt {
			case 1:
				month = time.January
			case 2:
				month = time.February
			case 3:
				month = time.March
			case 4:
				month = time.April
			case 5:
				month = time.May
			case 6:
				month = time.June
			case 7:
				month = time.July
			case 8:
				month = time.August
			case 9:
				month = time.September
			case 10:
				month = time.October
			case 11:
				month = time.November
			case 12:
				month = time.December
			}
			day = newDate[2]

			goto search
		}

		log.Logger.Info("stopAt is empty string, stop searching.")
	}

store:
	log.Logger.Info("Start storing repositories now.")
	for _, repo := range result {
	repeatStore:
		err = StoreRepo(&repo, client)
		if err != nil {
			if _, ok = err.(*github.RateLimitError); ok {
				log.Logger.Error("StoreRepo hit limit error, it's time to change client.", zap.Error(err))

				go func() {
					wg.Add(1)
					defer wg.Done()

					git.PutClient(client, resp)
				}()

				client = clientManager.GetClient()
				client.Manager = clientManager

				goto repeatStore
			} else if e, ok = err.(*github.AbuseRateLimitError); ok {
				log.Logger.Error("SearchReposByStartTime have triggered an abuse detection mechanism.", zap.Error(err))

				time.Sleep(*e.RetryAfter)
				goto repeatStore
			} else {
				log.Logger.Error("StoreRepo encounter this error, proceed to the next loop.", zap.Error(err))

				continue
			}
		}
	}

	wg.Wait()
	log.Logger.Info("All search and storage tasks have been successful.")
}
