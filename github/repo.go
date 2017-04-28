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

package github

import (
	"context"

	"github.com/google/go-github/github"

	"Andariel/models"
)

// 根据库 ID 获取库信息并存储到数据库
func GetRepoByID(repoID int) error {
	// 调用官方 API 获取库信息
	repo, _, err := GitClient.Client.Repositories.GetByID(context.Background(), repoID)
	if err != nil {

		return err
	}

	err = models.StoreRepo(repo)
	if err != nil {
		return err
	}

	return nil
}

// 调用 github API 获取所有库信息并存储到数据库
func GetAllRepos(opt *github.RepositoryListAllOptions) error {
	repos, _, err := GitClient.Client.Repositories.ListAll(context.Background(), opt)
	if err != nil {

		return err
	}

	for _, repo := range repos {
		err = models.StoreRepo(repo)
		if err != nil {

			return err
		}
	}

	return nil
}
