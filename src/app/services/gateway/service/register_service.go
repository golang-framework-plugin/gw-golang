// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	err "github.com/golang-framework/mvc/modules/error"
	"github.com/golang-framework/mvc/storage"
	"src/app/db/models/demo"
	store "src/storage"
)

type RegisterService struct {
	modDemo *demo.DemoModel
}

func NewRegisterService() *RegisterService {
	return &RegisterService {
		modDemo: demo.NewDemoModel(),
	}
}

func (s *RegisterService) GetDemo() *storage.Tpl {
	res := storage.FwTpl(store.ErrGW("test"))

	conditions := &storage.Conditions {}

	d, e := s.modDemo.FetchOne(conditions)
	if e != nil {
		res.Res = &storage.Y {}
		return res
	}

	res.Res = &storage.Y {
		"demo": d,
	}

	return res
}

func (s *RegisterService) GetDemos() *storage.Tpl {
	res := storage.FwTpl(err.Err("SRV_GW", "test"))

	conditions := &storage.Conditions {}

	d, e := s.modDemo.FetchAll(conditions)
	if e != nil {
		res.Status = -1
		res.Msg = ""
		res.Res = &storage.Y {}
	}

	res.Res = &storage.Y {
		"demo": d,
	}

	return res
}
