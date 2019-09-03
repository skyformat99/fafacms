package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/session"
)

const AuthHeader = "Auth"

var (
	AuthDebug = false
	AdminUrl  map[string]int
)

// auth filter
// 授权过滤器
var AuthFilter = func(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		if resp.Error == nil {
			return
		}
		c.AbortWithStatusJSON(200, resp)
	}()

	// get session
	nowUser, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("filter err:%s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	// record log will need uid, monitor who op
	c.Set("uid", nowUser.Id)

	// skip next auth
	if AuthDebug {
		return
	}

	// root user can ignore next auth
	if nowUser.Id == -1 {
		return
	}

	// 超级管理员skr
	if nowUser.Name == "admin" {
		return
	}

	// 未激活不能进入
	if nowUser.Status == 0 {
		flog.Log.Errorf("filter err: not active")
		resp.Error = Error(UserNotActivate, "not active")
		return
	}

	// 被加入了黑名单
	if nowUser.Status == 2 {
		flog.Log.Errorf("filter err: black lock, contact admin")
		resp.Error = Error(UserIsInBlack, "black lock, contact admin")
		return
	}

	// resource is exist
	//r := new(model.Resource)
	url := c.Request.URL.Path
	//r.Url, _ = util2.Sha256([]byte(url))
	//r.Admin = true
	//
	//// resource not found can skip auth
	//if err := r.Get(); err != nil {
	//	flog.Log.Debugf("resource found url:%s, auth err:%s", url, err.Error())
	//	return
	//}

	// resource not found can skip auth
	resourceId, exist := AdminUrl[url]
	if !exist {
		return
	}

	// if group has this resource
	gr := new(model.GroupResource)
	gr.GroupId = nowUser.GroupId
	gr.ResourceId = resourceId
	exist, err = model.FafaRdb.Client.Exist(gr)
	if err != nil {
		flog.Log.Errorf("filter err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		// not found
		flog.Log.Errorf("filter err:%s", "resource not allow")
		resp.Error = Error(UserAuthPermit, "resource not allow")
		return
	}
}

// 获取用户信息，存于Session中的
func GetUserSession(c *gin.Context) (*model.User, error) {
	// 请求链只需查询一次用户信息
	if v, exist := c.Get("everAuth"); exist {
		return v.(*model.User), nil
	}
	// 检查并拉出用户信息
	token := c.GetHeader(AuthHeader)
	user, err := session.FafaSessionMgr.CheckToken(token)
	if err != nil {
		return nil, err
	}

	c.Set("everAuth", user)
	return user, nil
}

func SetUserSession(user *model.User) (string, error) {
	if user == nil {
		return "", errors.New("user nil")
	}

	// 只允许单点登录
	session.FafaSessionMgr.DeleteUserToken(user.Id)
	return session.FafaSessionMgr.SetToken(user, 24*3600*7)
}

func DeleteUserSession(c *gin.Context) error {
	token := c.GetHeader("Auth")
	err := session.FafaSessionMgr.DeleteToken(token)
	return err
}

func DeleteUserAllSession(id int) error {
	err := session.FafaSessionMgr.DeleteUserToken(id)
	return err
}

func RefreshUserSession(c *gin.Context) error {
	token := c.GetHeader("Auth")
	err := session.FafaSessionMgr.RefreshToken(token)
	return err
}
