package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/util"
	"math"
)

type ListMessageRequest struct {
	MessageId       int64    `json:"message_id"`
	MessageType     int      `json:"message_type" validate:"oneof=-1 0 1 2 3 4 5 6 7 8 9"`
	UserId          int64    `json:"user_id"`
	ReceiveStatus   int      `json:"receive_status" validate:"oneof=-1 0 1 2"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	Sort            []string `json:"sort"`
	PageHelp
}

type ListMessageResponse struct {
	Messages      []model.Message               `json:"messages"`
	Comments      map[int64]model.Comment       `json:"comments"`
	Contents      map[int64]model.ContentHelper `json:"contents"`
	ExtraUsers    map[int64]model.UserHelper    `json:"extra_users"`
	ExtraComments map[int64]model.CommentHelper `json:"extra_comments"`
	UnRead        map[string]int                `json:"un_read"`
	PageHelp
}

func ListMessageHelper(c *gin.Context, isAdmin bool) {
	resp := new(Resp)

	respResult := new(ListMessageResponse)
	req := new(ListMessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("ListMessageHelper err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	var yourUserId int64 = 0
	var all = true

	if !isAdmin {
		uu, err := GetUserSession(c)
		if err != nil {
			flog.Log.Errorf("ListMessageHelper err: %s", err.Error())
			resp.Error = Error(GetUserSessionError, err.Error())
			return
		}

		yourUserId = uu.Id
		all = false
	}

	// new query list session
	session := model.FaFaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Message)).Where("1=1")

	if req.MessageId != 0 {
		session.And("id=?", req.MessageId)
	}

	if req.MessageType != -1 {
		session.And("message_type=?", req.MessageType)
	}

	if req.ReceiveStatus != -1 {
		session.And("receive_status=?", req.ReceiveStatus)
	}

	if req.CreateTimeBegin > 0 {
		session.And("create_time>=?", req.CreateTimeBegin)
	}

	if req.CreateTimeEnd > 0 {
		session.And("create_time<?", req.CreateTimeEnd)
	}

	if !all {
		// search your message which not delete
		session.And("receive_user_id=?", yourUserId).And("receive_status!=?", 2)
	}

	if all {
		if req.UserId != 0 {
			session.And("receive_user_id=?", req.UserId)
		}
	}

	// count all message num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	countMap, err := model.GroupCount(yourUserId)
	if err != nil {
		flog.Log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	cs := make([]model.Message, 0)
	p := &req.PageHelp
	if total == 0 {
		if p.Limit == 0 {
			p.Limit = 20
		}
	} else {
		// sql build
		p.build(session, req.Sort, model.MessageSortName)
		// do query
		err = session.Find(&cs)
		if err != nil {
			flog.Log.Errorf("ListMessageHelper err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	contentIds := make(map[int64]struct{})
	userIds := make(map[int64]struct{})
	commentIds := make(map[int64]struct{})
	for _, v := range cs {
		// content ids collect
		if v.ContentId != 0 {
			contentIds[v.ContentId] = struct{}{}
		}

		switch v.MessageType {
		case model.MessageTypeCommentForContent, model.MessageTypeCommentForComment:
			// if comment is anonymous must hide user id
			if v.CommentAnonymous == 1 && v.IsYourSelf == 0 {
				if !all {
					v.UserId = 0
				}
			}
		}

		// user id collect
		if v.UserId != 0 {
			userIds[v.UserId] = struct{}{}
		}

		// comment id collect
		if v.CommentId != 0 {
			commentIds[v.CommentId] = struct{}{}
		}
	}

	// get all none delete comment
	comments, err := model.GetComment(util.MapToArray(commentIds), all)
	if err != nil {
		flog.Log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// extra comment id collect
	commentIds2 := make(map[int64]struct{})
	for k := range comments {
		commentIds2[k] = struct{}{}
	}

	// get extra comment info
	comments2, user2, err := model.GetCommentAndCommentUser(util.MapToArray(commentIds2), all, util.MapToArray(userIds), yourUserId)
	if err != nil {
		flog.Log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// get content base info
	contents, err := model.GetContentHelper(util.MapToArray(contentIds), all, yourUserId)
	if err != nil {
		flog.Log.Errorf("ListMessageHelper err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// result
	respResult.Messages = cs
	respResult.UnRead = countMap
	respResult.Comments = comments
	respResult.ExtraComments = comments2
	respResult.ExtraUsers = user2
	respResult.Contents = contents
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	p.Total = int(total)
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

func ListMessage(c *gin.Context) {
	ListMessageHelper(c, false)
}

func ListAllMessage(c *gin.Context) {
	ListMessageHelper(c, true)
}

type MessageRequest struct {
	Id int64 `json:"id"`
}

func ReadMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(MessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.Id == 0 {
		flog.Log.Errorf("ReadMessage err: %s", "message_id empty")
		resp.Error = Error(ParasError, "message_id empty")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ReadMessage err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	m := new(model.Message)
	m.Id = req.Id
	m.ReceiveUserId = uu.Id
	m.ReceiveStatus = 1
	err = m.Update()
	if err != nil {
		flog.Log.Errorf("ReadMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}

func DeleteMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(MessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.Id == 0 {
		flog.Log.Errorf("DeleteMessage err: %s", "message_id empty")
		resp.Error = Error(ParasError, "message_id empty")
		return
	}
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("DeleteMessage err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	m := new(model.Message)
	m.Id = req.Id
	m.ReceiveUserId = uu.Id
	m.ReceiveStatus = 2
	err = m.Update()
	if err != nil {
		flog.Log.Errorf("DeleteMessage err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Flag = true
}

type SendPrivateMessageRequest struct {
	UserId  int64  `json:"user_id"`
	Message string `json:"message"`
}

// todo
func SendPrivateMessage(c *gin.Context) {
	resp := new(Resp)
	req := new(SendPrivateMessageRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.UserId == 0 {
		flog.Log.Errorf("SendPrivateMessage err: %s", "user_id empty")
		resp.Error = Error(ParasError, "user_id empty")
		return
	}

	//uu, err := GetUserSession(c)
	//if err != nil {
	//	flog.Log.Errorf("SendPrivateMessage err: %s", err.Error())
	//	resp.Error = Error(GetUserSessionError, err.Error())
	//	return
	//}
	resp.Flag = true
}
