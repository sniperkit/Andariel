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
	"errors"

	"Andariel/models"
	"gopkg.in/mgo.v2"
)

// 根据库 ID 获取库信息并存储到数据库
func GetRepoByID(repoID uint64) error {
	// 调用官方 API 获取库信息
	repo, _, err := GitClient.Client.Repositories.GetByID(context.Background(), int(repoID))
	if err != nil {

		return err
	}

	// fork 的库不保存
	if repo.Fork {
		err = errors.New("this repos is forked from others")

		return err
	}

	// 判断数据库中是否有此作者信息
	oldUserID, err := models.GitUserService.GetUserID(string(repo.Owner.Name))
	if err != nil {
		if err != mgo.ErrNotFound {

			return err
		}

		// User 数据库中无此作者信息
		newUserID, err := models.GitUserService.Create(repo.Owner)
		if err != nil {

			return err
		}

		err = models.GitReposService.Create(repo, &newUserID)
		if err != nil {

			return err
		}
	}

	// User 数据库中有此作者信息
	err = models.GitReposService.Create(repo, &oldUserID)
	if err != nil {

		return err
	}

	return nil
}
