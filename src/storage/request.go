// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package storage

type (
	Users struct {
		Username	string 	`form:"username" binding:"required"`
		Password	string 	`form:"password" binding:"required"`
	}

	SaveUser struct {
		Id			int		`form:"id"`
		Uid			string	`form:"uid"`
		Password	string	`form:"password"`
		Nickname	string	`form:"nickname"`
		DeleteFlag	int		`form:"delete_flag"`
	}
)
