package model

const (
	// private message 1:1
	MessageTypePrivate = 0

	// system message inside
	MessageTypeSystem = 2

	// global message 1:N
	MessageTypeGlobal = 3

	// who follow you message
	BodyTypeFollow = 0

	// who comment your content or comment message
	BodyTypeComment = 1

	// who good your content or comment message
	BodyTypeGoodContent = 2
	BodyTypeGoodComment = 3

	// comment or content be ban by system message
	BodyTypeContentBan = 4
	BodyTypeCommentBan = 5

	// comment or content be ban by recover message
	BodyTypeContentRecover = 6
	BodyTypeCommentRecover = 7

	// who you follow publish content
	BodyTypeContentPublish = 8

	// who send a message to you
	BodyTypePrivate = 9

	// global send a message to you
	BodyTypeGlobal = 10
)

// Message inside
type Message struct {
	Id            int64  `json:"id" xorm:"bigint pk autoincr"`
	SendUserId    int64  `json:"send_user_id" xorm:"bigint index"`
	ReceiveUserId int64  `json:"receive_user_id" xorm:"bigint index"`
	CreateTime    int64  `json:"create_time"`
	ReadTime      int64  `json:"read_time"`
	DeleteTime    int64  `json:"delete_time"`
	SendStatus    int    `json:"send_status" xorm:"not null comment('0 waiting,1 read,2 delete') TINYINT(1) index"`
	ReceiveStatus int    `json:"receive_status" xorm:"not null comment('0 waiting,1 read,2 delete') TINYINT(1) index"`
	Body          string `json:"body"`
	BodyType      int    `json:"body_type"`
	MessageType   int    `json:"message_type"` // private 1:1, system 1:1 special, global: 1:N
}

var MessageSortName = []string{"=id", "-create_time", "=status", "=type", "=send_user_id", "=receive_user_id"}
