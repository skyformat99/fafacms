package model

import (
	"errors"
	"fmt"
	"github.com/hunterhug/fafacms/core/util"
	"time"
)

// 用户表
type User struct {
	Id                  int    `json:"id" xorm:"bigint pk autoincr"`
	Name                string `json:"name" xorm:"varchar(100) notnull unique"`  // 独一无二的标志
	NickName            string `json:"nick_name" xorm:"varchar(100) notnull"`    // 昵称，如小花花，随便改
	Email               string `json:"email" xorm:"varchar(100) notnull unique"` // 邮箱，独一无二
	WeChat              string `json:"wechat" xorm:"varchar(100)"`
	WeiBo               string `json:"weibo" xorm:"TEXT"`
	Github              string `json:"github" xorm:"TEXT"`
	QQ                  string `json:"qq" xorm:"varchar(100)"`
	Password            string `json:"password,omitempty" xorm:"varchar(100)"` // 明文的密码，就是这么强
	Gender              int    `json:"gender" xorm:"not null comment('0 unknow,1 boy,2 girl') TINYINT(1)"`
	Describe            string `json:"describe" xorm:"TEXT"`
	HeadPhoto           string `json:"head_photo" xorm:"varchar(700)"`
	CreateTime          int64  `json:"create_time"`
	UpdateTime          int64  `json:"update_time,omitempty"`
	ActivateTime        int64  `json:"activate_time,omitempty"`              // activate time
	ActivateCode        string `json:"activate_code,omitempty" xorm:"index"` // activate code
	ActivateCodeExpired int64  `json:"activate_code_expired,omitempty"`      // activate code expired time
	Status              int    `json:"status" xorm:"not null comment('0 unactive, 1 normal, 2 black') TINYINT(1) index"`
	GroupId             int    `json:"group_id,omitempty" xorm:"bigint index"`
	ResetCode           string `json:"reset_code,omitempty" xorm:"index"` // forget password code
	ResetCodeExpired    int64  `json:"reset_code_expired,omitempty"`      // forget password code expired
	LoginTime           int64  `json:"login_time,omitempty"`              // login time last time
	LoginIp             string `json:"login_ip,omitempty"`                // login ip last time
}

var UserSortName = []string{"=id", "=name", "-activate_time", "-create_time", "-update_time", "-gender"}

// 获取用户信息，不存在用户报错
func (u *User) Get() (err error) {
	var exist bool
	exist, err = FafaRdb.Client.Get(u)
	if err != nil {
		return
	}
	if !exist {
		return fmt.Errorf("user not found")
	}
	return
}

// 原生获取用户信息
func (u *User) GetRaw() (bool, error) {
	return FafaRdb.Client.Get(u)
}

func (u *User) Exist() (bool, error) {
	if u.Id == 0 && u.Name == "" && u.GroupId == 0 {
		return false, errors.New("where is empty")
	}

	s := FafaRdb.Client.Table(u)
	s.Where("1=1")

	if u.Id != 0 {
		s.And("id=?", u.Id)
	}

	if u.Name != "" {
		s.And("name=?", u.Name)
	}

	if u.GroupId != 0 {
		s.And("group_id=?", u.GroupId)
	}

	c, err := s.Count()

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (u *User) IsNameRepeat() (bool, error) {
	if u.Name == "" {
		return false, errors.New("where is empty")
	}
	c, err := FafaRdb.Client.Table(u).Where("name=?", u.Name).Count()

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (u *User) IsEmailRepeat() (bool, error) {
	if u.Email == "" {
		return false, errors.New("where is empty")
	}
	c, err := FafaRdb.Client.Table(u).Where("email=?", u.Email).Count()

	if c >= 1 {
		return true, nil
	}

	return false, err
}

func (u *User) InsertOne() error {
	u.CreateTime = time.Now().Unix()
	_, err := FafaRdb.Insert(u)
	return err
}

func (u *User) IsActivateCodeExist() (bool, error) {
	if u.ActivateCode == "" || u.Email == "" {
		return false, errors.New("where is empty")
	}
	c, err := FafaRdb.Client.Get(u)
	return c, err
}

func (u *User) UpdateActivateStatus() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}
	u.ActivateTime = time.Now().Unix()
	_, err := FafaRdb.Client.Where("id=?", u.Id).Cols("status", "activate_time").Update(u)
	return err
}

func (u *User) UpdateActivateCode() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}
	u.UpdateTime = time.Now().Unix()
	u.ActivateCode = util.GetGUID()
	u.ActivateCodeExpired = time.Now().Add(5 * time.Minute).Unix()
	_, err := FafaRdb.Client.Where("id=?", u.Id).Cols("activate_code", "activate_code_expired", "update_time").Update(u)
	return err
}

func (u *User) GetUserByEmail() (bool, error) {
	if u.Email == "" {
		return false, errors.New("where is empty")
	}
	c, err := FafaRdb.Client.Get(u)
	return c, err
}

func (u *User) UpdateCode() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}
	u.UpdateTime = time.Now().Unix()
	u.ResetCode = util.GetGUID()[0:6]
	u.ResetCodeExpired = time.Now().Unix() + 300
	_, err := FafaRdb.Client.Where("id=?", u.Id).Cols("reset_code", "reset_code_expired", "update_time").Update(u)
	return err
}

func (u *User) UpdatePassword() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}
	u.UpdateTime = time.Now().Unix()
	u.ResetCode = ""
	u.ResetCodeExpired = 0
	_, err := FafaRdb.Client.Where("id=?", u.Id).Cols("reset_code", "reset_code_expired", "update_time", "password").Update(u)
	return err
}

func (u *User) UpdateInfo() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}

	u.UpdateTime = time.Now().Unix()
	_, err := FafaRdb.Client.Where("id=?", u.Id).Omit("id").Update(u)
	return err
}

func (u *User) UpdateLoginInfo() error {
	if u.Id == 0 {
		return errors.New("where is empty")
	}

	_, err := FafaRdb.Client.Where("id=?", u.Id).Cols("login_time", "login_ip").Update(u)
	return err
}
