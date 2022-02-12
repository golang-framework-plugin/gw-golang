// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gwdb

import (
	"time"
)

type Users struct {
	Id         	int			`xorm:"not null pk autoincr INT(11)" json:"id,omitempty"`
	Uid       	string  	`xorm:"unique 'uid' VARCHAR(50)" json:"uuid,omitempty"`
	Username   	string  	`xorm:"unique 'username' VARCHAR(50)" json:"username,omitempty"`
	Password   	string		`xorm:"VARCHAR(32)" json:"password,omitempty"`
	Nickname   	string  	`xorm:"VARCHAR(50)" json:"nickname,omitempty"`
	CreateTime 	int 		`xorm:"created INT(11)" json:"create_time,omitempty"`
	UpdateTime 	*time.Time 	`xorm:"updated TIMESTAMP" json:"update_time,omitempty"`
	DeleteFlag 	int	    	`xorm:"default 1 TINYINT(1)" json:"delete_flag,omitempty"`
}
