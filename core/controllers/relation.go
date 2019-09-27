package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"math"
)

type AddRelationRequest struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
}

func AddRelation(c *gin.Context) {
	resp := new(Resp)
	req := new(AddRelationRequest)
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

type MinuteRelationRequest struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
}

func MinuteRelation(c *gin.Context) {
	resp := new(Resp)
	req := new(MinuteRelationRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.UserId == 0 && req.UserName == "" {
		flog.Log.Errorf("RelationMinute err: %s", "user info empty")
		resp.Error = Error(ParasError, "user info empty")
		return
	}

	who := new(model.User)
	who.Id = req.UserId
	who.Name = req.UserName
	ok, err := who.GetRaw()
	if err != nil {
		flog.Log.Errorf("RelationMinute err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		flog.Log.Errorf("RelationMinute err: %s", "user not fund")
		resp.Error = Error(UserNotFound, "")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("RelationMinute err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	r := new(model.Relation)
	r.UserAId = uu.Id
	r.UserBId = who.Id
	err = r.Minute()
	if err != nil {
		flog.Log.Errorf("RelationMinute err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}

type ListRelationRequest struct {
	UserAId         int64    `json:"user_a_id"`
	UserBId         int64    `json:"user_b_id"`
	UserAName       string   `json:"user_a_name"`
	UserBName       string   `json:"user_b_name"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	Sort            []string `json:"sort"`
	PageHelp
}

type ListRelationResponse struct {
	Relations []model.Relation `json:"relations"`
	PageHelp
}

// who follow you
func ListFollowedRelation(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ListFollowedRelation err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	ListRelation(c, 0, uid)
}

// you follow who
func ListFollowingRelation(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ListFollowingRelation err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	ListRelation(c, uid, 0)
}

// who follow who
func ListAllRelation(c *gin.Context) {
	ListRelation(c, 0, 0)
}

func ListRelation(c *gin.Context, userAId int64, userBId int64) {
	resp := new(Resp)
	req := new(ListRelationRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	// new query list session
	session := model.FaFaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Relation)).Where("1=1")

	if userAId != 0 {
		req.UserAId = userAId
		req.UserAName = ""
	}

	if userBId != 0 {
		req.UserBId = userBId
		req.UserBName = ""
	}

	if req.UserAId != 0 {
		session.And("user_a_id=?", req.UserAId)
	}

	if req.UserBId != 0 {
		session.And("user_b_id=?", req.UserBId)
	}

	if req.UserAName != "" {
		session.And("user_a_name=?", req.UserAName)
	}

	if req.UserBName != "" {
		session.And("user_b_name=?", req.UserBName)
	}

	if req.CreateTimeBegin > 0 {
		session.And("create_time>=?", req.CreateTimeBegin)
	}

	if req.CreateTimeEnd > 0 {
		session.And("create_time<?", req.CreateTimeEnd)
	}

	// count num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("ListRelation err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	cs := make([]model.Relation, 0)
	p := &req.PageHelp
	if total == 0 {
		if p.Limit == 0 {
			p.Limit = 20
		}
	} else {
		// sql build
		p.build(session, req.Sort, model.RelationSortName)
		// do query
		err = session.Find(&cs)
		if err != nil {
			flog.Log.Errorf("ListRelation err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	for k, v := range cs {
		temp := new(model.Relation)
		temp.UserAId = v.UserBId
		temp.UserBId = v.UserAId
		num, err := temp.Count()
		if err != nil {
			flog.Log.Errorf("ListRelation err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if num > 0 {
			cs[k].IsBoth = true
		}
	}
	respResult := new(ListRelationResponse)
	respResult.Relations = cs
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}
