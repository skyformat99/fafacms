package model

import (
	"errors"
	"time"
)

var HistoryRecord = true

type Content struct {
	Id               int64  `json:"id" xorm:"bigint pk autoincr"`
	Seo              string `json:"seo" xorm:"index"`
	Title            string `json:"title" xorm:"varchar(200) notnull"`
	PreTitle         string `json:"pre_title" xorm:"varchar(200) notnull"`
	UserId           int64  `json:"user_id" xorm:"bigint index"`
	UserName         string `json:"user_name" xorm:"index"`
	NodeId           int64  `json:"node_id" xorm:"bigint index"`
	NodeSeo          string `json:"node_seo" xorm:"index"`
	Status           int    `json:"status" xorm:"not null comment('0 normal, 1 hide，2 ban, 3 rubbish') TINYINT(1) index"`
	Top              int    `json:"top" xorm:"not null comment('0 normal, 1 top') TINYINT(1) index"`
	Describe         string `json:"describe" xorm:"TEXT"`
	PreDescribe      string `json:"pre_describe" xorm:"TEXT"`
	PreFlush         int    `json:"pre_flush" xorm:"not null comment('1 flush') TINYINT(1)"`
	CloseComment     int    `json:"close_comment" xorm:"not null comment('0 close, 1 open, 2 direct open') TINYINT(1)"`
	Version          int    `json:"version"`
	CreateTime       int64  `json:"create_time"`
	UpdateTime       int64  `json:"update_time,omitempty"`
	FirstPublishTime int64  `json:"first_publish_time,omitempty"`
	PublishTime      int64  `json:"publish_time,omitempty"`
	ImagePath        string `json:"image_path" xorm:"varchar(700)"`
	Views            int64  `json:"views"`
	Password         string `json:"password,omitempty"`
	SortNum          int64  `json:"sort_num"`
	Bad              int64  `json:"bad"`
	Cool             int64  `json:"cool"`
	CommentNum       int64  `json:"comment_num"`
}

var ContentSortName = []string{"=id", "-user_id", "-top", "+sort_num", "-first_publish_time", "-publish_time", "-create_time", "-update_time", "-views", "=version", "+status", "=seo"}
var ContentSortName2 = []string{
	"=id",
	"-user_id",
	"-top",
	"+sort_num",
	"-first_publish_time",
	"-publish_time",
	"-views",
	"=seo",}

type ContentCool struct {
	Id         int64 `json:"id" xorm:"bigint pk autoincr"`
	UserId     int64 `json:"user_id" xorm:"bigint index"`
	ContentId  int64 `json:"content_id,omitempty" xorm:"bigint index"`
	CreateTime int64 `json:"create_time"`
}

type ContentBad struct {
	Id         int64 `json:"id" xorm:"bigint pk autoincr"`
	UserId     int64 `json:"user_id" xorm:"bigint index"`
	ContentId  int64 `json:"content_id,omitempty" xorm:"bigint index"`
	CreateTime int64 `json:"create_time"`
}

type ContentHistory struct {
	Id         int64  `json:"id" xorm:"bigint pk autoincr"`
	ContentId  int64  `json:"content_id" xorm:"bigint index"`
	Title      string `json:"title" xorm:"varchar(200) notnull"`
	UserId     int64  `json:"user_id" xorm:"bigint index"`
	NodeId     int64  `json:"node_id" xorm:"bigint index"`
	Describe   string `json:"describe" xorm:"TEXT"`
	Types      int    `json:"types" xorm:"not null comment('0 update save, 1 publish, 2 restore') TINYINT(1)"`
	CreateTime int64  `json:"create_time"`
}

var ContentHistorySortName = []string{"=id", "-user_id", "-create_time", "-content_id"}

func (c *Content) CountNumUnderNode() (int64, error) {
	if c.UserId == 0 || c.NodeId == 0 {
		return 0, errors.New("where is empty")
	}

	allNum, err := FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("node_id=?", c.NodeId).Count()
	if err != nil {
		return 0, err
	}
	return allNum, nil
}

func (c *Content) CheckSeoValid() (bool, error) {
	if c.UserId == 0 || c.Seo == "" {
		return false, errors.New("where is empty")
	}

	num, err := FafaRdb.Client.Table(c).Where("user_id=?", c.UserId).And("seo=?", c.Seo).Count()

	if num >= 1 {
		return true, nil
	}
	return false, err
}

func (c *Content) Insert() (int64, error) {
	c.CreateTime = time.Now().Unix()
	return FafaRdb.InsertOne(c)
}

func (c *Content) Get() (bool, error) {
	if c.Id == 0 {
		return false, errors.New("where is empty")
	}

	return FafaRdb.Client.Get(c)
}

func (c *Content) GetByRaw() (bool, error) {
	return FafaRdb.Client.Get(c)
}

func (c *Content) UpdateDescribeAndHistory(save bool) error {
	if c.UserId == 0 || c.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()

	var err error

	now := time.Now().Unix()
	if HistoryRecord && save && c.PreFlush != 1 {
		history := new(ContentHistory)
		history.NodeId = c.NodeId
		history.CreateTime = now

		history.Title = c.PreTitle
		history.Describe = c.PreDescribe
		history.ContentId = c.Id
		history.UserId = c.UserId

		history.Types = 0
		_, err = session.InsertOne(history)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	c.UpdateTime = now
	c.PreDescribe = c.Describe
	c.PreTitle = c.Title
	c.PreFlush = 0

	_, err = session.Cols("update_time", "pre_title", "pre_describe", "pre_flush").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
	}

	return err
}

func (c *Content) UpdateSeo() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	return FafaRdb.Client.Cols("seo").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

func (c *Content) UpdateImage() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	return FafaRdb.Client.Cols("image_path").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

func (c *Content) UpdateStatus() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	return FafaRdb.Client.Cols("status").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

func (c *Content) UpdateTop() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	return FafaRdb.Client.Cols("top").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

// update comment
func (c *Content) UpdateComment() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	return FafaRdb.Client.Cols("close_comment").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

// update password
func (c *Content) UpdatePassword() (int64, error) {
	if c.UserId == 0 || c.Id == 0 {
		return 0, errors.New("where is empty")
	}
	return FafaRdb.Client.Cols("password").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
}

func (n *Content) UpdateNode(beforeNodeId int64) error {
	if n.UserId == 0 || n.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}

	_, err = session.Exec("update fafacms_content SET sort_num=sort_num-1 where sort_num > ? and user_id = ? and node_id = ?", n.SortNum, n.UserId, beforeNodeId)
	if err != nil {
		session.Rollback()
		return err
	}

	c, err := session.Table(n).Where("user_id=?", n.UserId).And("node_id=?", n.NodeId).Count()
	if err != nil {
		session.Rollback()
		return err
	}

	n.SortNum = c

	_, err = session.Exec("update fafacms_content SET sort_num=?, node_id=?, node_seo=? where id = ? and user_id = ?", n.SortNum, n.NodeId, n.NodeSeo, n.Id, n.UserId)
	if err != nil {
		session.Rollback()
		return err
	}

	_, err = session.Exec("update fafacms_content_history SET node_id=? where content_id = ? and user_id = ?", n.NodeId, n.Id, n.UserId)
	if err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		return err
	}

	return nil
}

func (c *Content) UpdateView() {
	FafaRdb.Client.ID(c.Id).Incr("views").Update(new(Content))
}

// get brother content
func (c *Content) GetBrotherContent() (pre, next *Content, err error) {
	pre = new(Content)
	next = new(Content)
	session1 := FafaRdb.Client.Where("sort_num > ?", c.SortNum).And("id!=?", c.Id).And("node_id=?", c.NodeId).Omit("describe", "pre_describe")
	session1.And("status=?", 0).And("version>?", 0).Asc("sort_num").Limit(1)
	_, err = session1.Get(next)
	if err != nil {
		return
	}

	session2 := FafaRdb.Client.Where("sort_num < ?", c.SortNum).And("id!=?", c.Id).And("node_id=?", c.NodeId).Omit("describe", "pre_describe")
	session2.And("status=?", 0).And("version>?", 0).Desc("sort_num").Limit(1)
	_, err = session2.Get(pre)
	if err != nil {
		return
	}

	return
}

func (c *Content) PublishDescribe() error {
	if c.UserId == 0 || c.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	now := time.Now().Unix()
	if HistoryRecord {
		history := new(ContentHistory)
		history.NodeId = c.NodeId
		history.CreateTime = now

		history.Title = c.PreTitle
		history.Describe = c.PreDescribe
		history.ContentId = c.Id
		history.UserId = c.UserId

		history.Types = 1
		_, err := session.InsertOne(history)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	if c.Version == 0 {
		c.FirstPublishTime = now
	}

	c.Version = c.Version + 1
	c.UpdateTime = now
	c.PreFlush = 1
	c.Title = c.PreTitle
	c.Describe = c.PreDescribe
	c.PublishTime = now
	_, err := session.Cols("title", "describe", "pre_flush", "update_time", "publish_time", "first_publish_time", "version").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
	if err != nil {
		session.Rollback()
		return err
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		return err
	}
	return nil
}

func (c *Content) ResetDescribe(save bool) error {
	if c.UserId == 0 || c.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	now := time.Now().Unix()
	if c.PreFlush != 1 && HistoryRecord && save {
		history := new(ContentHistory)
		history.NodeId = c.NodeId
		history.CreateTime = now
		history.Title = c.PreTitle
		history.Describe = c.PreDescribe
		history.ContentId = c.Id
		history.UserId = c.UserId

		history.Types = 2
		_, err := session.InsertOne(history)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	c.UpdateTime = now
	c.PreFlush = 0
	c.PreTitle = c.Title
	c.PreDescribe = c.Describe
	_, err := session.Cols("pre_title", "pre_describe", "pre_flush", "update_time").Where("id=?", c.Id).And("user_id=?", c.UserId).Update(c)
	if err != nil {
		session.Rollback()
		return err
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		return err
	}
	return nil
}

func (c *Content) Delete() error {
	if c.UserId == 0 || c.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	if _, err := session.Where("id=?", c.Id).And("user_id=?", c.UserId).Delete(new(Content)); err != nil {
		session.Rollback()
		return err
	}

	if _, err := session.Where("content_id=?", c.Id).And("user_id=?", c.UserId).Delete(new(ContentHistory)); err != nil {
		session.Rollback()
		return err
	}

	if err := session.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *ContentHistory) GetRaw() (bool, error) {
	return FafaRdb.Client.Get(c)
}

func (c *ContentHistory) Delete() (int64, error) {
	return FafaRdb.Client.Where("id=?", c.Id).And("user_id=?", c.UserId).Delete(new(ContentHistory))
}
