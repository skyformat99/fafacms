package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"strings"
)

type LoginRequest struct {
	UserName string `json:"user_name"`
	PassWd   string `json:"pass_wd"`
}

func Login(c *gin.Context) {
	resp := new(Resp)
	req := new(LoginRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	// check session
	userInfo, _ := GetUserSession(c)
	if userInfo != nil {
		//c.Set("skipLog", true)
		c.Set("uid", userInfo.Id)
		resp.Flag = true
		return
	}

	// paras not empty
	if req.UserName == "" || req.PassWd == "" {
		flog.Log.Errorf("login err:%s", "paras wrong")
		resp.Error = Error(ParasError, "field username or pass_wd")
		return
	}

	// super root user login
	if req.UserName == "hunterhug" && req.PassWd == "hunterhug" {
		u := new(model.User)
		u.Id = -1
		u.Name = "hunterhug"
		u.Status = 1
		token, err := SetUserSession(u)
		if err != nil {
			flog.Log.Errorf("login set root err:%s", err.Error())
			resp.Error = Error(SetUserSessionError, err.Error())
			return
		}

		c.Set("uid", u.Id)
		resp.Data = token
		resp.Flag = true
		return
	}

	// common people login
	uu := new(model.User)
	if strings.Contains(req.UserName, "@") {
		uu.Email = req.UserName
	} else {
		uu.Name = req.UserName
	}
	uu.Password = req.PassWd
	ok, err := uu.GetRaw()
	if err != nil {
		flog.Log.Errorf("login err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		flog.Log.Errorf("login err:%s", "user or password wrong")
		resp.Error = Error(LoginWrong, "user or password wrong")
		return
	}

	c.Set("uid", uu.Id)

	// 就算未激活，或者黑名单都可以登录，但授权的API无法使用，激活用户的时候session会生成一个新的，即新的token，并且用户缓存会刷新
	token, err := SetUserSession(uu)
	if err != nil {
		flog.Log.Errorf("login err:%s", err.Error())
		resp.Error = Error(SetUserSessionError, err.Error())
		return
	}

	resp.Data = token
	resp.Flag = true
}

func Logout(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSON(c, 200, resp)
	}()
	user, err := GetUserSession(c)

	if err != nil {
		flog.Log.Errorf("logout err:%s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	if user != nil {
		err = DeleteUserSession(c)
		if err != nil {
			flog.Log.Errorf("logout err:%s", err.Error())
			resp.Error = Error(DeleteUserSessionError, err.Error())
			return
		}
	}
	resp.Flag = true
}

func Refresh(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSON(c, 200, resp)
	}()
	user, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("refresh err:%s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	if user != nil {
		err = RefreshUserSession(c)
		if err != nil {
			flog.Log.Errorf("refresh err:%s", err.Error())
			resp.Error = Error(RefreshUserCacheError, err.Error())
			return
		}
	}
	resp.Flag = true
}
