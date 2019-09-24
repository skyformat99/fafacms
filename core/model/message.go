package model

const (
	MessageTypePrivate = 0
	MessageTypePublic  = 1 //  not use
	MessageTypeGlobal  = 2
)

// https://my.oschina.net/cmcm/blog/397323
// Message inside
type Message struct {
	Id            int64 `json:"id" xorm:"bigint pk autoincr"`
	ReceiveUserId int64 `json:"receive_user_id" xorm:"bigint index"`
	CreateTime    int64 `json:"create_time"`
	ReadTime      int64 `json:"read_time"`
	DeleteTime    int64 `json:"delete_time"`
	Status        int   `json:"status" xorm:"not null comment('0 waiting,1 read,2 delete') TINYINT(1) index"`
}

type MessageText struct {
	Id         int64  `json:"id" xorm:"bigint pk autoincr"`
	SendUserId int64  `json:"send_user_id" xorm:"bigint index"`
	CreateTime int64  `json:"create_time"`
	Body       string `json:"body"`
	GroupId    int64  `json:"group_id"`
	Type       int    `json:"type"`
}
