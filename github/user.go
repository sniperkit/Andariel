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

	"Andariel/models"
)

// 调用 API 获取作者信息
func GetOwnerByID(ownerID int) (*models.User, error) {
	owner, _, err := GitClient.Client.Users.GetByID(context.Background(), ownerID)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:                uint64(owner.ID),
		HTMLURL:           string(owner.HTMLURL),
		Name:              string(owner.Name),
		Email:             string(owner.Email),
		PublicRepos:       uint64(owner.PublicRepos),
		PublicGists:       uint64(owner.PublicGists),
		Followers:         uint64(owner.Followers),
		Following:         uint64(owner.Following),
		CreatedAt:         owner.CreatedAt,
		UpdatedAt:         owner.UpdatedAt,
		SuspendedAt:       owner.SuspendedAt,
		Type:              string(owner.Type),
		TotalPrivateRepos: uint64(owner.TotalPrivateRepos),
		OwnedPrivateRepos: uint64(owner.OwnedPrivateRepos),
		PrivateGists:      uint64(owner.PrivateGists),
	}

	return user, nil
}
