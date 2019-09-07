package model

import (
	"errors"
	"time"
)

// 内容节点，最多两层
type ContentNode struct {
	Id           int    `json:"id" xorm:"bigint pk autoincr"`
	UserId       int    `json:"user_id" xorm:"bigint index"`
	UserName     string `json:"user_name" xorm:"index"`
	Seo          string `json:"seo" xorm:"index"`
	Status       int    `json:"status" xorm:"not null comment('0 normal,1 hide') TINYINT(1) index"`
	Name         string `json:"name" xorm:"varchar(100) notnull"`
	Describe     string `json:"describe" xorm:"TEXT"`
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time,omitempty"`
	ImagePath    string `json:"image_path" xorm:"varchar(700)"`
	ParentNodeId int    `json:"parent_node_id" xorm:"bigint"`
	Level        int    `json:"level"`
	SortNum      int    `json:"sort_num"` //  排序，数字越大排越后
}

// 内容节点排序专用
var ContentNodeSortName = []string{"=id", "+sort_num", "-create_time", "-update_time", "+status", "=seo"}

// 检查节点数量
func (n *ContentNode) CountNodeNum() (int, error) {
	if n.UserId == 0 {
		return 0, errors.New("where is empty")
	}

	// 创建时，要在同一层排序最大，这样排在最后
	num, err := FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).Where("parent_node_id=?", n.ParentNodeId).Count()
	return int(num), err
}

// 节点检查SEO的子路径是否有效
func (n *ContentNode) CheckSeoValid() (bool, error) {
	// 用户ID和SEO必须同时存在
	if n.UserId == 0 || n.Seo == "" {
		return false, errors.New("where is empty")
	}

	// 常规统计
	c, err := FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("seo=?", n.Seo).Count()

	// 如果大于1表示存在
	if c >= 1 {
		return true, nil
	}
	return false, err
}

// 节点检查 指定的父亲节点是否存在
func (n *ContentNode) CheckParentValid() (bool, error) {
	if n.UserId == 0 || n.ParentNodeId == 0 {
		return false, errors.New("where is empty")
	}

	// 只允许两层节点，Level必须为0
	c, err := FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("id=?", n.ParentNodeId).And("level=?", 0).Count()

	// 如果大于1表示存在
	if c >= 1 {
		return true, nil
	}
	return false, err
}

// 检查节点下的儿子节点数量
func (n *ContentNode) CheckChildrenNum() (int, error) {
	if n.UserId == 0 || n.Id == 0 {
		return 0, errors.New("where is empty")
	}
	num, err := FafaRdb.Client.Table(n).Where("user_id=?", n.UserId).And("parent_node_id=?", n.Id).Count()
	return int(num), err
}

// 节点常规插入
func (n *ContentNode) InsertOne() error {
	n.CreateTime = time.Now().Unix()
	_, err := FafaRdb.Insert(n)
	return err
}

// 节点常规获取，ID和SEO必须存在一者
func (n *ContentNode) Get() (bool, error) {
	if n.Id == 0 && n.Seo == "" {
		return false, errors.New("where is empty")
	}
	return FafaRdb.Client.Get(n)
}

// 获取某用户一个sort的节点
func (n *ContentNode) GetSortOneNode() (bool, error) {
	if n.UserId == 0 {
		return false, errors.New("where is empty")
	}
	return FafaRdb.Client.Get(n)
}

// 判断节点是否存在
func (n *ContentNode) Exist() (bool, error) {
	if n.Id == 0 {
		return false, errors.New("where is empty")
	}
	num, err := FafaRdb.Client.Table(n).Where("id=?", n.Id).And("user_id=?", n.UserId).Count()
	if err != nil {
		return false, err
	}

	return num >= 1, nil
}

// 更新节点SEO
func (n *ContentNode) UpdateSeo() error {
	if n.UserId == 0 || n.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}

	// 文章内容的节点SEO也要更改
	_, err = session.Exec("update fafacms_content SET node_seo=? where user_id=? and node_id=?", n.Seo, n.UserId, n.Id)
	if err != nil {
		session.Rollback()
		return err
	}

	n.UpdateTime = time.Now().Unix()
	_, err = session.Where("id=?", n.Id).And("user_id=?", n.UserId).Cols("seo", "update_time").Update(n)
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

// 更新节点详情
// 这种 update 要自己控制好，不要有多余字段被更新了
func (n *ContentNode) UpdateInfo() error {
	if n.UserId == 0 || n.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	n.UpdateTime = time.Now().Unix()
	_, err := session.Where("id=?", n.Id).And("user_id=?", n.UserId).MustCols("describe").Omit("id", "user_id").Update(n)
	return err
}

// 更新节点图片
func (n *ContentNode) UpdateImage() error {
	if n.UserId == 0 || n.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	n.UpdateTime = time.Now().Unix()
	_, err := session.Where("id=?", n.Id).And("user_id=?", n.UserId).MustCols("image_path").Omit("id", "user_id").Update(n)
	return err
}

// 更新节点状态，1表示隐藏
func (n *ContentNode) UpdateStatus() error {
	if n.UserId == 0 || n.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	n.UpdateTime = time.Now().Unix()
	_, err := session.Where("id=?", n.Id).And("user_id=?", n.UserId).Cols("status", "update_time").Update(n)
	return err
}

// 更新节点的父亲
func (n *ContentNode) UpdateParent(beforeParentNode int) error {
	if n.UserId == 0 || n.Id == 0 {
		return errors.New("where is empty")
	}

	session := FafaRdb.Client.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return err
	}

	// 同一层的先减一，假装删除这一个节点
	_, err = session.Exec("update fafacms_content_node SET sort_num=sort_num-1 where sort_num > ? and user_id = ? and parent_node_id = ?", n.SortNum, n.UserId, beforeParentNode)
	if err != nil {
		session.Rollback()
		return err
	}

	// 更新节点时，排序在该层最大
	c, err := session.Table(n).Where("user_id=?", n.UserId).And("parent_node_id", n.ParentNodeId).Count()
	if err != nil {
		session.Rollback()
		return err
	}

	n.SortNum = int(c)
	n.UpdateTime = time.Now().Unix()

	// 更新节点
	// 事务怕本ORM混淆，所以直接使用原生
	// 每次更改节点，他都会成为这一层最靓丽排得最前面的仔，不是啦，成为排最后的仔，排序越大越往后。
	_, err = session.Exec("update fafacms_content_node SET sort_num=?, update_time=?, level=?, parent_node_id=? where id = ? and user_id = ?", n.SortNum, n.UpdateTime, n.Level, n.ParentNodeId, n.Id, n.UserId)
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
