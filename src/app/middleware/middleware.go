// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-framework/mvc/modules/property"
	"github.com/golang-framework/mvc/src/components/caches/redis"
	"github.com/golang-framework/mvc/src/components/jwt"
	"github.com/golang-framework/mvc/src/middleware"
	"github.com/golang-framework/mvc/storage"
	"github.com/spf13/cast"
	store "src/storage"
)

var (
	InstanceGw      = new(gw)
	InstanceGwAdmin = new(gwAdmin)
)

type (
	gw      struct{}
	gwAdmin struct{}

	M struct {
		Parent *middleware.M
	}
)

func New() *M {
	return &M{
		Parent: middleware.New(),
	}
}

func (m *M) IsLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//toolTP := tool.New()

		pyJTS := property.Instance.Nest("JWT.%v.SUB", "JwTTag.GwSrv001")
		if pyJTS == nil || pyJTS.(string) == "" {
			m.Parent.Abort(ctx, store.ErrGW(store.KeyGW10020))
			return
		}

		cookieName, errCookieName := ctx.Cookie(pyJTS.(string))
		if errCookieName != nil {
			m.Parent.Abort(ctx, store.ErrGW(store.KeyGW21001))
			return
		}

		componentJwT := jwt.NewJwT("JwTTag.GwSrv001")
		_, payload, errJwTParse := componentJwT.Parse(cookieName)
		if errJwTParse != nil {
			m.Parent.Abort(ctx, store.ErrGW(store.KeyGW21002))
			return
		}

		userID, ok := payload.Inf.(map[string]interface{})["ID"]
		if ok == false {
			m.Parent.Abort(ctx, store.ErrGW(store.KeyGW21003))
			return
		}

		i := 0
		r := redis.New(i)
		r.SetPrefix(cast.ToString(
			property.Instance.Get("RedisPrefixKey.Gw", "").(interface{}),
		))
		rJwT, errRedisJwT := r.HGet(pyJTS.(string), cast.ToString(userID))
		if errRedisJwT != nil {
			m.Parent.Abort(ctx, store.ErrGW(store.KeyGW21004))
			return
		}

		if cookieName != rJwT {
			m.Parent.Abort(ctx, store.ErrGW(store.KeyGW21005))
			return
		}

		componentJwT.Inf = payload.Inf
		componentJwT.Exp = cast.ToDuration(30 * time.Second)

		c, errJwTRefresh := componentJwT.Refresh(cookieName)
		if errJwTRefresh != nil {
			m.Parent.Abort(ctx, errJwTRefresh)
			return
		}

		_, errRedisHSet := r.HSet(pyJTS.(string), map[string]interface{}{
			cast.ToString(userID): c.(string),
		})
		if errRedisHSet != nil {
			m.Parent.Abort(ctx, errRedisHSet)
			return
		}

		m.SetCookie(ctx, &storage.TplCookie{
			Name:     pyJTS.(string),
			Value:    c.(string),
			MaxAge:   30,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
			HttpOnly: false,
		})

		ctx.Next()
	}
}

func (_ *M) SetCookie(ctx *gin.Context, cookie *storage.TplCookie) {
	ctx.SetCookie(
		cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path,
		cookie.Domain, cookie.Secure, cookie.HttpOnly,
	)
}

func NoRouter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(storage.StatusOK, "no router for golang-fw")
		ctx.Abort()
		return
	}
}

func NoMethod() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(storage.StatusOK, "no method for golang-fw")
		ctx.Abort()
		return
	}
}
