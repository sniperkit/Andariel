package github

import (
	"context"
	"errors"

	"Andariel/models"
)

type GithubRequestService models.GithubClient

// 根据库 ID 获取库信息并存储到数据库
func (grs *GithubRequestService) GetReposByID(repoID uint64) error {
	var (
		ownerID string
	)

	repos, _, err := grs.Client.Repositories.GetByID(context.Background(), int(repoID))
	if err != nil {
		goto finish
	} else if repos.Fork {
		err = errors.New("this repos is forked from others")

		goto finish
	}

	// 查询数据库中是否有作者信息
	if ok := models.GitUserService.IsUserExists(string(repos.Owner.Name)); ok {
		goto getUserID
	} else {
		goto getOwnerInfo
	}

getUserID:
	ownerID, err = models.GitUserService.GetUserID(string(repos.Owner.Name))
	if err != nil {

		goto finish
	} else {
		goto storeRepos
	}

getOwnerInfo:
	err = grs.GetOwnerByID(int(repos.Owner.ID))
	if err != nil {

		goto finish
	} else {
		goto getUserID
	}

storeRepos:
	err = models.GitReposService.Create(repos, ownerID)

finish:
	return err
}

// 调用 API 获取作者信息并存储到数据库（只在判断数据库中没有作者信息时才调用此函数）
func (grs *GithubRequestService) GetOwnerByID(ownerID int) error {
	owner, _, err := grs.Client.Users.GetByID(context.Background(), ownerID)

	if err != nil {
		return err
	}

	err = models.GitUserService.Create(owner)

	if err != nil {
		return err
	}

	return nil
}
