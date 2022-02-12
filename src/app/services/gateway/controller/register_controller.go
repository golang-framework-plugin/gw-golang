// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-framework/mvc/modules/uuid"
	"github.com/golang-framework/mvc/src/components/caches/redis"
	"github.com/golang-framework/mvc/storage"
	"github.com/spf13/cast"
	"src/app/services/gateway/service"
	"strings"
)

type RegisterController struct {
	r *redis.Component
	srv *service.RegisterService
	d *uuid.M
}

func NewRegisterController() *RegisterController {
	return &RegisterController {
		r: redis.New(0),
		srv: service.NewRegisterService(),
		d: uuid.New(),
	}
}

func (c *RegisterController) To(ctx *gin.Context) {
	fmt.Println("To !!!")
	ctx.Abort()
	return
}

func (c *RegisterController) ToSuccess(ctx *gin.Context) {
	e := " asdfasd "
	f := " !dafda   "
	g := "lasdjf "
	h := " askjldfa"
	var a []*string = []*string{ &e,&f,&g,&h }
	fmt.Println(*a[0],*a[1],*a[2],*a[3])
	t(a ...)
	fmt.Println(*a[0],*a[1],*a[2],*a[3])

	ctx.JSON(storage.StatusOK, c.srv.GetDemo())
	ctx.Abort()
	return
}

func t(a ... *string) {

	for _, v := range a {
		*v = strings.Trim(cast.ToString(*v), " ")
	}
}
