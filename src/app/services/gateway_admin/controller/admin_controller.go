// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-framework/mvc/storage"
)

type AdminController struct {

}

func NewAdminController() *AdminController {
	return &AdminController {

	}
}

func (c *AdminController) ToSuccessAdmin(ctx *gin.Context) {
	ctx.String(storage.StatusOK, "gateway administration success")
	ctx.Abort()
	return
}
