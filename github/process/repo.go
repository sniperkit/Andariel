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
	"time"

	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"

	git "Andariel/github/general"
	"Andariel/log"
	"Andariel/models"
	"Andariel/utility"
)

var logger *log.AndarielLogger = log.AndarielCreateLogger(
	&log.AndarielLogTag{
		log.LogTagService: "github",
		log.LogTagType:    "process",
	},
	log.AndarielLogLevelDefault)

var ClientManager *git.ClientManager = git.NewClientManager()

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
	client := ClientManager.GetClient()

search:
	repos, resp, stopAt, err := git.SearchReposByStartTime(client, year, month, day, incremental, querySeg, opt)
	if err != nil {
		logger.Error("SearchReposByStartTime returned error:", err)
		return
	}

	// 将获取的库存储到数据库
	for _, repo := range repos {
		err = StoreRepo(&repo, client)
		if err != nil {
			logger.Error("StoreRepo returned error:", err)
			return
		}
	}

	// 判断 client 是否遇到速率限制并将该 client 放回 ClientManager，切换到下个 client 继续执行任务
	if resp != nil && resp.Remaining <= 1 {
		go git.PutClient(client)

		client = ClientManager.GetClient()

		if stopAt != "" {
			newDate, err := util.SplitDate(stopAt)
			if err != nil {
				logger.Error("SplitDate returned error:", err)
			}
			year = newDate[0]
			month = newDate[1]
			switch month {
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
		} else {
			logger.Error("stopAt is empty string")
			return
		}

		goto search
	}
}
