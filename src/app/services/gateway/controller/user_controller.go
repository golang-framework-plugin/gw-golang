// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-framework/mvc/modules/tool"
	"github.com/golang-framework/mvc/modules/uuid"
	"github.com/golang-framework/mvc/storage"
	"src/app/middleware"
	"src/app/services/gateway/service"
	store "src/storage"
)

type UserController struct {
	uid *uuid.M
	tfs *tool.M
	mdw *middleware.M
	srv *service.UserService
}

func NewUserController() *UserController {
	return &UserController {
		uid: uuid.New(),
		tfs: tool.New(),
		mdw: middleware.New(),
		srv: service.NewUserService(),
	}
}

func (c *UserController) Login(ctx *gin.Context) {
	users := &store.Users{}

	if errUsers := ctx.ShouldBind(&users); errUsers != nil {
		c.mdw.Parent.Abort(ctx, store.ErrGW(store.KeyGW30003))
		return
	}

	c.tfs.SourceFilter(
		&users.Username,
		&users.Password,
	)

	ctx.JSON(storage.StatusOK, c.srv.Login(users, ctx))
	ctx.Abort()

	return
}

func (c *UserController) Logout(ctx *gin.Context) {
	ctx.JSON(storage.StatusOK, c.srv.Logout(ctx))
	ctx.Abort()
	return
}

func (c *UserController) Add(ctx *gin.Context) {
	users := &store.Users {}

	if errUsers := ctx.ShouldBind(&users); errUsers != nil {
		c.mdw.Parent.Abort(ctx, store.ErrGW(store.KeyGW30001))
		return
	}

	c.tfs.SourceFilter(
		&users.Username,
		&users.Password,
	)

	ctx.JSON(storage.StatusOK, c.srv.Add(users))
	ctx.Abort()

	return
}

func (c *UserController) Save(ctx *gin.Context) {
	users := &store.SaveUser{}

	if errUserUpdate := ctx.ShouldBind(&users); errUserUpdate != nil {
		c.mdw.Parent.Abort(ctx, store.ErrGW(store.KeyGW30001))
		return
	}

	c.tfs.SourceFilter(
		&users.Uid,
		&users.Nickname,
		&users.Password,
	)

	ctx.JSON(storage.StatusOK, c.srv.Save(users))
	ctx.Abort()

	return
}
