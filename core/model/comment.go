package model

// 灌水表
type Comment struct {
	Id                int    `json:"id" xorm:"bigint pk autoincr"`
	UserId            int    `json:"user_id" xorm:"bigint index"`                                                   // 评论者的用户ID
	ObjectId          int    `json:"object_id" xorm:"bigint index"`                                                 // 评论对应的内容ID
	ObjectUserId      int    `json:"object_user_id" xorm:"bigint index"`                                            // 评论对应的内容所属用户ID
	CommentId         int    `json:"comment_id,omitempty"`                                                          //  对某评论的评论，某评论的ID
	CommentUserId     int    `json:"comment_user_id,omitempty"`                                                     //  对某评论的评论，某评论所属的用户ID
	Status            int    `json:"status" xorm:"not null comment('1 normal, 0 hide，2 deleted') TINYINT(1) index"` // 逻辑删除为2
	Describe          string `json:"describe" xorm:"TEXT"`
	CreateTime        int64  `json:"create_time"`
	UpdateTime        int64  `json:"update_time,omitempty"`
	SuggestUpdateTime int64  `json:"suggest_update_time,omitempty"` // 建议协程更新时间
}
