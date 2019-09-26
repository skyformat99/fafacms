package model

// a love b
type Relation struct {
	Id         int64  `json:"id" xorm:"bigint pk autoincr"`
	CreateTime int64  `json:"create_time"`
	UserAId    int64  `json:"user_a_id" xorm:"index"`
	UserBId    int64  `json:"user_b_id" xorm:"index"`
	UserAName  string `json:"user_a_name" xorm:"index"`
	UserBName  string `json:"user_b_name" xorm:"index"`
	IsBoth     int    `json:"is_both" xorm:"index"`
}

func (r *Relation) Get() {

}
