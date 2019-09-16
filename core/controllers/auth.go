package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/session"
)

// every API which need auth should take a HTTP header `Auth`
const AuthHeader = "Auth"

var (
	// if you want skip auth you can set it true
	AuthDebug = false

	// those api will be check resource
	AdminUrl map[string]int64
)

// api access auth filter
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

	// admin user is skr
	if nowUser.Name == "admin" {
		return
	}

	// not active will be refuse
	if nowUser.Status == 0 {
		flog.Log.Errorf("filter err: not active")
		resp.Error = Error(UserNotActivate, "not active")
		return
	}

	// black user will be refuse
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

	// resource not found in group will be refuse
	if !exist {
		flog.Log.Errorf("filter err:%s", "resource not allow")
		resp.Error = Error(UserAuthPermit, "resource not allow")
		return
	}
}

// get the info of user，will save in redis Session
func GetUserSession(c *gin.Context) (*model.User, error) {
	// get the info from context if exist
	if v, exist := c.Get("everAuth"); exist {
		return v.(*model.User), nil
	}

	// get token from HTTP header and check if it is exist
	token := c.GetHeader(AuthHeader)
	user, err := session.FafaSessionMgr.CheckToken(token)
	if err != nil {
		return nil, err
	}

	// set the info into context
	c.Set("everAuth", user)
	return user, nil
}

func SetUserSession(user *model.User) (string, error) {
	if user == nil {
		return "", errors.New("user nil")
	}

	// single login
	// we only allow one token exist, other token will be delete.
	session.FafaSessionMgr.DeleteUserToken(user.Id)
	return session.FafaSessionMgr.SetToken(user, 24*3600*7)
}

func DeleteUserSession(c *gin.Context) error {
	token := c.GetHeader(AuthHeader)
	err := session.FafaSessionMgr.DeleteToken(token)
	return err
}

func DeleteUserAllSession(id int64) error {
	err := session.FafaSessionMgr.DeleteUserToken(id)
	return err
}

func RefreshUserSession(c *gin.Context) error {
	token := c.GetHeader(AuthHeader)
	err := session.FafaSessionMgr.RefreshToken(token)
	return err
}
