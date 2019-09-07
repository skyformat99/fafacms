package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/parrot/util"
	"math"
	"time"
)

var TimeZone int64 = 0

func GetSecond2DateTimes(secord int64) string {
	secord = secord + 3600*TimeZone
	tm := time.Unix(secord, 0)
	return tm.UTC().Format("2006-01-02 15:04:05")

}

func Home(c *gin.Context) {
	resp := new(Resp)
	resp.Flag = true
	resp.Data = "FaFa CMS: https://github.com/hunterhug/fafacms"
	defer func() {
		c.JSON(200, resp)
	}()
}

type People struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	NickName   string `json:"nick_name"`
	Email      string `json:"email"`
	WeChat     string `json:"wechat"`
	WeiBo      string `json:"weibo"`
	Github     string `json:"github"`
	QQ         string `json:"qq"`
	Gender     int    `json:"gender"`
	Describe   string `json:"describe"`
	HeadPhoto  string `json:"head_photo"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time,omitempty"`
}

type PeoplesRequest struct {
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

	session := model.FafaRdb.Client.NewSession()
	defer session.Close()

	// 找出激活的用户
	session.Table(new(model.User)).Where("1=1").And("status=?", 1)

	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		flog.Log.Errorf("ListUser err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	users := make([]model.User, 0)
	p := &req.PageHelp
	if total == 0 {
	} else {
		p.build(session, req.Sort, model.UserSortName)
		err = session.Find(&users)
		if err != nil {
			flog.Log.Errorf("ListUser err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	peoples := make([]People, 0, len(users))
	for _, v := range users {
		p := People{}
		p.Id = v.Id
		p.Describe = v.Describe
		p.CreateTime = GetSecond2DateTimes(v.CreateTime)

		if v.UpdateTime > 0 {
			p.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
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
		peoples = append(peoples, p)
	}
	respResult.Users = peoples
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

// 返回的节点信息
type Node struct {
	Id            int    `json:"id"`
	Seo           string `json:"seo"`
	Name          string `json:"name"`
	Describe      string `json:"describe"`
	ImagePath     string `json:"image_path"`
	CreateTime    string `json:"create_time"`
	CreateTimeInt int64  `json:"create_time_int"`
	UpdateTime    string `json:"update_time"`
	UpdateTimeInt int64  `json:"update_time_int"`
	UserId        int    `json:"user_id"`
	UserName      string `json:"user_name"`
	SortNum       int    `json:"sort_num"`
	Level         int    `json:"level"`
	Status        int    `json:"status"`
	ParentNodeId  int    `json:"parent_node_id"`
	Son           []Node `json:"son,omitempty"`
}

type NodesInfoRequest struct {
	UserId   int      `json:"user_id"`
	UserName string   `json:"user_name"`
	Sort     []string `json:"sort"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}

// 列出全部节点
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
		flog.Log.Errorf("ListNode err:%s", "")
		resp.Error = Error(ParasError, "where is empty")
		return
	}

	session := model.FafaRdb.Client.NewSession()
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
		flog.Log.Errorf("ListNode err:%s", err.Error())
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

// 列出一个节点以及他的儿子们
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

	if req.UserId == 0 && req.UserName == "" {
		flog.Log.Errorf("Node err:%s", "")
		resp.Error = Error(ParasError, "where is empty")
		return
	}

	session := model.FafaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.ContentNode)).Where("1=1").And("status=?", 0)

	if req.UserId != 0 {
		session.And("user_id=?", req.UserId)
	}

	if req.UserName != "" {
		session.And("user_name=?", req.UserName)
	}

	isOne := false
	if req.Id != 0 {
		isOne = true
		session.And("id=?", req.Id)
	}

	if req.Seo != "" {
		isOne = true
		session.And("seo=?", req.Seo)
	}

	if !isOne {
		flog.Log.Errorf("Node err:%s", "id or seo empty")
		resp.Error = Error(ParasError, "id or seo empty")
		return
	}

	v := new(model.ContentNode)
	exist, err := session.Get(v)
	if err != nil {
		flog.Log.Errorf("Node err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("Node err:%s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, err.Error())
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

	// 是顶层且需要列出儿子
	if f.Level == 0 && req.ListSon {
		ns := make([]model.ContentNode, 0)
		err = model.FafaRdb.Client.Where("parent_node_id=?", f.Id).Find(ns)
		if err != nil {
			flog.Log.Errorf("Node err:%s", err.Error())
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
			f.Son = append(f.Son, ff)
		}
	}
	resp.Flag = true
	resp.Data = f
}

type UserInfoRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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
	exist, err := model.FafaRdb.Client.Where("status=?", 1).Get(user)
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
	p.CreateTime = GetSecond2DateTimes(v.CreateTime)

	if v.UpdateTime > 0 {
		p.UpdateTime = GetSecond2DateTimes(v.UpdateTime)
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
	resp.Flag = true
	resp.Data = p
}

type UserCountRequest struct {
	UserId   int    `json:"id"`
	UserName string `json:"user_name"`
}

type UserCountX struct {
	Count           int    `json:"count"`
	Days            string `json:"days"`
	CreateTimeBegin int64  `json:"create_time_begin"`
	CreateTimeEnd   int64  `json:"create_time_end"`
}
type UserCountResponse struct {
	Info     []UserCountX `json:"info"`
	UserId   int          `json:"user_id"`
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

	sql := "SELECT DATE_FORMAT(from_unixtime(create_time),'%Y%m%d') days,count(id) count FROM `fafacms_content` WHERE user_id=? and version>0 and status=0 group by days;"
	result, err := model.FafaRdb.Client.QueryString(sql, req.UserId)
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
		begin, _ := time.Parse("20060102", t.Days)
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
	NodeId          int      `json:"node_id"`
	NodeSeo         string   `json:"node_seo"`
	UserId          int      `json:"user_id"`
	UserName        string   `json:"user_name"`
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	Sort            []string `json:"sort"`
	PageHelp
}

type ContentsX struct {
	Id             int    `json:"id"`
	Seo            string `json:"seo"`
	Title          string `json:"title"`
	UserId         int    `json:"user_id"`
	UserName       string `json:"user_name"`
	NodeId         int    `json:"node_id"`
	NodeSeo        string `json:"node_seo"`
	Top            int    `json:"top"`
	CreateTime     string `json:"create_time"`
	PublishTime    string `json:"publish_time,omitempty"`
	CreateTimeInt  int64  `json:"create_time_int"`
	PublishTimeInt int64  `json:"publish_time_int"`
	ImagePath      string `json:"image_path"`
	Views          int    `json:"views"`
	IsLock         bool   `json:"is_lock"`
	Describe       string `json:"describe"`
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
	session := model.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.Content)).Where("1=1")

	if req.UserId != 0 {
		session.And("user_id=?", req.UserId)
	}

	if req.UserName != "" {
		session.And("user_name=?", req.UserName)
	}

	session.And("status=?", 0).And("version>?", 0)

	if req.NodeId != 0 {
		session.And("node_id=?", req.NodeId)
	}

	if req.NodeSeo != "" {
		session.And("node_seo=?", req.NodeSeo)
	}

	if req.CreateTimeBegin > 0 {
		session.And("create_time>=?", req.CreateTimeBegin)
	}

	if req.CreateTimeEnd > 0 {
		session.And("create_time<?", req.CreateTimeEnd)
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
	} else {
		// sql build
		p.build(session, req.Sort, model.ContentSortName)
		// do query
		err = session.Omit("describe", "pre_describe").Find(&cs)
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
		temp.NodeSeo = c.NodeSeo
		temp.UserName = c.UserName
		temp.Id = c.Id
		temp.Top = c.Top
		temp.Title = c.Title
		temp.NodeId = c.NodeId
		temp.Views = c.Views
		temp.CreateTime = GetSecond2DateTimes(c.CreateTime)
		temp.PublishTime = GetSecond2DateTimes(c.PublishTime)
		temp.ImagePath = c.ImagePath
		temp.CreateTimeInt = c.CreateTime
		temp.PublishTimeInt = c.PublishTime
		if c.Password != "" {
			temp.IsLock = true
		}
		bcs = append(bcs, temp)
	}

	respResult.Contents = bcs
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

type ContentRequest struct {
	Id       int    `json:"id"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Seo      string `json:"seo"`
	Password string `json:"password"`
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
		flog.Log.Errorf("TakeContent err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	content := new(model.Content)
	content.Id = req.Id
	content.UserId = req.UserId
	content.UserName = req.UserName
	content.Seo = req.Seo
	exist, err := content.GetByRaw()
	if err != nil {
		flog.Log.Errorf("TakeContent err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("TakeContent err: %s", "content not found")
		resp.Error = Error(ContentNotFound, "")
		return
	}

	if content.Status == 0 {

	} else if content.Status == 2 {
		flog.Log.Errorf("TakeContent err: %s", "content ban")
		resp.Error = Error(ContentBanPermit, "")
		return
	} else {
		flog.Log.Errorf("TakeContent err: %s", "content not found 1")
		resp.Error = Error(ContentNotFound, "")
		return
	}

	if content.Version == 0 {
		flog.Log.Errorf("TakeContent err: %s", "content not found 2")
		resp.Error = Error(ContentNotFound, "")
		return
	}

	if content.Password != "" && content.Password != req.Password {
		flog.Log.Errorf("TakeContent err: %s", "content password")
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
	temp.CreateTime = GetSecond2DateTimes(cx.CreateTime)
	temp.PublishTime = GetSecond2DateTimes(cx.PublishTime)
	temp.CreateTimeInt = cx.CreateTime
	temp.PublishTimeInt = cx.UpdateTime
	temp.ImagePath = cx.ImagePath
	if cx.Password != "" {
		temp.IsLock = true
	}

	temp.Describe = cx.Describe

	cx.UpdateView()

	resp.Flag = true
	resp.Data = temp
}
