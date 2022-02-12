// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-framework/mvc/storage"
)

type ApiController struct {

}

func NewApiController() *ApiController {
	return &ApiController {

	}
}

func (c *ApiController) T(ctx *gin.Context) {

	ctx.JSON(storage.StatusOK, "gw_api_success")
	ctx.Abort()
	return
}
