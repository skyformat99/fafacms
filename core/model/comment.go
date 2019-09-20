package model

import (
	"errors"
	"time"
)

const (
	CommentTypeOfContent     = 0
	CommentTypeOfRootComment = 1
	CommentTypeOfComment     = 2
	CommentAnonymous         = 1

	AnonymousUser                = "void"
	CommentContentDeleteDescribe = "content deleted"
	CommentDeleteDescribe        = "comment deleted"
	CommentBanDescribe           = "comment ban"
)

type UserHelper struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	NickName  string `json:"nick_name"`
	HeadPhoto string `json:"head_photo"`
}

func GetUserHelper(u *User) *UserHelper {
	if u == nil || u.Id == 0 {
		return nil
	}

	uu := new(UserHelper)
	uu.Id = u.Id
	uu.NickName = u.NickName
	uu.HeadPhoto = u.HeadPhoto
	uu.Name = u.Name
	return uu
}

type Comment struct {
	Id                        int64       `json:"id" xorm:"bigint pk autoincr"`
	UserId                    int64       `json:"-" xorm:"bigint index"`
	UserName                  string      `json:"-" xorm:"index"`
	ContentId                 int64       `json:"content_id" xorm:"bigint index"`
	ContentTitle              string      `json:"content_title"`
	ContentUserId             int64       `json:"-" xorm:"bigint index"`
	ContentUserName           string      `json:"-" xorm:"index"`
	CommentId                 int64       `json:"comment_id,omitempty" xorm:"bigint index"`
	CommentUserId             int64       `json:"-" xorm:"bigint index"`
	CommentUserName           string      `json:"-" xorm:"index"`
	RootCommentId             int64       `json:"root_comment_id,omitempty" xorm:"bigint index"`
	RootCommentUserId         int64       `json:"-" xorm:"bigint index"`
	RootCommentUserName       string      `json:"-" xorm:"xindex"`
	HelperUser                *UserHelper `json:"user,omitempty" xorm:"-"`
	HelperContentUser         *UserHelper `json:"comment_user,omitempty" xorm:"-"`
	HelperRootCommentUser     *UserHelper `json:"root_comment_user,omitempty" xorm:"-"`
	HelperCommentUser         *UserHelper `json:"comment_user,omitempty" xorm:"-"`
	Describe                  string      `json:"describe" xorm:"TEXT"`
	HelperRootCommentDescribe string      `json:"root_comment_describe,omitempty" xorm:"-"`
	HelperCommentDescribe     string      `json:"comment_describe,omitempty" xorm:"-"`
	CreateTime                int64       `json:"create_time"`
	Status                    int         `json:"status" xorm:"not null comment('0 normal, 1 ban') TINYINT(1) index"`
	Cool                      int64       `json:"cool"`
	Bad                       int64       `json:"bad"`
	CommentType               int         `json:"comment_type"`
	CommentAnonymous          int         `json:"comment_anonymous"`
	IsDelete                  int         `json:"is_delete,omitempty"`
	DeleteTime                int64       `json:"delete_time,omitempty"`
}

type CommentCool struct {
	Id         int64 `json:"id" xorm:"bigint pk autoincr"`
	UserId     int64 `json:"user_id" xorm:"bigint index(gr)"`
	CommentId  int64 `json:"comment_id,omitempty" xorm:"bigint index(gr)"`
	ContentId  int64 `json:"content_id" xorm:"bigint index"`
	CreateTime int64 `json:"create_time"`
}

type CommentBad struct {
	Id         int64 `json:"id" xorm:"bigint pk autoincr"`
	UserId     int64 `json:"user_id" xorm:"bigint index(gr)"`
	CommentId  int64 `json:"comment_id,omitempty" xorm:"bigint index(gr)"`
	ContentId  int64 `json:"content_id" xorm:"bigint index"`
	CreateTime int64 `json:"create_time"`
}

func (c *Comment) InsertOne() error {
	se := FaFaRdb.Client.NewSession()
	err := se.Begin()
	if err != nil {
		return err
	}

	c.CreateTime = time.Now().Unix()
	num, err := se.InsertOne(c)
	if err != nil {
		se.Rollback()
		return err
	}

	if num == 0 {
		se.Rollback()
		return errors.New("some err")
	}

	num, err = se.Where("id=?", c.ContentId).Incr("comment_num").Update(new(Content))
	if err != nil {
		se.Rollback()
		return err
	}

	if num == 0 {
		se.Rollback()
		return errors.New("some err")
	}

	err = se.Commit()
	if err != nil {
		se.Rollback()
		return err
	}

	return nil
}

func (c *Comment) Get() (bool, error) {
	if c.Id == 0 {
		return false, errors.New("where is empty")
	}
	return FaFaRdb.Client.Get(c)
}

func (c *Comment) Delete() (err error) {
	if c.Id == 0 {
		return errors.New("where is empty")
	}

	se := FaFaRdb.Client.NewSession()
	err = se.Begin()
	if err != nil {
		return err
	}

	c.IsDelete = 1
	c.DeleteTime = time.Now().Unix()
	num, err := se.Where("id=?", c.Id).Cols("is_delete", "delete_time").Update(c)
	if err != nil {
		return err
	}

	if num == 0 {
		se.Rollback()
		return errors.New("some err")
	}

	num, err = se.Where("id=?", c.ContentId).Decr("comment_num").Update(new(Content))
	if err != nil {
		se.Rollback()
		return err
	}

	if num == 0 {
		se.Rollback()
		return errors.New("some err")
	}

	return
}
