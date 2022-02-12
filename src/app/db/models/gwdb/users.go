// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gwdb

import (
	"github.com/golang-framework/mvc/src/models"
	"github.com/golang-framework/mvc/storage"
)

type UserModel struct {
	mod *models.Model
}

func NewUserModel() *UserModel {
	return &UserModel {
		mod: models.New(1),
	}
}

func (m *UserModel) AddUser(d *Users) (int64, error) {
	return m.mod.Insert(d)
}

func (m *UserModel) SetUser(d *Users, conditions *storage.Conditions) (int64, error) {
	return m.mod.Update(conditions, d)
}

func (m *UserModel) GetUser(conditions *storage.Conditions) (*Users, error) {
	conditions.Types = storage.SelectOne

	d := &Users{}
	_, errUser := m.mod.Select(conditions, d)

	return d, errUser
}