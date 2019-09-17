package model

type Message struct {
	Id                int64  `json:"id" xorm:"bigint pk autoincr"`
	UserId            int64  `json:"user_id" xorm:"bigint index"`
	//ContentId         int64  `json:"content_id" xorm:"bigint index"`
	//ContentUserId     int64  `json:"content_user_id" xorm:"bigint index"`
	//CommentId         int64  `json:"comment_id,omitempty" xorm:"bigint index"`
	//CommentUserId     int64  `json:"comment_user_id,omitempty" xorm:"bigint index"`
	//RootCommentId     int64  `json:"root_comment_id,omitempty" xorm:"bigint index"`
	//RootCommentUserId int64  `json:"root_comment_user_id,omitempty" xorm:"bigint index"`
	//Describe          string `json:"describe" xorm:"TEXT"`
	//CreateTime        int64  `json:"create_time"`
	//Status            int    `json:"status" xorm:"not null comment('0 waiting,1 read') TINYINT(1) index"`
	//MessageType int `json:"message_type" xorm:"not null TINYINT(1) index"`
}