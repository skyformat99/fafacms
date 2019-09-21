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

	AnonymousUser = "匿名" // just ignore this, not use at all
)

type ContentHelper struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	IsHide      bool   `json:"is_hide"`
	IsInRubbish bool   `json:"is_in_rubbish"`
	IsBan       bool   `json:"is_ban"`
	UserId      int64  `json:"user_id"`
	UserName    string `json:"user_name"`
	Seo         string `json:"seo"`
	Status      int    `json:"status"`
}

// get contents from id, if all false, which is deleted or hide will not include in map
func GetContentHelper(ids []int64, all bool, yourUserId int64) (back map[int64]ContentHelper, err error) {
	back = make(map[int64]ContentHelper)
	cs := make([]Content, 0)
	err = FaFaRdb.Client.Cols("id", "title", "status", "seo", "user_name", "user_id").In("id", ids).Find(&cs)
	if err != nil {
		return
	}

	for _, v := range cs {
		temp := ContentHelper{
			Id:       v.Id,
			Title:    v.Title,
			Seo:      v.Seo,
			UserName: v.UserName,
			UserId:   v.UserId,
			Status:   v.Status,
		}

		if v.Status == 1 {
			temp.IsHide = true
		}
		if v.Status == 2 {
			temp.IsBan = true
		}
		if v.Status == 3 {
			temp.IsInRubbish = true
		}

		if v.Status == 1 || v.Status == 3 {
			// if is yourself will still back the data
			if !all && yourUserId != temp.UserId {
				continue
			}
		}
		back[v.Id] = temp
	}

	return
}

type UserHelper struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	NickName  string `json:"nick_name"`
	HeadPhoto string `json:"head_photo"`
	IsVip     bool   `json:"is_vip"`
}

type CommentHelper struct {
	Id            int64  `json:"id"`
	Describe      string `json:"describe"`
	CreateTime    int64  `json:"create_time"`
	CommentDelete bool   `json:"is_delete"`
	IsBan         bool   `json:"is_ban"`
	IsAnonymous   bool   `json:"is_anonymous"`
	UserId        int64  `json:"user_id"`
	Cool          int64  `json:"cool"`
	Bad           int64  `json:"bad"`
}

// get comments from id, if all false, some will not show, and user info related will also show
func GetCommentAndCommentUser(ids []int64, all bool, yourUserId int64) (comments map[int64]CommentHelper, users map[int64]UserHelper, err error) {
	comments = make(map[int64]CommentHelper)
	users = make(map[int64]UserHelper)

	cms := make([]Comment, 0)
	err = FaFaRdb.Client.Cols("id", "user_id", "create_time", "describe", "is_delete", "bad", "cool", "status", "comment_anonymous").In("id", ids).Find(&cms)
	if err != nil {
		return
	}

	for _, v := range cms {
		temp := CommentHelper{
			Id:            v.Id,
			CreateTime:    v.CreateTime,
			Describe:      v.Describe,
			CommentDelete: v.IsDelete == 1,
			IsBan:         v.Status == 1,
			IsAnonymous:   v.CommentAnonymous == CommentAnonymous,
			UserId:        v.UserId,
			Cool:          v.Cool,
			Bad:           v.Bad,
		}
		if !all {
			// delete will not show others
			if temp.CommentDelete {
				temp = CommentHelper{Id: v.Id, CommentDelete: true}
			} else {

				// userId is you, do nothing
				if yourUserId != temp.UserId {
					// user info hide
					if temp.IsAnonymous {
						temp.UserId = 0
					}

					// describe hide
					if temp.IsBan {
						temp.Describe = ""
					}
				}
			}
		}

		comments[v.Id] = temp
	}

	// user id trip repeat
	userHelper := make(map[int64]struct{})
	for _, v := range comments {
		if v.UserId != 0 {
			userHelper[v.UserId] = struct{}{}
		}
	}

	userIds := make([]int64, 0)
	for k := range userHelper {
		userIds = append(userIds, k)
	}

	uu := make([]User, 0)
	err = FaFaRdb.Client.Cols("id", "vip", "name", "nick_name", "head_photo").In("id", userIds).Find(&uu)
	if err != nil {
		return
	}

	for _, v := range uu {
		temp := UserHelper{
			Id:        v.Id,
			Name:      v.Name,
			NickName:  v.NickName,
			HeadPhoto: v.HeadPhoto,
		}

		temp.IsVip = v.Vip == 1
		users[v.Id] = temp
	}

	return
}

type Comment struct {
	Id                  int64  `json:"id" xorm:"bigint pk autoincr"`
	UserId              int64  `json:"-" xorm:"bigint index"`
	UserName            string `json:"-" xorm:"index"`
	ContentId           int64  `json:"content_id" xorm:"bigint index"`
	ContentTitle        string `json:"content_title"` // may be content delete so this field keep
	ContentUserId       int64  `json:"-" xorm:"bigint index"`
	ContentUserName     string `json:"-" xorm:"index"`
	CommentId           int64  `json:"comment_id" xorm:"bigint index"`
	CommentUserId       int64  `json:"-" xorm:"bigint index"`
	CommentUserName     string `json:"-" xorm:"index"`
	RootCommentId       int64  `json:"root_comment_id" xorm:"bigint index"`
	RootCommentUserId   int64  `json:"-" xorm:"bigint index"`
	RootCommentUserName string `json:"-" xorm:"index"`
	Describe            string `json:"-" xorm:"TEXT"`
	CreateTime          int64  `json:"-"`
	Status              int    `json:"-" xorm:"not null comment('0 normal, 1 ban') TINYINT(1) index"`
	Cool                int64  `json:"-"`
	Bad                 int64  `json:"-"`
	CommentType         int    `json:"comment_type"` // 0 comment to content, 1 comment to comment, 2 comment to comment more
	CommentAnonymous    int    `json:"-"`
	IsDelete            int    `json:"-"`
	DeleteTime          int64  `json:"-"`
}

type CommentExtra struct {
	Users    map[int64]UserHelper    `json:"users"`
	Comments map[int64]CommentHelper `json:"comments"`
	Contents map[int64]ContentHelper `json:"contents"`
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

	err = se.Commit()
	if err != nil {
		se.Rollback()
		return err
	}

	return
}
