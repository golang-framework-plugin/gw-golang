// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (_ *gwAdmin) GwAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("test gw admin before!")
		ctx.Next()
	}
}
