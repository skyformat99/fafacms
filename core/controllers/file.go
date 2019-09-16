package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hunterhug/fafacms/core/config"
	. "github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/util/oss"
	"github.com/hunterhug/go_image"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"
	"time"
	myutil "github.com/hunterhug/fafacms/core/util"
)

var scaleType = []string{"jpg", "jpeg", "png"}
var FileAllow = map[string][]string{
	"image": {
		"jpg", "jpeg", "png", "gif"},
	"flash": {
		"swf", "flv"},
	"media": {
		"swf", "flv", "mp3", "wav", "wma", "wmv", "mid", "avi", "mpg", "asf", "rm", "rmvb"},
	"file": {
		"doc", "docx", "xls", "xlsx", "ppt", "htm", "html", "txt", "zip", "rar", "gz", "bz2", "pdf"},
	"other": {
		"jpg", "jpeg", "png", "bmp", "gif", "swf", "flv", "mp3",
		"wav", "wma", "wmv", "mid", "avi", "mpg", "asf", "rm", "rmvb",
		"doc", "docx", "xls", "xlsx", "ppt", "htm", "html", "txt", "zip", "rar", "gz", "bz2"}}

var FileBytes = 1 << 25 // (1<<25)/1000.0/1000.0 33.54 不能超出33M

type UploadResponse struct {
	FileName       string `json:"file_name"`
	ReallyFileName string `json:"really_file_name"`
	Size           int64  `json:"size"`
	Url            string `json:"url"`
	UrlX           string `json:"url_x"`
	IsPicture      bool   `json:"is_picture"`
	Addon          string `json:"addon"`
	Oss            bool   `json:"oss"`
}

/*
	file: 文件form名称
	type: 上传类型，分别为image、flash、media、file、other
	describe: 备注
*/
func UploadFile(c *gin.Context) {
	resp := new(Resp)
	data := UploadResponse{}
	defer func() {
		JSONL(c, 200, nil, resp)
	}()

	uu, err := GetUserSession(c)
	if err != nil {
		Log.Errorf("upload err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		return
	}

	uName := uu.Name

	fileType := c.DefaultPostForm("type", "other")
	if fileType == "" {
		fileType = "other"
	}

	tag := c.DefaultPostForm("tag", "other")
	if tag == "" {
		tag = "other"
	}

	describe := c.DefaultPostForm("describe", "")
	h, err := c.FormFile("file")
	if err != nil {
		Log.Errorf("upload err:%s", err.Error())
		resp.Error = Error(UploadFileError, err.Error())
		return
	}

	fileAllowArray, ok := FileAllow[fileType]
	if !ok {
		Log.Errorf("upload err: type not permit")
		resp.Error = Error(UploadFileTypeNotPermit, "")
		return
	}

	fileSuffix := myutil.GetFileSuffix(h.Filename)
	if !myutil.InArray(fileAllowArray, fileSuffix) {
		Log.Errorf("upload err: file suffix: %s not permit", fileSuffix)
		resp.Error = Error(UploadFileTypeNotPermit, fmt.Sprintf("file suffix: %s not permit", fileSuffix))
		return
	}

	if h.Size > int64(FileBytes) {
		Log.Errorf("upload err: file size too big: %d", h.Size)
		resp.Error = Error(UploadFileTooMaxLimit, fmt.Sprintf(" file size too big: %d", h.Size))
		return
	}

	// 打开文件流
	f, err := h.Open()
	if err != nil {
		Log.Errorf("upload err:%s", err.Error())
		resp.Error = Error(UploadFileError, err.Error())
		return
	}

	defer f.Close()

	// 读取二进制
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		Log.Errorf("upload err:%s", err.Error())
		resp.Error = Error(UploadFileError, err.Error())
		return
	}

	// 二进制空那么报错
	fileSize := len(raw)
	if fileSize == 0 {
		Log.Errorf("upload err:%s", "file empty")
		resp.Error = Error(UploadFileError, "file empty")
		return
	}

	// HashCode
	fileHashCode, err := myutil.Sha256(raw)
	if err != nil {
		Log.Errorf("upload err:%s", err.Error())
		resp.Error = Error(UploadFileError, err.Error())
		return
	}

	// MD5需要再加上用户唯一标志，方便不同用户可以上传一样的文件
	fileHashCode = uName + "_" + fileHashCode
	fileName := fileHashCode + "." + fileSuffix

	// 判断数据库文件是否存在
	p := new(model.File)
	p.HashCode = fileHashCode
	exist, err := p.Get()
	if err != nil {
		resp.Error = Error(DBError, err.Error())
		return
	}

	// 文件不存在
	helpPath := fmt.Sprintf("storage/%s/%s", uName, fileType)
	if !exist {
		fileDir := filepath.Join(config.FafaConfig.DefaultConfig.StoragePath, uName, fileType)
		fileAbName := filepath.Join(fileDir, fileName)

		// 本地存储模式
		if config.FafaConfig.DefaultConfig.StorageOss != true {
			// 磁盘模式
			err := myutil.MakeDir(fileDir)
			if err != nil {
				Log.Errorf("upload err:%s", err.Error())
				resp.Error = Error(UploadFileError, err.Error())
				return
			}

			err = myutil.SaveToFile(fileAbName, raw)
			if err != nil {
				Log.Errorf("upload err:%s", err.Error())
				resp.Error = Error(UploadFileError, err.Error())
				return
			}

			p.Url = fmt.Sprintf("/%s/%s", helpPath, fileName)
		} else {
			// 阿里OSS模式
			p.StoreType = 1
			p.Url = fmt.Sprintf("%s.%s/%s/%s", config.FafaConfig.OssConfig.BucketName, config.FafaConfig.OssConfig.Endpoint, helpPath, fileName)
			err = oss.SaveFile(config.FafaConfig.OssConfig, p.Url, raw)
			if err != nil {
				Log.Errorf("upload err:%s", err.Error())
				resp.Error = Error(UploadFileError, err.Error())
				return
			}
		}

		p.UrlHashCode, _ = myutil.Sha256([]byte(p.Url))

		// 如果是图片进行裁剪
		if myutil.InArray(scaleType, fileSuffix) {
			p.IsPicture = 1

			// 本地存储模式，裁剪图静态路径为  /storage_x
			if config.FafaConfig.DefaultConfig.StorageOss != true {
				fileScaleDir := filepath.Join(config.FafaConfig.DefaultConfig.StoragePath+"_x", uName, fileType)
				fileScaleAbName := filepath.Join(fileScaleDir, fileName)

				// 裁剪
				err = myutil.MakeDir(fileScaleDir)
				if err != nil {
					Log.Errorf("upload err:%s", err.Error())
					resp.Error = Error(UploadFileError, err.Error())
					return
				}
				err := go_image.ScaleF2F(fileAbName, fileScaleAbName, 200)
				if err != nil {
					Log.Errorf("upload err:%s", err.Error())
					resp.Error = Error(UploadFileError, err.Error())
					return
				}
			} else {
				// 阿里OSS模式
				outRaw, err := go_image.ScaleB2B(raw, 200)
				if err != nil {
					Log.Errorf("upload err:%s", err.Error())
					resp.Error = Error(UploadFileError, err.Error())
					return
				}

				err = oss.SaveFile(config.FafaConfig.OssConfig, strings.Replace(helpPath, "storage/", "storage_x/", -1)+"/"+fileName, outRaw)
				if err != nil {
					Log.Errorf("upload err:%s", err.Error())
					resp.Error = Error(UploadFileError, err.Error())
					return
				}
			}
		}

		p.Type = fileType
		p.FileName = fileName
		p.ReallyFileName = h.Filename
		p.CreateTime = time.Now().Unix()
		p.Describe = describe
		p.UserId = uu.Id
		p.UserName = uName
		p.Tag = tag
		p.Size = int64(fileSize)
		_, err = model.FafaRdb.InsertOne(p)
		if err != nil {
			Log.Errorf("upload err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	} else {
		// 文件存在
		data.Addon = "file the same in server"
		if p.Status != 0 {
			// 如果被隐藏了额，应该改回来
			p.Status = 0
			p.UpdateStatus()
		}
	}

	// 返回基本信息
	data.FileName = p.FileName
	data.ReallyFileName = p.ReallyFileName
	data.IsPicture = p.IsPicture == 1
	data.Size = p.Size
	data.Url = p.Url
	data.Oss = p.StoreType == 1
	if data.IsPicture {
		data.UrlX = strings.Replace(p.Url, "/storage", "/storage_x", -1)
	}

	resp.Data = data
	resp.Flag = true
	return
}

type ListFileAdminRequest struct {
	CreateTimeBegin int64    `json:"create_time_begin"`
	CreateTimeEnd   int64    `json:"create_time_end"`
	UpdateTimeBegin int64    `json:"update_time_begin"`
	UpdateTimeEnd   int64    `json:"update_time_end"`
	SizeBegin       int64    `json:"size_begin"`
	SizeEnd         int64    `json:"size_end"`
	Sort            []string `json:"sort"`
	HashCode        string   `json:"hash_code"`
	Url             string   `json:"url"`
	StoreType       int      `json:"store_type" validate:"oneof=-1 0 1"`
	Status          int      `json:"status" validate:"oneof=-1 0 1"`
	Type            string   `json:"type"`
	Tag             string   `json:"tag"`
	UserId          int64    `json:"user_id"`
	Id              int64    `json:"id"`
	IsPicture       int      `json:"is_picture" validate:"oneof=-1 0 1"`
	PageHelp
}

type ListFileAdminResponse struct {
	Files []model.File `json:"files"`
	PageHelp
}

func ListFileAdminHelper(c *gin.Context, userId int64) {
	resp := new(Resp)

	respResult := new(ListFileAdminResponse)
	req := new(ListFileAdminRequest)
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
		Log.Errorf("ListFileAdmin err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// new query list session
	session := model.FafaRdb.Client.NewSession()
	defer session.Close()

	// group list where prepare
	session.Table(new(model.File)).Where("1=1")

	// query prepare
	if req.Id != 0 {
		session.And("id=?", req.Id)
	}
	if req.HashCode != "" {
		session.And("hash_code=?", req.HashCode)
	}

	if req.Status != -1 {
		// 只要不暴露这个参数的话，那些隐藏的文件不会被用户看到
		session.And("status=?", req.Status)
	}

	if req.Url != "" {
		urlHashCode, _ := myutil.Sha256([]byte(req.Url))
		session.And("url_hash_code=?", urlHashCode)
	}

	if req.IsPicture != -1 {
		session.And("is_picture=?", req.IsPicture)
	}

	if req.Type != "" {
		session.And("type=?", req.Type)
	}

	if req.StoreType != -1 {
		session.And("store_type=?", req.StoreType)
	}

	if req.Tag != "" {
		session.And("tag=?", req.Tag)
	}

	if userId != 0 {
		session.And("user_id=?", userId)
	} else {
		if req.UserId != 0 {
			session.And("user_id=?", req.UserId)
		}
	}

	if req.CreateTimeBegin > 0 {
		session.And("create_time>=?", req.CreateTimeBegin)
	}

	if req.CreateTimeEnd > 0 {
		session.And("create_time<?", req.CreateTimeEnd)
	}

	if req.UpdateTimeBegin > 0 {
		session.And("update_time>=?", req.UpdateTimeBegin)
	}

	if req.UpdateTimeEnd > 0 {
		session.And("update_time<?", req.UpdateTimeEnd)
	}

	if req.SizeBegin > 0 {
		session.And("size>=?", req.SizeBegin)
	}

	if req.SizeEnd > 0 {
		session.And("size<?", req.SizeEnd)
	}

	// count num
	countSession := session.Clone()
	defer countSession.Close()
	total, err := countSession.Count()
	if err != nil {
		Log.Errorf("ListFileAdmin err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	// if count>0 start list
	files := make([]model.File, 0)
	p := &req.PageHelp
	if total == 0 {
		if p.Limit == 0 {
			p.Limit = 20
		}
	} else {
		// sql build
		p.build(session, req.Sort, model.FileSortName)
		// do query
		err = session.Find(&files)
		if err != nil {
			Log.Errorf("ListFileAdmin err:%s", err.Error())
			resp.Error = Error(DBError, err.Error())
			return
		}
	}

	// result
	respResult.Files = files
	p.Pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	respResult.PageHelp = *p
	resp.Data = respResult
	resp.Flag = true
}

// 列出所有用户的文件信息，管理员权限
func ListFileAdmin(c *gin.Context) {
	ListFileAdminHelper(c, 0)
}

// 列出自己的文件
func ListFile(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		Log.Errorf("ListFile err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	ListFileAdminHelper(c, uid)
}

type UpdateFileRequest struct {
	Id       int64  `json:"id" validate:"required"`
	Tag      string `json:"tag"`
	Hide     bool   `json:"hide"`
	Describe string `json:"describe"`
}

func UpdateFileAdminHelper(c *gin.Context, userId int64) {
	resp := new(Resp)
	req := new(UpdateFileRequest)
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
		Log.Errorf("UpdateFileAdmin err: %s", err.Error())
		resp.Error = Error(ParasError, err.Error())
		return
	}

	// 可以修改文件Tag和描述，方便分组
	f := new(model.File)
	f.Id = req.Id
	f.Tag = req.Tag
	f.Describe = req.Describe
	f.UserId = userId

	// 更改文件，可以将文件设置为隐藏，文件一旦上传，不能删除
	ok, err := f.Update(req.Hide)
	if err != nil {
		Log.Errorf("UpdateFileAdmin err:%s", err.Error())
		resp.Error = Error(DBError, err.Error())
		return
	}

	resp.Data = ok
	resp.Flag = true
}

// 可以更改其他人的文件信息，管理员权限
func UpdateFileAdmin(c *gin.Context) {
	UpdateFileAdminHelper(c, 0)
}

// 更新自己的文件信息
func UpdateFile(c *gin.Context) {
	resp := new(Resp)
	uu, err := GetUserSession(c)
	if err != nil {
		Log.Errorf("UpdateFile err: %s", err.Error())
		resp.Error = Error(GetUserSessionError, err.Error())
		JSONL(c, 200, nil, resp)
		return
	}

	uid := uu.Id
	UpdateFileAdminHelper(c, uid)
}
