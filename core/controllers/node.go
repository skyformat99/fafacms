package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
)

type CreateNodeRequest struct {
	Seo          string `json:"seo" validate:"omitempty,alphanumunicode"`
	Name         string `json:"name" validate:"required"`
	Describe     string `json:"describe"`
	ImagePath    string `json:"image_path"`
	ParentNodeId int    `json:"parent_node_id"`
}

func CreateNode(c *gin.Context) {
	resp := new(Resp)
	req := new(CreateNodeRequest)
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
		flog.Log.Errorf("CreateNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("CreateNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	n := new(model.ContentNode)
	n.UserId = uu.Id

	// 如果SEO非空，检查是否已经存在
	if req.Seo != "" {
		n.Seo = req.Seo
		exist, err := n.CheckSeoValid()
		if err != nil {
			flog.Log.Errorf("CreateNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		if exist {
			// 存在报错
			flog.Log.Errorf("CreateNode err: %s", "node seo already be use")
			resp.Error = Error(ContentNodeSeoAlreadyBeUsed, "")
			return
		}
	}

	// 如果指定了父亲节点
	if req.ParentNodeId != 0 {
		n.ParentNodeId = req.ParentNodeId
		exist, err := n.CheckParentValid()
		if err != nil {
			flog.Log.Errorf("CreateNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		if !exist {
			// 父亲节点不存在，报错
			flog.Log.Errorf("CreateNode err: %s", "parent content node not found")
			resp.Error = Error(ContentParentNodeNotFound, "")
			return
		}

		n.Level = 1
	}

	// if image not empty
	if req.ImagePath != "" {
		n.ImagePath = req.ImagePath
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			flog.Log.Errorf("CreateNode err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("CreateNode err: image not exist")
			resp.Error = Error(FileCanNotBeFound, "image url not exist")
			return
		}
	}
	n.Name = req.Name
	n.Describe = req.Describe
	n.ParentNodeId = req.ParentNodeId
	n.UserName = uu.Name
	n.SortNum, _ = n.CountNodeNum()
	err = n.InsertOne()
	if err != nil {
		flog.Log.Errorf("CreateNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
	resp.Data = n
}

type UpdateInfoOfNodeRequest struct {
	Id       int    `json:"id" validate:"required"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
}

type UpdateImageOfNodeRequest struct {
	Id        int    `json:"id" validate:"required"`
	ImagePath string `json:"image_path" validate:"required"`
}

type UpdateStatusOfNodeRequest struct {
	Id     int `json:"id" validate:"required"`
	Status int `json:"status" validate:"oneof=0 1"`
}

type UpdateSeoOfNodeRequest struct {
	Id  int    `json:"id" validate:"required"`
	Seo string `json:"seo" validate:"required,alphanumunicode"`
}

type UpdateParentOfNodeRequest struct {
	Id           int  `json:"id" validate:"required"`
	ToBeRoot     bool `json:"to_be_root"` // 升级为最上层节点
	ParentNodeId int  `json:"parent_node_id"`
}

func UpdateSeoOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateSeoOfNodeRequest)
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
		flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateSeoOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id

	seoChange := false

	// 和之前的SEO不一样
	if req.Seo != n.Seo {
		after.Seo = req.Seo
		seoChange = true
		// 检查是否存在SEO
		exist, err := after.CheckSeoValid()
		if err != nil {
			flog.Log.Errorf("UpdateSeoOfNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		if exist {
			// SEO存在了，报错
			flog.Log.Errorf("UpdateSeoOfNode err: %s", "seo been used")
			resp.Error = Error(ContentNodeSeoAlreadyBeUsed, "")
			return
		}
	}

	if seoChange {
		// 更新
		err = after.UpdateSeo()
		if err != nil {
			flog.Log.Errorf("UpdateSeoOfNode err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}
	resp.Flag = true
}

func UpdateInfoOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateInfoOfNodeRequest)
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
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, "")
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateInfoOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id

	// 以下只要存在不一致性才替换
	if req.Name != "" {
		if req.Name != n.Name {
			after.Name = req.Name
		}
	}

	after.Describe = req.Describe

	// 更新
	err = after.UpdateInfo()
	if err != nil {
		flog.Log.Errorf("UpdateNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

func UpdateImageOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateImageOfNodeRequest)
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
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, "")
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateInfoOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateInfoOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id

	if req.ImagePath != n.ImagePath {
		after.ImagePath = req.ImagePath
		p := new(model.File)
		p.Url = req.ImagePath
		ok, err := p.Exist()
		if err != nil {
			flog.Log.Errorf("UpdateInfoOfNode err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		if !ok {
			flog.Log.Errorf("UpdateInfoOfNode err: image not exist")
			resp.Error = Error(FileCanNotBeFound, "")
			return
		}

		// 更新
		err = after.UpdateImage()
		if err != nil {
			flog.Log.Errorf("UpdateNode err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	resp.Flag = true
}

func UpdateStatusOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateStatusOfNodeRequest)
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
		flog.Log.Errorf("UpdateStatusOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateStatusOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateStatusOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateStatusOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id
	after.Status = req.Status

	// 更新
	err = after.UpdateStatus()
	if err != nil {
		flog.Log.Errorf("UpdateStatusOfNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

func UpdateParentOfNode(c *gin.Context) {
	resp := new(Resp)
	req := new(UpdateParentOfNodeRequest)
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
		flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	if req.ParentNodeId == req.Id {
		flog.Log.Errorf("UpdateParentOfNode err: %s", "self can not be parent")
		resp.Error = Error(ParasError, "self can not be parent")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("UpdateParentOfNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	// 有儿子的节点不能成为儿子，毕竟只有两层
	childNum, err := n.CheckChildrenNum()
	if err != nil {
		flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if childNum > 0 {
		flog.Log.Errorf("UpdateParentOfNode err: %s", "has child")
		resp.Error = Error(ContentNodeHasChildren, "has child")
		return
	}

	beforeParentNode := n.ParentNodeId

	after := new(model.ContentNode)
	after.UserId = n.UserId
	after.Id = n.Id
	after.SortNum = n.SortNum

	// 成为根节点
	if req.ToBeRoot {
		// 已经是了
		if n.ParentNodeId == 0 {
			resp.Flag = true
			return
		}
		// 没有指定父亲节点，归零
		after.Level = 0
		after.ParentNodeId = 0
	} else {
		if n.ParentNodeId == req.ParentNodeId {
			resp.Flag = true
			return
		}

		after.ParentNodeId = req.ParentNodeId

		// 检查该父亲节点是否存在
		exist, err := after.CheckParentValid()
		if err != nil {
			flog.Log.Errorf("UpdateParentOfNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		if !exist {
			// 不存在父亲节点，报错
			flog.Log.Errorf("UpdateParentOfNode err: %s", "parent content node not found")
			resp.Error = Error(ContentParentNodeNotFound, "")
			return
		}
		after.Level = 1
	}

	// 更新
	err = after.UpdateParent(beforeParentNode)
	if err != nil {
		flog.Log.Errorf("UpdateParentOfNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

type DeleteNodeRequest struct {
	Id int `json:"id" validate:"required"`
}

// 删除节点
// 删除的时候，后面的节点排序值要依次顶上
func DeleteNode(c *gin.Context) {
	resp := new(Resp)
	req := new(DeleteNodeRequest)
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
		flog.Log.Errorf("DeleteNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("DeleteNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}
	n := new(model.ContentNode)
	n.Id = req.Id
	n.UserId = uu.Id

	// 获取节点，节点会携带所有内容
	exist, err := n.Get()
	if err != nil {
		flog.Log.Errorf("DeleteNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	if !exist {
		// 不存在节点，报错
		flog.Log.Errorf("DeleteNode err: %s", "content node not found")
		resp.Error = Error(ContentNodeNotFound, "")
		return
	}

	// 删除节点时节点下不能有节点
	childNum, err := n.CheckChildrenNum()
	if err != nil {
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if childNum >= 1 {
		// 不能删除
		flog.Log.Errorf("DeleteNode err:%s", "has node child")
		resp.Error = Error(ContentNodeHasChildren, "")
		return
	}

	content := new(model.Content)
	content.UserId = uu.Id
	content.NodeId = n.Id

	// 删除节点时，节点下不能有内容
	normalContentNum, err := content.CountNumUnderNode()
	if err != nil {
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if normalContentNum >= 1 {
		// 有内容，不能删除
		flog.Log.Errorf("DeleteNode err:%s", "has content child")
		resp.Error = Error(ContentNodeHasContentCanNotDelete, "")
		return
	}

	session := model.FafaRdb.Client.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// 同一层的大于删除节点排序的节点-1 ,依次顶上
	_, err = session.Exec("update fafacms_content_node SET sort_num=sort_num-1 where sort_num > ? and user_id = ? and parent_node_id = ?", n.SortNum, n.UserId, n.ParentNodeId)
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// 可以删除了
	_, err = session.Where("id=?", n.Id).Delete(new(model.ContentNode))
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("DeleteNode err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
}

func TakeNode(c *gin.Context) {
	resp := new(Resp)
	req := new(NodeInfoRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("TakeNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	session := model.FafaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.ContentNode)).Where("1=1").And("user_id=?", uu.Id)

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
	f.Status = v.Status

	// 是顶层且需要列出儿子
	if f.Level == 0 && req.ListSon {
		ns := make([]model.ContentNode, 0)
		err = model.FafaRdb.Client.Where("parent_node_id=?", f.Id).Find(&	ns)
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
			ff.Status = vv.Status
			f.Son = append(f.Son, ff)
		}
	}
	resp.Flag = true
	resp.Data = f

}

func ListNode(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("ListNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	ListNodeHelper(c, uid)
}

func ListNodeAdmin(c *gin.Context) {
	ListNodeHelper(c, 0)
}

// 可以查别人的节点
func ListNodeHelper(c *gin.Context, userId int) {
	resp := new(Resp)

	respResult := new(NodesResponse)
	req := new(NodesInfoRequest)
	defer func() {
		JSONL(c, 200, req, resp)
	}()

	if errResp := ParseJSON(c, req); errResp != nil {
		resp.Error = errResp
		return
	}

	if userId != 0 {
		req.UserId = userId
		req.UserName = ""
	}

	if req.UserId == 0 && req.UserName == "" {
		flog.Log.Errorf("ListNode err:%s", "")
		resp.Error = Error(ParasError, "where is empty")
		return
	}

	session := model.FafaRdb.Client.NewSession()
	defer session.Close()

	session.Table(new(model.ContentNode)).Where("1=1")

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
		f.Status = v.Status
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
				s.Status = vv.Status
				f.Son = append(f.Son, s)
			}
		}

		n = append(n, f)
	}

	respResult.Nodes = n
	resp.Flag = true
	resp.Data = respResult
}

// 将X节点放在Y节点的下面
// 所以想把一个X节点拖到父节点的最顶层是做不到的，需要拖两次
// 上面拖两次的问题可以解决了，就是YID为空时，直接把他拖到最前面
type SortNodeRequest struct {
	XID int `json:"xid" validate:"required"`
	YID int `json:"yid"`
}

//  节点
//  拖曳排序超级函数
func SortNode(c *gin.Context) {
	resp := new(Resp)
	req := new(SortNodeRequest)
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
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	if req.XID == req.YID {
		flog.Log.Errorf("SortNode err: %s", "xid=yid not right")
		resp.Error = Error(ParasError, "xid=yid not right")
		return
	}

	uu, err := GetUserSession(c)
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	x := new(model.ContentNode)
	x.Id = req.XID
	x.UserId = uu.Id
	exist, err := x.GetSortOneNode()
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("SortNode err: %s", "x node not found")
		resp.Error = Error(ContentNodeNotFound, "x node not found")
		return
	}

	// x节点要拉到最顶层
	if req.YID == 0 {
		session := model.FafaRdb.Client.NewSession()
		defer session.Close()

		err = session.Begin()
		if err != nil {
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		// 比x小的都往下走，因为x要做大哥
		//  --- a  0		---
		//  --- x  1   ==》	--- a x 1
		//  --- b  2		--- b 2
		_, err = session.Exec("update fafacms_content_node SET sort_num=sort_num+1 where sort_num < ? and user_id = ? and parent_node_id = ?", x.SortNum, uu.Id, x.ParentNodeId)
		if err != nil {
			session.Rollback()
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}

		// x做大哥了
		_, err = session.Exec("update fafacms_content_node SET sort_num=0 where user_id = ? and parent_node_id = ? and id = ?", uu.Id, x.ParentNodeId, x.Id)
		if err != nil {
			session.Rollback()
			flog.Log.Errorf("SortNode err: %s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
		resp.Flag = true
		return
	}

	y := new(model.ContentNode)
	y.Id = req.YID
	y.UserId = uu.Id
	exist, err = y.GetSortOneNode()
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	if !exist {
		flog.Log.Errorf("SortNode err: %s", "y node not found")
		resp.Error = Error(ContentNodeNotFound, "y node not found")
		return
	}

	// x是y的爸爸了，怎么可以和儿子做兄弟
	if y.ParentNodeId == x.Id {
		flog.Log.Errorf("SortNode err: %s", "can not move node to be his child's brother")
		resp.Error = Error(ContentNodeSortConflict, "can not move node to be his child's brother")
		return
	}

	children, err := x.CheckChildrenNum()
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// 当y是儿子节点，但x想当y兄弟，然而x已经有一堆儿子了，当不了兄弟
	if y.Level == 1 && children > 0 {
		flog.Log.Errorf("SortNode err: %s", "x has child can not move to be other's child's brother")
		resp.Error = Error(ContentNodeSortConflict, "x has child can not move to be other's child's brother")
		return
	}

	session := model.FafaRdb.Client.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// 先把x假装删掉，比x大的都-1，依次顶上x的位置
	_, err = session.Exec("update fafacms_content_node SET sort_num=sort_num-1 where sort_num > ? and user_id = ? and parent_node_id = ?", x.SortNum, uu.Id, x.ParentNodeId)
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// 把大于y排序的节点都+1，腾出位置给x
	_, err = session.Exec("update fafacms_content_node SET sort_num=sort_num+1 where sort_num > ? and user_id = ? and parent_node_id = ?", y.SortNum, uu.Id, y.ParentNodeId)
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}


	// 先把x假装删掉，比x大的都-1，依次顶上x的位置，把大于y排序的节点都+1，腾出位置给x
	// 同一级
	// y>x  y=5 x=2
	//   --- a	0		--- a 0			--- a 0
	//   --- b  1 ==>	--- b 1 	==>	--- b 1
	//   --- x	2		--- xc 2		--- xc 2
	//   --- c	3		--- d 3			--- d 3
	//   --- d  4		--- y 4			--- y 4
	//   --- y	5		--- e 5			---			==> x=5
	//   --- e	6		---  			--- e 6

	// y<x  y=2 x=5
	//   --- a	0		--- a 0			--- a 0
	//   --- b  1 ==>	--- b 1 	==>	--- b 1
	//   --- y	2		--- y 2			--- y 2		==> x=3
	//   --- c	3		--- c 3			--- c 4
	//   --- d  4		--- d 4			--- d 5
	//   --- x	5		--- xe 5		---	xe 6
	//   --- e	6		---

	// 不同级
	// y=1
	//   --- a	0		--- a 0				--- a 0
	//   --- b  1 	==>	--- b 1 		==>	--- b 1
	//   	--- c 0			--- c 0				--- c 0
	//   	--- y 1			--- y 1				--- y 1		==> x=2
	//   	--- d 2			--- d 2				--- d 3
	//   --- x	2		--- xe 2			---	xe 2
	//   --- e	3

	// x顶上, 如果同一级的，且y>x
	if x.ParentNodeId == y.ParentNodeId && y.SortNum > x.SortNum {
		_, err = session.Exec("update fafacms_content_node SET sort_num=?,level=?,parent_node_id=? where user_id = ? and id = ?", y.SortNum, y.Level, y.ParentNodeId, uu.Id, x.Id)
	} else {
		// 否则
		_, err = session.Exec("update fafacms_content_node SET sort_num=?,level=?,parent_node_id=? where user_id = ? and id = ?", y.SortNum+1, y.Level, y.ParentNodeId, uu.Id, x.Id)
	}
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		flog.Log.Errorf("SortNode err: %s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}
	resp.Flag = true
	return
}
