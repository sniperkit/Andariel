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
