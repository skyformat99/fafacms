package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"strings"
)

var htmlEscaper = strings.NewReplacer(
	`<`, "&lt;",
	`>`, "&gt;",
)

type CreateCommentRequest struct {
	ContentId   int64  `json:"content_id"`
	CommentId   int64  `json:"comment_id"`
	IsToComment bool   `json:"is_to_comment"`
	Body        string `json:"body"`
}

func CreateComment(c *gin.Context) {
	resp := new(Resp)
	req := new(CreateCommentRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	req.Body = strings.TrimSpace(req.Body)
	if len(req.Body) == 0 {
		flog.Log.Errorf("CreateComment err: %s", "body empty")
		resp.Error = Error(ParasError, "body empty")
		return
	}

	req.Body = htmlEscaper.Replace(req.Body)
	if !req.IsToComment && req.ContentId == 0 {
		flog.Log.Errorf("CreateComment err: %s", "content_id empty")
		resp.Error = Error(ParasError, "content_id empty")
		return
	}

	if req.IsToComment && req.CommentId == 0 {
		flog.Log.Errorf("CreateComment err: %s", "comment_id empty")
		resp.Error = Error(ParasError, "comment_id empty")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("CreateComment err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	// comment to content
	if !req.IsToComment {
		content := new(model.Content)
		content.Id = req.ContentId
		ok, err := content.GetByRaw()
		if err != nil {
			flog.Log.Errorf("CreateComment err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("CreateComment err: %s", "content not found")
			resp.Error = Error(ContentNotFound, "")
			return
		}

		if content.Status == 0 && content.Version > 0 {
			cm := new(model.Comment)
			cm.ContentId = content.Id
			cm.ContentUserId = content.UserId
			cm.UserId = uu.Id
			cm.Describe = req.Body
			err = cm.InsertOne()
			if err != nil {
				flog.Log.Errorf("CreateComment err: %s", err.Error())
				resp.Error = Error(DBError, err.Error())
				return
			}
			resp.Data = cm.Id
		} else {
			flog.Log.Errorf("CreateComment err: %s", "content status not 0 or not publish")
			if content.Status == 2 {
				resp.Error = Error(ContentBanPermit, "")
			} else {
				resp.Error = Error(ContentNotFound, "")
			}
			return
		}

		resp.Flag = true
		return
	}

	targetComment := new(model.Comment)
	targetComment.Id = req.CommentId
	ok, err := targetComment.Get()
	if err != nil {
		flog.Log.Errorf("CreateComment err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		flog.Log.Errorf("CreateComment err: %s", "comment not found")
		resp.Error = Error(CommentNotFound, "")
		return
	}

	if targetComment.Status == 1 {
		flog.Log.Errorf("CreateComment err: %s", "comment ban")
		resp.Error = Error(CommentBanPermit, "")
		return
	}

	content := new(model.Content)
	content.Id = targetComment.ContentId
	ok, err = content.GetByRaw()
	if err != nil {
		flog.Log.Errorf("CreateComment err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !ok {
		flog.Log.Errorf("CreateComment err: %s", "content not found")
		resp.Error = Error(ContentNotFound, "")
		return
	}

	if content.Status != 0 && content.Version > 0 {
		flog.Log.Errorf("CreateComment err: %s", "content status not 0 or not publish")
		if content.Status == 2 {
			resp.Error = Error(ContentBanPermit, "")
		} else {
			resp.Error = Error(ContentNotFound, "")
		}
		return
	}

	newComment := new(model.Comment)
	newComment.ContentId = content.Id
	newComment.ContentUserId = content.UserId
	newComment.UserId = uu.Id
	newComment.Describe = req.Body
	if targetComment.RootCommentId == 0 {
		newComment.RootCommentId = targetComment.Id
		newComment.RootCommentUserId = targetComment.UserId
	} else {
		newComment.CommentId = targetComment.Id
		newComment.CommentUserId = targetComment.UserId
		newComment.RootCommentId = targetComment.RootCommentId
		newComment.RootCommentUserId = targetComment.RootCommentUserId
	}

	err = newComment.InsertOne()
	if err != nil {
		flog.Log.Errorf("CreateComment err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Data = newComment.Id
	resp.Flag = true
}

func UpdateComment(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func DeleteComment(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func TakeComment(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func ListComment(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}

func ListHomeComment(c *gin.Context) {
	resp := new(Resp)
	defer func() {
		JSONL(c, 200, nil, resp)
	}()
}
