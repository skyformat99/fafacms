package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/util"
	"math"
	"time"
)

var TimeZone int64 = 0

// Local time format you know
func GetSecond2DateTimes(second int64) string {
	second = second + 3600*TimeZone
	tm := time.Unix(second, 0)
	return tm.UTC().Format("2006-01-02 15:04:05")

}

func Home(c *gin.Context) {
	resp := new(Resp)
	resp.Flag = true
	resp.Data = "FaFa CMS: https://github.com/hunterhug/fafacms Version:" + config.Version
	defer func() {
		c.JSON(200, resp)
	}()
}

type People struct {
	Id                    int64  `json:"id"`
	Name                  string `json:"name"`
	NickName              string `json:"nick_name"`
	Email                 string `json:"email"`
	WeChat                string `json:"wechat"`
	WeiBo                 string `json:"weibo"`
	Github                string `json:"github"`
	QQ                    string `json:"qq"`
	Gender                int    `json:"gender"`
	Describe              string `json:"describe"`
	ShortDescribe         string `json:"short_describe"`
	HeadPhoto             string `json:"head_photo"`
	CreateTime            string `json:"create_time"`
	CreateTimeInt         int64  `json:"create_time_int"`
	UpdateTime            string `json:"update_time,omitempty"`
	UpdateTimeInt         int64  `json:"update_time_int,omitempty"`
	ActivateTime          string `json:"activate_time,omitempty"`
	ActivateTimeInt       int64  `json:"activate_time_int,omitempty"`
	LoginTime             string `json:"login_time,omitempty"`
	LoginTimeInt          int64  `json:"login_time_int,omitempty"`
	LoginIp               string `json:"login_ip,omitempty"`
	NickNameUpdateTimeInt int64  `json:"nick_name_update_time,omitempty"`
	NickNameUpdateTime    string `json:"nick_name_update_time,omitempty"`
	IsInBlack             bool   `json:"is_in_black"`
	IsVip                 bool   `json:"is_vip"`
	FollowedNum           int64  `json:"followed_num"`
	FollowingNum          int64  `json:"following_num"`
	ContentNum            int64  `json:"content_num"`      // normal publish content num
	ContentCoolNum        int64  `json:"content_cool_num"` // normal content cool num
}

type PeoplesRequest struct {
	Vip  int      `json:"vip" validate:"oneof=-1 0 1"`
	Sort []string `json:"sort"`
	PageHelp
}

type PeoplesResponse struct {
	Users []People `json:"users"`
	PageHelp
}

func Peoples(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	respResult := new(PeoplesResponse)
	req := new(PeoplesRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("Peoples err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	session := model.FaFaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.User)).Where("1=1").And("status!=?", 0).And("name!=?", "admin")

	if req.Vip != -1 {
		if req.Vip == 0 {
			session.And("vip=?", 0)
		} else {
			session.And("vip=?", 1)
		}
	}

	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("Peoples err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	users := make([]model.User, 0)
	p := &req.PageHelp
	if total == 0 {
		if p.Limit == 0 {
			p.Limit = 20
		}
	} else {
		p.build(session, req.Sort, model.UserSortName)
		err = session.Find(&users)
		if err != nil {
			flog.Log.Errorf("Peoples err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	peoples := make([]People, 0, len(users))
	for _, v := range users {
		p := People{}
		p.Id = v.Id
		p.ShortDescribe = v.ShortDescribe
		p.Describe = v.Describe
		p.CreateTimeInt = v.CreateTime
		p.CreateTime = GetSecond2DateTimes(v.CreateTime)

		p.UpdateTimeInt = v.UpdateTime
		if v.UpdateTime > 0 {
			p.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
		}

		p.ActivateTimeInt = v.ActivateTime
		if v.ActivateTime > 0 {
			p.ActivateTime = GetSecond2DateTimes(v.ActivateTime)
		}

		p.LoginTimeInt = v.LoginTime
		if v.LoginTime > 0 {
			p.LoginTime = GetSecond2DateTimes(v.LoginTime)
		}

		if v.Status == 2 {
			p.IsInBlack = true
		}
		p.Email = v.Email
		p.Github = v.Github
		p.Name = v.Name
		p.NickName = v.NickName
		p.HeadPhoto = v.HeadPhoto
		p.QQ = v.QQ
		p.WeChat = v.WeChat
		p.WeiBo = v.WeiBo
		p.Gender = v.Gender
		p.IsVip = v.Vip == 1
		p.FollowedNum = v.FollowedNum
		p.FollowingNum = v.FollowingNum
		p.ContentNum = v.ContentNum
		p.ContentCoolNum = v.ContentCoolNum
		peoples = append(peoples, p)
	}
	respResult.Users = peoples
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	p.Total = int(total)
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type Node struct {
	Id            int64  `json:"id"`
	Seo           string `json:"seo"`
	Name          string `json:"name"`
	Describe      string `json:"describe"`
	ImagePath     string `json:"image_path"`
	CreateTime    string `json:"create_time"`
	CreateTimeInt int64  `json:"create_time_int"`
	UpdateTime    string `json:"update_time,omitempty"`
	UpdateTimeInt int64  `json:"update_time_int,omitempty"`
	UserId        int64  `json:"user_id"`
	UserName      string `json:"user_name"`
	SortNum       int64  `json:"sort_num"`
	Level         int    `json:"level"`
	Status        int    `json:"status"`
	ParentNodeId  int64  `json:"parent_node_id"`
	Son           []Node `json:"son,omitempty"`
	ContentNum    int64  `json:"content_num"` // normal publish content num
}

type NodesInfoRequest struct {
	UserId   int64    `json:"user_id"`
	UserName string   `json:"user_name"`
	Sort     []string `json:"sort"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}

func NodesInfo(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	respResult := new(NodesResponse)
	req := new(NodesInfoRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.UserId == 0 && req.UserName == "" {
		flog.Log.Errorf("NodesInfo err:%s", "")
		resp.Error = Error(ParasError, "user info empty")
		return
	}

	session := model.FaFaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.ContentNode)).Where("1=1").And("status=?", 0)

	if req.UserId != 0 {
		session.And("user_id=?", req.UserId)
	}

	if req.UserName != "" {
		session.And("user_name=?", req.UserName)
	}

	nodes := make([]model.ContentNode, 0)
	Build(session, req.Sort, model.ContentNodeSortName)
	err := session.Find(&nodes)
	if err != nil {
		flog.Log.Errorf("NodesInfo err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	father := make([]model.ContentNode, 0)
	son := make([]model.ContentNode, 0)
	for _, v := range nodes {
		if v.Level == 0 {
			father = append(father, v)
		} else {
			son = append(son, v)
		}
	}

	n := make([]Node, 0)
	for _, v := range father {
		f := Node{}
		f.Id = v.Id
		f.Seo = v.Seo
		f.Describe = v.Describe
		f.ImagePath = v.ImagePath
		f.Name = v.Name
		if v.UpdateTime > 0 {
			f.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
			f.UpdateTimeInt = v.UpdateTime
		}
		f.CreateTime = GetSecond2DateTimes(v.CreateTime)
		f.CreateTimeInt = v.CreateTime
		f.SortNum = v.SortNum
		f.UserName = v.UserName
		f.UserId = v.UserId
		f.Level = v.Level
		f.ParentNodeId = v.ParentNodeId
		f.ContentNum = v.ContentNum
		for _, vv := range son {
			if vv.ParentNodeId == f.Id {
				s := Node{}
				s.Id = vv.Id
				s.Seo = vv.Seo
				s.Describe = vv.Describe
				s.ImagePath = vv.ImagePath
				s.Name = vv.Name
				if vv.UpdateTime > 0 {
					s.UpdateTimeInt = vv.UpdateTime
					s.UpdateTime = GetSecond2DateTimes(vv.UpdateTime)
				}
				s.CreateTime = GetSecond2DateTimes(vv.CreateTime)
				s.CreateTimeInt = vv.CreateTime
				s.SortNum = vv.SortNum
				s.UserId = vv.UserId
				s.UserName = vv.UserName
				s.Level = vv.Level
				s.ContentNum = vv.ContentNum
				s.ParentNodeId = vv.ParentNodeId
				f.Son = append(f.Son, s)
			}
		}

		n = append(n, f)
	}

	respResult.Nodes = n
	resp.Flag = true
	resp.Data = respResult
}

type NodeInfoRequest struct {
	Id       int    `json:"id"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Seo      string `json:"seo"`
	ListSon  bool   `json:"list_son"`
}

func NodeInfo(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	req := new(NodeInfoRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.Id == 0 && req.Seo == "" {
		flog.Log.Errorf("NodeInfo err: %s", "content node id or seo empty")
		resp.Error = Error(ParasError, "content node id or seo empty")
		return
	}

	if req.Id == 0 && req.Seo != "" {
		if req.UserId == 0 && req.UserName == "" {
			flog.Log.Errorf("NodeInfo err: %s", "content node seo exist but user info empty")
			resp.Error = Error(ParasError, "content node seo exist but user info empty")
			return
		}
	}

	session := model.FaFaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.ContentNode)).Where("1=1").And("status=?", 0)

	if req.UserId != 0 {
		session.And("user_id=?", req.UserId)
	}

	if req.UserName != "" {
		session.And("user_name=?", req.UserName)
	}

	if req.Id != 0 {
		session.And("id=?", req.Id)
	}

	if req.Seo != "" {
		session.And("seo=?", req.Seo)
	}

	v := new(model.ContentNode)
	exist, err := session.Get(v)
	if err != nil {
		flog.Log.Errorf("NodeInfo err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("NodeInfo err:%s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	f := Node{}
	f.Id = v.Id
	f.Seo = v.Seo
	f.Describe = v.Describe
	f.ImagePath = v.ImagePath
	f.Name = v.Name
	if v.UpdateTime > 0 {
		f.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
		f.UpdateTimeInt = v.UpdateTime
	}
	f.CreateTime = GetSecond2DateTimes(v.CreateTime)
	f.CreateTimeInt = v.CreateTime
	f.SortNum = v.SortNum
	f.UserName = v.UserName
	f.UserId = v.UserId
	f.Level = v.Level
	f.ParentNodeId = v.ParentNodeId
	f.ContentNum = v.ContentNum

	// 是顶层且需要列出儿子
	if f.Level == 0 && req.ListSon {
		ns := make([]model.ContentNode, 0)
		err = model.FaFaRdb.Client.Where("parent_node_id=?", f.Id).And("status=?", 0).Find(&ns)
		if err != nil {
			flog.Log.Errorf("NodeInfo err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		for _, vv := range ns {
			ff := Node{}
			ff.Id = vv.Id
			ff.Seo = vv.Seo
			ff.Describe = vv.Describe
			ff.ImagePath = vv.ImagePath
			ff.Name = vv.Name
			if vv.UpdateTime > 0 {
				ff.UpdateTime = GetSecond2DateTimes(vv.UpdateTime)
				ff.UpdateTimeInt = vv.UpdateTime
			}
			ff.CreateTime = GetSecond2DateTimes(vv.CreateTime)
			ff.CreateTimeInt = vv.CreateTime
			ff.SortNum = vv.SortNum
			ff.UserName = vv.UserName
			ff.UserId = vv.UserId
			ff.Level = vv.Level
			ff.ParentNodeId = vv.ParentNodeId
			ff.ContentNum = vv.ContentNum
			f.Son = append(f.Son, ff)
		}
	}
	resp.Flag = true
	resp.Data = f
}

type UserInfoRequest struct {
	Id   int64  `json:"user_id"`
	Name string `json:"user_name"`
}

func UserInfo(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	req := new(UserInfoRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.Id == 0 && req.Name == "" {
		resp.Error = Error(ParasError, "where is empty")
		return
	}

	user := new(model.User)
	user.Id = req.Id
	user.Name = req.Name
	exist, err := model.FaFaRdb.Client.Where("status!=?", 0).Get(user)
	if err != nil {
		flog.Log.Errorf("UserInfo err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("UserInfo err:%s", "user  not found")
		resp.Error = Error(UserNotFound, "")
		return
	}

	v := user
	p := People{}
	p.Id = v.Id
	p.Describe = v.Describe
	p.ShortDescribe = v.ShortDescribe
	p.CreateTime = GetSecond2DateTimes(v.CreateTime)
	p.CreateTimeInt = v.CreateTime

	if v.Status == 2 {
		p.IsInBlack = true
	}

	p.UpdateTimeInt = v.UpdateTime
	if v.UpdateTime > 0 {
		p.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
	}

	p.LoginTimeInt = v.LoginTime
	if v.LoginTime > 0 {
		p.LoginTime = GetSecond2DateTimes(v.LoginTime)
	}

	p.ActivateTimeInt = v.ActivateTime
	if v.ActivateTime > 0 {
		p.ActivateTime = GetSecond2DateTimes(v.ActivateTime)
	}
	p.Email = v.Email
	p.Github = v.Github
	p.Name = v.Name
	p.NickName = v.NickName
	p.HeadPhoto = v.HeadPhoto
	p.QQ = v.QQ
	p.WeChat = v.WeChat
	p.WeiBo = v.WeiBo
	p.Gender = v.Gender
	p.FollowingNum = v.FollowingNum
	p.FollowedNum = v.FollowedNum
	p.ContentNum = v.ContentNum
	p.ContentCoolNum = v.ContentCoolNum
	p.IsVip = v.Vip == 1
	resp.Flag = true
	resp.Data = p
}

type UserCountRequest struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
}

type UserCountX struct {
	Count           int    `json:"count"`
	Days            string `json:"days"`
	CreateTimeBegin int64  `json:"first_publish_time_begin"`
	CreateTimeEnd   int64  `json:"first_publish_time_end"`
}
type UserCountResponse struct {
	Info     []UserCountX `json:"info"`
	UserId   int64        `json:"user_id"`
	UserName string       `json:"user_name"`
}

func UserCount(c *gin.Context) {
	resp := new(Resp)

	defer func() {
		JSON(c, 200, resp)
	}()

	req := new(UserCountRequest)
	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if req.UserId == 0 && req.UserName == "" {
		resp.Error = Error(ParasError, "where is empty")
		return
	}

	user := new(model.User)
	user.Id = req.UserId
	user.Name = req.UserName
	user.Status = 1
	exist, err := user.GetRaw()
	if err != nil {
		flog.Log.Errorf("UserCount err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("UserCount err:%s", "user not found")
		resp.Error = Error(UserNotFound, "")
		return
	}

	req.UserId = user.Id

	sql := fmt.Sprintf("SELECT DATE_FORMAT(from_unixtime(first_publish_time + %d * 3600)", TimeZone) + ",'%Y%m%d') as days,count(id) as count FROM `fafacms_content` WHERE first_publish_time!=0 and user_id=? and version>0 and status!=1 and status!=3 group by days;"
	result, err := model.FaFaRdb.Client.QueryString(sql, req.UserId)
	if err != nil {
		flog.Log.Errorf("UserCount err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	back := make([]UserCountX, 0)
	for _, v := range result {
		t := UserCountX{}
		t.Count, _ = util.SI(v["count"])
		t.Days = v["days"]
		begin, _ := time.ParseInLocation("20060102", t.Days, time.UTC)
		begin = begin.Add(time.Second * time.Duration(3600*TimeZone))
		end := begin.AddDate(0, 0, 1)
		t.CreateTimeBegin = begin.Unix()
		t.CreateTimeEnd = end.Unix()
		back = append(back, t)
	}

	resp.Flag = true
	resp.Data = UserCountResponse{
		Info:     back,
		UserId:   user.Id,
		UserName: user.Name,
	}
}

type ContentsRequest struct {
	NodeId                int64    `json:"node_id"`
	NodeSeo               string   `json:"node_seo"`
	UserId                int64    `json:"user_id"`
	UserName              string   `json:"user_name"`
	FirstPublishTimeBegin int64    `json:"first_publish_time_begin"`
	FirstPublishTimeEnd   int64    `json:"first_publish_time_end"`
	PublishTimeBegin      int64    `json:"publish_time_begin"`
	PublishTimeEnd        int64    `json:"publish_time_end"`
	Sort                  []string `json:"sort"`
	PageHelp
}

type ContentsX struct {
	Id                  int64      `json:"id"`
	Seo                 string     `json:"seo"`
	Title               string     `json:"title"`
	UserId              int64      `json:"user_id"`
	UserName            string     `json:"user_name"`
	NodeId              int64      `json:"node_id"`
	NodeSeo             string     `json:"node_seo"`
	Top                 int        `json:"top"`
	FirstPublishTime    string     `json:"first_publish_time"`
	PublishTime         string     `json:"publish_time,omitempty"`
	FirstPublishTimeInt int64      `json:"first_publish_time_int"`
	PublishTimeInt      int64      `json:"publish_time_int"`
	ImagePath           string     `json:"image_path"`
	Views               int64      `json:"views"`
	IsLock              bool       `json:"is_lock"`
	Describe            string     `json:"describe"`
	Next                *ContentsX `json:"next,omitempty"`
	Pre                 *ContentsX `json:"pre,omitempty"`
	SortNum             int64      `json:"sort_num"`
	Bad                 int64      `json:"bad"`
	Cool                int64      `json:"cool"`
	CommentNum          int64      `json:"comment_num"`
	IsBan               bool       `json:"is_ban"`
}

type ContentsResponse struct {
	Contents []ContentsX `json:"contents"`
	PageHelp
}

func Contents(c *gin.Context) {
	resp := new(Resp)

	respResult := new(ContentsResponse)
	req := new(ContentsRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("Contents err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := model.FaFaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Content)).Where("1=1")

	if req.UserId != 0 {
		session.And("user_id=?", req.UserId)
	}

	if req.UserName != "" {
		session.And("user_name=?", req.UserName)
	}

	session.And("status!=?", 1).And("status!=?", 3).And("version>?", 0)

	if req.NodeId != 0 {
		session.And("node_id=?", req.NodeId)
	}

	if req.NodeSeo != "" {
		session.And("node_seo=?", req.NodeSeo)
	}

	if req.FirstPublishTimeBegin > 0 {
		session.And("first_publish_time>=?", req.FirstPublishTimeBegin)
	}

	if req.FirstPublishTimeEnd > 0 {
		session.And("first_publish_time<?", req.FirstPublishTimeEnd)
	}

	if req.PublishTimeBegin > 0 {
		session.And("publish_time>=?", req.PublishTimeBegin)
	}

	if req.PublishTimeEnd > 0 {
		session.And("publish_time<?", req.PublishTimeEnd)
	}

	// count num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("Contents err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	cs := make([]model.Content, 0)
	p := &req.PageHelp
	if total == 0 {
		if p.Limit == 0 {
			p.Limit = 20
		}
	} else {
		// sql build
		p.build(session, req.Sort, model.ContentSortName2)
		// do query
		err = session.Omit("pre_describe", "pre_title").Find(&cs)
		if err != nil {
			flog.Log.Errorf("Contents err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	bcs := make([]ContentsX, 0, len(cs))
	for _, c := range cs {
		temp := ContentsX{}
		temp.UserId = c.UserId
		temp.Seo = c.Seo
		temp.SortNum = c.SortNum
		temp.NodeSeo = c.NodeSeo
		temp.UserName = c.UserName
		temp.Id = c.Id
		temp.Top = c.Top
		temp.Title = c.Title
		temp.NodeId = c.NodeId
		temp.Views = c.Views
		temp.ImagePath = c.ImagePath
		temp.FirstPublishTime = GetSecond2DateTimes(c.FirstPublishTime)
		temp.PublishTime = GetSecond2DateTimes(c.PublishTime)
		temp.FirstPublishTimeInt = c.FirstPublishTime
		temp.PublishTimeInt = c.PublishTime
		temp.CommentNum = c.CommentNum
		temp.Bad = c.Bad
		temp.Cool = c.Cool
		if c.Status == 2 {
			temp.IsBan = true
		}

		if c.Password != "" {
			temp.IsLock = true
		}

		if len(c.Describe) > 50 {
			temp.Describe = c.Describe[:50]
		} else {
			temp.Describe = c.Describe
		}
		bcs = append(bcs, temp)
	}

	respResult.Contents = bcs
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	p.Total = int(total)
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type ContentRequest struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Seo      string `json:"seo"`
	Password string `json:"password"`
	More     bool   `json:"more"`
}

func Content(c *gin.Context) {
	resp := new(Resp)
	req := new(ContentRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	var validate = validator.New()
	err := validate.Struct(req)
	if err != nil {
		flog.Log.Errorf("Content err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	if req.Id == 0 && req.Seo == "" {
		flog.Log.Errorf("Content err: %s", "content id or seo empty")
		resp.Error = Error(ParasError, "content id or seo empty")
		return
	}

	if req.Id == 0 && req.Seo != "" {
		if req.UserId == 0 && req.UserName == "" {
			flog.Log.Errorf("Content err: %s", "content seo exist but user info empty")
			resp.Error = Error(ParasError, "content seo exist but user info empty")
			return
		}
	}

	content := new(model.Content)
	content.Id = req.Id
	content.UserId = req.UserId
	content.UserName = req.UserName
	content.Seo = req.Seo
	exist, err := content.GetByRawAll()
	if err != nil {
		flog.Log.Errorf("Content err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("Content err: %s", "content not found")
		resp.Error = Error(ContentNotFound, "")
		return
	}

	if content.Status == 0 {

	} else if content.Status == 2 {
		flog.Log.Errorf("Content err: %s", "content ban")
		resp.Error = Error(ContentBanPermit, "")
		return
	} else {
		flog.Log.Errorf("Content err: %s", "content not found for it hide")
		resp.Error = Error(ContentNotFound, "")
		return
	}

	if content.Version == 0 {
		flog.Log.Errorf("Content err: %s", "content not found for it not publish")
		resp.Error = Error(ContentNotFound, "")
		return
	}

	if content.Password != "" && content.Password != req.Password {
		flog.Log.Errorf("Content err: %s", "content password")
		resp.Error = Error(ContentPasswordWrong, "")
		return
	}

	cx := content
	temp := ContentsX{}
	temp.UserId = cx.UserId
	temp.Seo = cx.Seo
	temp.NodeSeo = cx.NodeSeo
	temp.UserName = cx.UserName
	temp.Id = cx.Id
	temp.Top = cx.Top
	temp.Title = cx.Title
	temp.NodeId = cx.NodeId
	temp.Views = cx.Views
	temp.SortNum = cx.SortNum
	temp.FirstPublishTime = GetSecond2DateTimes(cx.FirstPublishTime)
	temp.PublishTime = GetSecond2DateTimes(cx.PublishTime)
	temp.FirstPublishTimeInt = cx.FirstPublishTime
	temp.PublishTimeInt = cx.PublishTime
	temp.ImagePath = cx.ImagePath
	temp.CommentNum = cx.CommentNum
	temp.Bad = cx.Bad
	temp.Cool = cx.Cool
	if cx.Password != "" {
		temp.IsLock = true
	}

	temp.Describe = cx.Describe

	cx.UpdateView()

	if req.More {
		cxx := new(model.Content)
		cxx.SortNum = cx.SortNum
		cxx.NodeId = cx.NodeId
		cxx.Id = cx.Id
		pre, next, err := cxx.GetBrotherContent()

		if err != nil {
			flog.Log.Errorf("Content err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		if pre.Id != 0 {
			temp1 := new(ContentsX)
			temp1.UserId = pre.UserId
			temp1.Seo = pre.Seo
			temp1.NodeSeo = pre.NodeSeo
			temp1.UserName = pre.UserName
			temp1.Id = pre.Id
			temp1.Top = pre.Top
			temp1.Title = pre.Title
			temp1.NodeId = pre.NodeId
			temp1.Views = pre.Views
			temp1.SortNum = pre.SortNum
			temp1.FirstPublishTime = GetSecond2DateTimes(pre.FirstPublishTime)
			temp1.PublishTime = GetSecond2DateTimes(pre.PublishTime)
			temp1.FirstPublishTimeInt = pre.FirstPublishTime
			temp1.PublishTimeInt = pre.PublishTime
			temp1.ImagePath = pre.ImagePath
			temp1.CommentNum = pre.CommentNum
			temp1.Bad = pre.Bad
			temp1.Cool = pre.Cool
			if pre.Password != "" {
				temp1.IsLock = true
			}
			temp.Pre = temp1
		}
		if next.Id != 0 {
			temp2 := new(ContentsX)
			temp2.UserId = next.UserId
			temp2.Seo = next.Seo
			temp2.NodeSeo = next.NodeSeo
			temp2.UserName = next.UserName
			temp2.Id = next.Id
			temp2.Top = next.Top
			temp2.Title = next.Title
			temp2.NodeId = next.NodeId
			temp2.Views = next.Views
			temp2.SortNum = next.SortNum
			temp2.FirstPublishTime = GetSecond2DateTimes(next.FirstPublishTime)
			temp2.PublishTime = GetSecond2DateTimes(next.PublishTime)
			temp2.FirstPublishTimeInt = next.FirstPublishTime
			temp2.PublishTimeInt = next.PublishTime
			temp2.ImagePath = next.ImagePath
			temp2.CommentNum = next.CommentNum
			temp2.Bad = next.Bad
			temp2.Cool = next.Cool
			if next.Password != "" {
				temp2.IsLock = true
			}
			temp.Next = temp2
		}
	}
	resp.Flag = true
	resp.Data = temp
}
