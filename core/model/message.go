package model

import (
	"errors"
	"fmt"
	"github.com/hunterhug/fafacms/core/util"
	"time"
)

const (
	// who comment your content or comment message
	MessageTypeCommentForContent = 0
	MessageTypeCommentForComment = 1

	// who good your content or comment message
	MessageTypeGoodContent = 2
	MessageTypeGoodComment = 3

	// comment or content be ban by system message
	MessageTypeContentBan = 4
	MessageTypeCommentBan = 5

	// comment or content be ban by recover message
	MessageTypeContentRecover = 6
	MessageTypeCommentRecover = 7

	// who follow you message
	MessageTypeFollow = 8

	// who you follow publish content
	MessageTypeContentPublish = 9

	// who send a message to you
	MessageTypePrivate = 10

	// global send a message to you
	MessageTypeGlobal = 11
)

// Message inside
type Message struct {
	Id               int64  `json:"id" xorm:"bigint pk autoincr"`
	SendUserId       int64  `json:"send_user_id,omitempty" xorm:"bigint index"`
	SendMessage      string `json:"send_message,omitempty"`
	ReceiveUserId    int64  `json:"receive_user_id" xorm:"bigint index"`
	CreateTime       int64  `json:"create_time"`
	ReadTime         int64  `json:"read_time"`
	DeleteTime       int64  `json:"delete_time,omitempty"`
	SendStatus       int    `json:"-" xorm:"not null comment('0 waiting,1 read,2 delete') TINYINT(1) index"`
	ReceiveStatus    int    `json:"receive_status" xorm:"not null comment('0 waiting,1 read,2 delete') TINYINT(1) index"`
	UserId           int64  `json:"user_id" xorm:"bigint index"`
	ContentId        int64  `json:"content_id" xorm:"bigint index"`
	ContentTitle     string `json:"content_title"`
	CommentId        int64  `json:"comment_id" xorm:"bigint index"`
	CommentDescribe  string `json:"comment_describe"`
	CommentAnonymous int    `json:"comment_anonymous"`
	PublishAgain     int    `json:"publish_again"`
	MessageType      int    `json:"message_type" xorm:"index"`
	IsYourSelf       int    `json:"is_your_self"`
}

var MessageSortName = []string{"=id", "-create_time", "=receive_status", "=message_type", "=send_user_id", "=receive_user_id"}

func CommentAbout(userId int64, receiveUserId int64, contentId int64, contentTitle string, commentId int64, commentDescribe string, messageType int, commentAnonymous bool) error {
	m := new(Message)
	m.UserId = userId
	m.CommentId = commentId
	m.CommentDescribe = commentDescribe
	m.ContentId = contentId
	m.ContentTitle = contentTitle
	m.ReceiveUserId = receiveUserId
	m.MessageType = messageType
	if commentAnonymous {
		m.CommentAnonymous = 1
	}
	return m.Insert()
}

func ContentAbout(userId int64, receiveUserId int64, contentId int64, contentTitle string, messageType int) error {
	m := new(Message)
	m.UserId = userId
	m.ContentId = contentId
	m.ContentTitle = contentTitle
	m.ReceiveUserId = receiveUserId
	m.MessageType = messageType
	return m.Insert()
}

func (m *Message) Insert() error {
	m.CreateTime = time.Now().Unix()
	if m.UserId == m.ReceiveUserId {
		m.IsYourSelf = 1
	}
	_, err := FaFaRdb.Client.InsertOne(m)
	return err
}

func GoodContent(userId int64, receiveUserId int64, contentId int64, contentTitle string) error {
	return ContentAbout(userId, receiveUserId, contentId, contentTitle, MessageTypeGoodContent)
}

func BanContent(receiveUserId int64, contentId int64, contentTitle string) error {
	return ContentAbout(0, receiveUserId, contentId, contentTitle, MessageTypeContentBan)
}

func GoodComment(userId int64, receiveUserId int64, contentId int64, contentTitle string, commentId int64, commentDescribe string) error {
	return CommentAbout(userId, receiveUserId, contentId, contentTitle, commentId, commentDescribe, MessageTypeGoodComment, false)
}

func BanComment(receiveUserId int64, contentId int64, contentTitle string, commentId int64, commentDescribe string) error {
	return CommentAbout(0, receiveUserId, contentId, contentTitle, commentId, commentDescribe, MessageTypeCommentBan, false)
}

func RecoverContent(receiveUserId int64, contentId int64, contentTitle string) error {
	return ContentAbout(0, receiveUserId, contentId, contentTitle, MessageTypeContentRecover)
}

func RecoverComment(receiveUserId int64, contentId int64, contentTitle string, commentId int64, commentDescribe string) error {
	return CommentAbout(0, receiveUserId, contentId, contentTitle, commentId, commentDescribe, MessageTypeCommentRecover, false)
}

func CommentForContent(userId int64, receiveUserId int64, contentId int64, contentTitle string, commentId int64, commentDescribe string, commentAnonymous bool) error {
	return CommentAbout(userId, receiveUserId, contentId, contentTitle, commentId, commentDescribe, MessageTypeCommentForContent, commentAnonymous)
}

func CommentForComment(userId int64, receiveUserId int64, contentId int64, contentTitle string, commentId int64, commentDescribe string, commentAnonymous bool) error {
	return CommentAbout(userId, receiveUserId, contentId, contentTitle, commentId, commentDescribe, MessageTypeCommentForComment, commentAnonymous)
}

func FollowYou(userId int64, receiveUserId int64) error {
	m := new(Message)
	m.UserId = userId
	m.ReceiveUserId = receiveUserId
	m.MessageType = MessageTypeFollow
	return m.Insert()
}

func PublishContent(userId int64, receiveUserId int64, contentId int64, contentTitle string, PublishAgain bool) error {
	if receiveUserId != 0 {
		m := new(Message)
		m.UserId = userId
		m.ReceiveUserId = receiveUserId
		m.ContentId = contentId
		m.ContentTitle = contentTitle
		m.MessageType = MessageTypeContentPublish
		if PublishAgain {
			m.PublishAgain = 1
		}
		return m.Insert()
	}

	us := make([]Relation, 0)
	err := FaFaRdb.Client.Where("user_b_id=?", userId).Cols("user_a_id").Find(&us)
	if err != nil {
		return err
	}

	for _, u := range us {
		m := new(Message)
		m.UserId = userId
		m.ReceiveUserId = u.UserAId
		m.ContentId = contentId
		m.ContentTitle = contentTitle
		m.MessageType = MessageTypeContentPublish
		if PublishAgain {
			m.PublishAgain = 1
		}
		err = m.Insert()
		if err != nil {
			return err
		}
	}

	return nil
}

func Private(sendUserId, receiveUserId int64, sendMessage string) error {
	m := new(Message)
	m.MessageType = MessageTypePrivate
	m.SendUserId = sendUserId
	m.ReceiveUserId = receiveUserId
	m.SendMessage = sendMessage
	return m.Insert()
}

func (m *Message) Update() error {
	if m.Id == 0 || m.ReceiveUserId == 0 {
		return errors.New("where is empty")
	}

	sess := FaFaRdb.Client.Where("id=?", m.Id).And("receive_user_id=?", m.ReceiveUserId).And("receive_status!=?", 2).And("receive_status!=?", m.ReceiveStatus).Cols("receive_status")
	if m.ReceiveStatus == 1 {
		m.ReadTime = time.Now().Unix()
		sess.Cols("read_time")
	} else if m.ReceiveStatus == 2 {
		m.DeleteTime = time.Now().Unix()
		sess.Cols("delete_time")
	}

	_, err := sess.Update(m)
	return err
}

func GroupCount(userId int64) (countMap map[string]int, err error) {
	countMap = make(map[string]int)
	andSQL := ""
	if userId != 0 {
		andSQL = fmt.Sprintf("and receive_user_id=%d", userId)
	}

	sql := fmt.Sprintf("SELECT message_type,count(message_type) as count FROM `fafacms_message` WHERE receive_status=0 %s group by message_type;", andSQL)
	result, err := FaFaRdb.Client.QueryString(sql)
	if err != nil {
		return
	}

	for _, v := range result {
		i, err1 := util.SI(v["count"])
		if err1 != nil {
			continue
		}
		countMap[v["message_type"]] = i
	}
	return
}
