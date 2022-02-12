// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package storage

import (
	err "github.com/golang-framework/mvc/modules/error"
	"github.com/golang-framework/mvc/storage"
)

const (
	KeyGW10001 = "gw_srv_add*user_01"
	KeyGW10002 = "gw_srv_add*user*info_01"
	KeyGW10003 = "gw_srv_add*user*info_02"
	KeyGW10004 = "gw_srv_add*user_04"
	KeyGW10005 = "gw_srv_add*user_05"
	KeyGW10006 = "gw_srv_add*user_06"
	KeyGW10007 = "gw_srv_add*user_07"
	KeyGW10010 = "gw_srv_set*user_01"
	KeyGW10011 = "gw_srv_set*user_02"
	KeyGW10015 = "gw_srv_set*user*info_01"
	KeyGW10016 = "gw_srv_set*user*info_02"
	KeyGW10017 = "gw_srv_login_01"
	KeyGW10018 = "gw_srv_login_02"
	KeyGW10019 = "gw_srv_login_03"
	KeyGW10020 = "gw_srv_login_04"

	valGW10001 = "[10001]添加用户失败"
	valGW10002 = "[10002]生成用户唯一编码失败"
	valGW10003 = "[10003]用户密码加密失败"
	valGW10004 = "[10004]用户名或密码为空"
	valGW10005 = "[10005]用户名必须包含小写字母和数字且长度大于6，小于20"
	valGW10006 = "[10006]密码必须包含大小写字母、数字和特殊字符且长度大于8"
	valGW10007 = "[10007]用户名已存在"
	valGW10010 = "[10010]用户唯一编号不存在"
	valGW10011 = "[10011]用户不存在"
	valGW10015 = "[10015]用户唯一编号错误"
	valGW10016 = "[10016]用户密码转换错误"
	valGW10017 = "[10017]用户不存在"
	valGW10018 = "[10018]用户或密码错误"
	valGW10019 = "[10019]JWT 标签不存在"
	valGW10020 = "[10020]JWT 标签主题不存在"

	KeyGW20001 = "gw_mdw_m_refresh*jwt_01"
	KeyGW21001 = "gw_mdw_is*login_01"
	KeyGW21002 = "gw_mdw_is*login_02"
	KeyGW21003 = "gw_mdw_is*login_03"
	KeyGW21004 = "gw_mdw_is*login_04"
	KeyGW21005 = "gw_mdw_is*login_05"

	valGW20001 = "[20001]JWT 验证码实效"
	valGW21001 = "[21001]JWT Cookie获取失败"
	valGW21002 = "[21002]JWT INF获取失败"
	valGW21003 = "[21003]JWT INF实效"
	valGW21004 = "[21004]JWT 缓存信息获取失败"
	valGW21005 = "[21005]JWT Cookie实效"

	KeyGW30001 = "gw_ctr_add_user_01"
	KeyGW30002 = "gw_ctr_set_user_01"
	KeyGW30003 = "gw_ctr_login_01"

	valGW30001 = "[30001]用户名或密码未填"
	valGW30002 = "[30002]用户需要修改的数据不正确"
	valGW30003 = "[30003]用户名或密码不为空"
)

var EMsg = &storage.E {
	"gw": {
		KeyGW10001: valGW10001,
		KeyGW10002: valGW10002,
		KeyGW10003: valGW10003,
		KeyGW10004: valGW10004,
		KeyGW10005: valGW10005,
		KeyGW10006: valGW10006,
		KeyGW10007: valGW10007,
		KeyGW10010: valGW10010,
		KeyGW10011: valGW10011,
		KeyGW10015: valGW10015,
		KeyGW10016: valGW10016,
		KeyGW10017: valGW10017,
		KeyGW10018: valGW10018,
		KeyGW10019: valGW10019,
		KeyGW10020: valGW10020,

		KeyGW20001: valGW20001,
		KeyGW21001: valGW21001,
		KeyGW21002: valGW21002,
		KeyGW21003: valGW21003,
		KeyGW21004: valGW21004,
		KeyGW21005: valGW21005,

		KeyGW30001: valGW30001,
		KeyGW30002: valGW30002,
		KeyGW30003: valGW30003,
	},
	"gw_administration": {

	},
}

func ErrGW(k string, content ... interface{}) error {
	pfx := "gw"
	p(&pfx, k)

	return err.Err(pfx, k, content ...)
}

func NumGW(k string) string {
	pfx := "gw"
	p(&pfx, k)

	return err.Num(pfx, k)
}

func ErrGWAdmin(k string, content ... interface{}) error {
	pfx := "gw_administration"
	p(&pfx, k)

	return err.Err(pfx, k, content ...)
}

func p(pfx *string, k string) {
	msg := (*EMsg)[*pfx]
	_, ok := msg[k]
	if ok == false {
		*pfx = storage.ErrPrefix
	}
}
