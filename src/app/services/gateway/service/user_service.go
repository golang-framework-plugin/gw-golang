// Copyright 2022 YuWenYu  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-framework/mvc/modules/crypto"
	"github.com/golang-framework/mvc/modules/property"
	"github.com/golang-framework/mvc/modules/tool"
	"github.com/golang-framework/mvc/modules/uuid"
	"github.com/golang-framework/mvc/src/components/caches/redis"
	"github.com/golang-framework/mvc/src/components/jwt"
	"github.com/golang-framework/mvc/storage"
	"github.com/spf13/cast"
	"src/app/db/models/gwdb"
	"src/app/middleware"
	store "src/storage"
	"strings"
	"time"
)

type UserService struct {
	toolsMP *tool.M
	toolsUd *uuid.M
	modUser *gwdb.UserModel

	redisId int
}

func NewUserService() *UserService {
	return &UserService{
		toolsMP: tool.New(),
		toolsUd: uuid.New(),
		modUser: gwdb.NewUserModel(),

		redisId: 0,
	}
}

func (s *UserService) Login(users *store.Users, ctx *gin.Context) *storage.Tpl {
	res := storage.FwTpl(store.ErrGW(storage.SuccessOK))

	if users.Username == "" || users.Password == "" {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10004).Error()

		return res
	}

	if s.toolsMP.MatchPattern(users.Username, storage.PatternType01, 6, 20) == -1 {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10005).Error()

		return res
	}

	if s.toolsMP.MatchPattern(users.Password, storage.PatternType02, 8) == -1 {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10006).Error()

		return res
	}

	if s.isUserExistByUsername(users.Username) == -1 {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10017).Error()

		return res
	}

	pwd, errPwd := s.generatePwd(users.Password)
	if errPwd != nil {
		res.Status = storage.StatusUnknown
		res.Msg = errPwd.Error()

		return res
	}

	conditions := &storage.Conditions{
		Query:     "username=? AND password=?",
		QueryArgs: []interface{}{users.Username, pwd},
		Columns:   []string{"id"},
	}

	d, errUserLogin := s.getUser(conditions)
	if errUserLogin != nil || d.Id < 1 {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10018).Error()

	} else {
		pyJTS := property.Instance.Nest("JWT.%v.SUB", "JwTTag.GwSrv001")
		if pyJTS == nil || pyJTS.(string) == "" {
			res.Status = storage.StatusUnknown
			res.Msg = store.ErrGW(store.KeyGW10020).Error()

			return res
		}

		componentJwT := jwt.NewJwT("JwTTag.GwSrv001")
		componentJwT.Inf = &storage.Y{
			"ID":       d.Id,
			"USERNAME": users.Username,
		}
		componentJwT.Exp = cast.ToDuration(30 * time.Second)

		c, errJwTProduce := componentJwT.Produce()
		if errJwTProduce != nil {
			res.Status = storage.StatusUnknown
			res.Msg = errJwTProduce

			return res
		}

		r := redis.New(s.redisId).SetPrefix(cast.ToString(
			property.Instance.Get("RedisPrefixKey.Gw", "").(interface{}),
		))
		_, errRedisHSet := r.HSet(pyJTS.(string), map[string]interface{}{
			cast.ToString(d.Id): c.(string),
		})

		if errRedisHSet != nil {
			res.Status = storage.StatusUnknown
			res.Msg = errRedisHSet

			return res
		}

		(&middleware.M{}).SetCookie(ctx, &storage.TplCookie{
			Name:     pyJTS.(string),
			Value:    c.(string),
			MaxAge:   30,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
			HttpOnly: false,
		})
	}

	return res
}

func (s *UserService) Logout(ctx *gin.Context) *storage.Tpl {
	res := storage.FwTpl(store.ErrGW(storage.SuccessOK))

	pyJTS := property.Instance.Nest("JWT.%v.SUB", "JwTTag.GwSrv001")
	if pyJTS == nil || pyJTS.(string) == "" {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10020).Error()

		return res
	}

	cookieName, errCookieName := ctx.Cookie(pyJTS.(string))
	if errCookieName != nil {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW21001).Error()

		return res
	}

	componentJwT := jwt.NewJwT("JwTTag.GwSrv001")
	_, payload, errJwTParse := componentJwT.Parse(cookieName)
	if errJwTParse != nil {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW21002).Error()

		return res
	}

	userID, ok := payload.Inf.(map[string]interface{})["ID"]
	if ok == false {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW21003).Error()

		return res
	}

	r := redis.New(s.redisId).SetPrefix(cast.ToString(
		property.Instance.Get("RedisPrefixKey.Gw", "").(interface{}),
	))
	_, errRedisHExist := r.HExist(pyJTS.(string), cast.ToString(userID))
	if errRedisHExist != nil {
		res.Status = storage.StatusUnknown
		res.Msg = errRedisHExist

		return res
	}

	_, errRedisHDel := r.HDel(pyJTS.(string), cast.ToString(userID))
	if errRedisHDel != nil {
		res.Status = storage.StatusUnknown
		res.Msg = errRedisHDel

		return res
	}

	(&middleware.M{}).SetCookie(ctx, &storage.TplCookie{
		Name:     pyJTS.(string),
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "127.0.0.1",
		Secure:   true,
		HttpOnly: false,
	})

	return res
}

func (s *UserService) Add(users *store.Users) *storage.Tpl {
	res := storage.FwTpl(store.ErrGW(storage.SuccessOK))

	if users.Username == "" || users.Password == "" {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10004).Error()

		return res
	}

	if s.toolsMP.MatchPattern(users.Username, storage.PatternType01, 6, 20) == -1 {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10005).Error()

		return res
	}

	if s.isUserExistByUsername(users.Username) == 1 {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10007).Error()

		return res
	}

	if s.toolsMP.MatchPattern(users.Password, storage.PatternType02, 8) == -1 {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10006).Error()

		return res
	}

	d, errUserInfo := s.addInfo(users)
	if errUserInfo != nil {
		res.Status = storage.StatusUnknown
		res.Msg = errUserInfo.Error()

		return res
	}

	_, errUser := s.modUser.AddUser(d)
	if errUser != nil {
		res.Status = storage.StatusUnknown
		res.Msg = errUser.Error()

		return res
	}

	return res
}

func (s *UserService) Save(users *store.SaveUser) *storage.Tpl {
	res := storage.FwTpl(store.ErrGW(storage.SuccessOK))

	if users.Id <= 0 && users.Uid == "" {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10010).Error()

		return res
	}

	if (users.Id > 0 && s.isUserExistById(users.Id) == -1) ||
		(users.Uid != "" && s.isUserExistByUid(users.Uid) == -1) {
		res.Status = storage.StatusUnknown
		res.Msg = store.ErrGW(store.KeyGW10011).Error()

		return res
	}

	d, conditions, errUserInfo := s.setInfo(users)
	if errUserInfo != nil {
		res.Status = storage.StatusUnknown
		res.Msg = errUserInfo.Error()

		return res
	}

	_, errUser := s.modUser.SetUser(d, conditions)
	if errUser != nil {
		res.Status = storage.StatusUnknown
		res.Msg = errUser.Error()

		return res
	}

	return res
}

func (s *UserService) addInfo(users *store.Users) (*gwdb.Users, error) {
	s.toolsUd.D = []interface{}{
		storage.Md5,
		cast.ToString(time.Now().UnixNano()),
	}

	userUUID, errUserUUID := s.toolsUd.Generate()
	if errUserUUID != nil {
		return nil, store.ErrGW(store.KeyGW10002)
	}

	pwd, errPwd := s.generatePwd(strings.Trim(users.Password, " "))
	if errPwd != nil {
		return nil, store.ErrGW(store.KeyGW10003)
	}

	return &gwdb.Users{
		Uid:      userUUID.(string),
		Username: users.Username,
		Password: pwd.(string),
	}, nil
}

func (s *UserService) setInfo(users *store.SaveUser) (*gwdb.Users, *storage.Conditions, error) {
	if users.Id == 0 && users.Uid == "" {
		return nil, nil, store.ErrGW(store.KeyGW10015)
	}

	conditions := &storage.Conditions{}
	if users.Id > 0 {
		conditions.Query = "id=?"
		conditions.QueryArgs = []interface{}{users.Id}
	}

	if users.Uid != "" {
		conditions.Query = "uid=?"
		conditions.QueryArgs = []interface{}{users.Uid}
	}

	dbUsers := &gwdb.Users{}

	if users.Nickname != "" {
		dbUsers.Nickname = users.Nickname
	}

	if users.Password != "" {
		pwd, errPwd := s.generatePwd(users.Password)
		if errPwd != nil {
			return nil, nil, store.ErrGW(store.KeyGW10016)
		}

		dbUsers.Password = pwd.(string)
	}

	if users.DeleteFlag == 0 || users.DeleteFlag == 1 {
		dbUsers.DeleteFlag = users.DeleteFlag
	}

	return dbUsers, conditions, nil
}

func (s *UserService) isUserExistById(id int) int8 {
	conditions := &storage.Conditions{
		Query:     "id=?",
		QueryArgs: []interface{}{id},
		Columns:   []string{"id"},
	}

	d, errUser := s.getUser(conditions)
	if errUser != nil || d.Id < 1 {
		return -1
	} else {
		return 1
	}
}

func (s *UserService) isUserExistByUid(uid string) int8 {
	conditions := &storage.Conditions{
		Query:     "uid=?",
		QueryArgs: []interface{}{uid},
		Columns:   []string{"id"},
	}

	d, errUser := s.getUser(conditions)
	if errUser != nil || d.Id < 1 {
		return -1
	} else {
		return 1
	}
}

func (s *UserService) isUserExistByUsername(username string) int8 {
	conditions := &storage.Conditions{
		Query:     "username=?",
		QueryArgs: []interface{}{username},
		Columns:   []string{"id"},
	}

	d, errUser := s.getUser(conditions)
	if errUser != nil || d.Id < 1 {
		return -1
	} else {
		return 1
	}
}

func (s *UserService) getUser(conditions *storage.Conditions) (*gwdb.Users, error) {
	d, errUser := s.modUser.GetUser(conditions)
	if errUser != nil {
		return nil, errUser
	}

	return d, nil
}

func (s *UserService) generatePwd(password string) (interface{}, error) {
	cry := crypto.New()

	cry.Mode = storage.Common
	cry.D = []interface{}{storage.Md5, password}

	return cry.Engine()
}
