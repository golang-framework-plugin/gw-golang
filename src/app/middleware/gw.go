// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/gin-gonic/gin"
)

func (_ *gw) Gw() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
