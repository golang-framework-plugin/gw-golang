// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package demo

import (
	"github.com/golang-framework/mvc/src/models"
	"github.com/golang-framework/mvc/storage"
)

type DemoModel struct {
	mod *models.Model
}

func NewDemoModel() *DemoModel {
	return &DemoModel {
		mod: models.New(0),
	}
}

func (m *DemoModel) FetchOne(conditions *storage.Conditions) (*Demo, error) {
	conditions.Types = storage.SelectOne

	var d = &Demo{}
	_, e := m.mod.Select(conditions, d)

	return d, e
}

func (m *DemoModel) FetchAll(conditions *storage.Conditions) ([]Demo, error) {
	conditions.Types = storage.SelectAll

	var d = make([]Demo, 0)
	_, e := m.mod.Select(conditions, &d)

	return d, e
}
