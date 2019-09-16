package model

type Comment struct {
	Id                int64  `json:"id" xorm:"bigint pk autoincr"`
	UserId            int64  `json:"user_id" xorm:"bigint index"`
	ContentId         int64  `json:"content_id" xorm:"bigint index"`
	ContentUserId     int64  `json:"content_user_id" xorm:"bigint index"`
	CommentId         int64  `json:"comment_id,omitempty" xorm:"bigint index"`
	CommentUserId     int64  `json:"comment_user_id,omitempty" xorm:"bigint index"`
	RootCommentId     int64  `json:"root_comment_id,omitempty" xorm:"bigint index"`
	RootCommentUserId int64  `json:"root_comment_user_id,omitempty" xorm:"bigint index"`
	Describe          string `json:"describe" xorm:"TEXT"`
	CreateTime        int64  `json:"create_time"`
	Status            int    `json:"status" xorm:"not null comment('0 normal, 1 ban') TINYINT(1) index"`
	Cool              int64  `json:"cool"`
	Bad               int64  `json:"bad"`
}

type CommentCool struct {
	Id         int64 `json:"id" xorm:"bigint pk autoincr"`
	UserId     int64 `json:"user_id" xorm:"bigint index(gr)"`
	CommentId  int64 `json:"comment_id,omitempty" xorm:"bigint index(gr)"`
	CreateTime int64 `json:"create_time"`
}

type CommentBad struct {
	Id         int64 `json:"id" xorm:"bigint pk autoincr"`
	UserId     int64 `json:"user_id" xorm:"bigint index(gr)"`
	CommentId  int64 `json:"comment_id,omitempty" xorm:"bigint index(gr)"`
	CreateTime int64 `json:"create_time"`
}
