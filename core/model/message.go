package model

const (
	MessageTypeCoolForContent    = 0
	MessageTypeCoolForComment    = 1
	MessageTypeCommentForComment = 2
	MessageTypeCommentForContent = 3
	MessageTypeBanForContent     = 4
	MessageTypeBanForComment     = 5
	MessageTypeSystem            = 6
)

type Message struct {
	Id              int64  `json:"id" xorm:"bigint pk autoincr"`
	TargetUserId    int64  `json:"target_user_id" xorm:"bigint index"`
	CreateTime      int64  `json:"create_time"`
	Status          int    `json:"status" xorm:"not null comment('0 waiting,1 read') TINYINT(1) index"`
	MessageType     int    `json:"message_type" xorm:"not null TINYINT(1) index"`
	ContentId       int64  `json:"content_id" xorm:"bigint index"` // 034
	OriginUserId    int64  `json:"origin_user_id" xorm:"bigint index"`
	OriginCommentId int64  `json:"origin_comment_id" xorm:"bigint index"` // 23
	TargetCommentId int64  `json:"target_comment_id" xorm:"bigint index"` // 125
	MessageBody     string `json:"message_body"`                          // 6
}
