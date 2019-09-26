package model

import (
	"errors"
)

// a love b
type Relation struct {
	Id         int64  `json:"id" xorm:"bigint pk autoincr"`
	CreateTime int64  `json:"create_time"`
	UserAId    int64  `json:"user_a_id" xorm:"index"`
	UserBId    int64  `json:"user_b_id" xorm:"index"`
	UserAName  string `json:"user_a_name" xorm:"index"`
	UserBName  string `json:"user_b_name" xorm:"index"`
}

func (r *Relation) Add() (err error) {
	if r.UserAId == 0 || r.UserBId == 0 {
		return errors.New("user id empty")
	}

	num, err := FaFaRdb.Client.Where("user_a_id=?", r.UserAId).And("user_b_id=?", r.UserBId).Count(new(Relation))
	if err != nil {
		return err
	}

	if num > 0 {
		return
	}

	se := FaFaRdb.Client.NewSession()
	defer se.Close()
	err = se.Begin()
	if err != nil {
		return
	}

	return
}
