package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
)

type RelationAddRequest struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
}

func RelationAdd(c *gin.Context) {
	resp := new(Resp)
	req := new(RelationAddRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.UserId == 0 && req.UserName == "" {
		flog.Log.Errorf("RelationAdd err: %s", "user info empty")
		resp.Error = Error(ParasError, "user info empty")
		return
	}

	who := new(model.User)
	who.Id = req.UserId
	who.Name = req.UserName
	ok, err := who.GetRaw()
	if err != nil {
		flog.Log.Errorf("RelationAdd err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		flog.Log.Errorf("RelationAdd err: %s", "user not fund")
		resp.Error = Error(UserNotFound, "")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("RelationAdd err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	r := new(model.Relation)
	r.UserAId = uu.Id
	r.UserBId = who.Id
	r.UserAName = uu.Name
	r.UserBName = who.Name
	err = r.Add()
	if err != nil {
		flog.Log.Errorf("RelationAdd err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}
