package controllers

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"strings"
)

// error code
const (
	GetUserSessionError               = 100000
	SetUserSessionError               = 100001
	UserNoLogin                       = 100002
	UserNotFound                      = 100003
	UserNotActivate                   = 100004
	UserIsInBlack                     = 100005
	UserAuthPermit                    = 100006
	ParasError                        = 100010
	ParseJsonError                    = 100011
	LoginWrong                        = 100020
	CloseRegisterError                = 100021
	UserNameAlreadyBeUsed             = 100022
	EmailAlreadyBeUsed                = 100023
	ActivateCodeWrong                 = 100024
	ActivateCodeExpired               = 100025
	ActivateCodeNotExpired            = 100026
	EmailNotFound                     = 100027
	ResetCodeExpiredTimeNotReach      = 100028
	RestCodeWrong                     = 100029
	FileCanNotBeFound                 = 100030
	GroupNameAlreadyBeUsed            = 100040
	GroupNotFound                     = 100041
	GroupHasResourceHookIn            = 100042
	GroupHasUserHookIn                = 100043
	ResourceCountNumNotRight          = 100050
	UploadFileError                   = 100100
	UploadFileTypeNotPermit           = 100101
	UploadFileTooMaxLimit             = 100102
	ContentNodeSeoAlreadyBeUsed       = 101000
	ContentNodeNotFound               = 101001
	ContentParentNodeNotFound         = 101002
	ContentNodeSortConflict           = 101003
	ContentNodeHasChildren            = 101004
	ContentNodeHasContentCanNotDelete = 101005
	ContentNotFound                   = 110000
	ContentPasswordWrong              = 110001
	ContentBanPermit                  = 110002
	ContentSeoAlreadyBeUsed           = 110003
	ContentInRubbish                  = 110004
	ContentsAreInDifferentNode        = 110005
	ContentHistoryNotFound            = 110006
	ContentCanNotDelete               = 110007
	AddUserCacheError                 = 120000
	DeleteUserCacheError              = 120001
	RefreshUserCacheError             = 120002
	DeleteUserAllSessionError         = 120003
	DeleteUserSessionError            = 120004
	RefreshUserSessionError           = 120005

	DBError    = 200000
	DbNotFound = 200001
	DbRepeat   = 200002
	DbHookIn   = 200003

	EmailSendError = 300000
	SystemProblem  = 300001

	LazyError = 99997
	I500      = 99998
	Unknown   = 99999
)

// error code message map
var ErrorMap = map[int]string{
	AddUserCacheError:                 "add user cache err",
	DeleteUserCacheError:              "delete user cache err",
	RefreshUserCacheError:             "refresh user cache err",
	DeleteUserAllSessionError:         "delete user all session err",
	DeleteUserSessionError:            "delete user session err",
	GetUserSessionError:               "get user session err",
	SetUserSessionError:               "set user session err",
	RefreshUserSessionError:           "refresh user session err",
	UserNoLogin:                       "user no login",
	UserNotFound:                      "user not found",
	UserIsInBlack:                     "user is in black",
	UserNotActivate:                   "user not active",
	UserAuthPermit:                    "user auth permit",
	ParasError:                        "paras input not right",
	DBError:                           "db operation err",
	LoginWrong:                        "username or password wrong",
	CloseRegisterError:                "register close",
	ParseJsonError:                    "json parse err",
	UserNameAlreadyBeUsed:             "user name already be used",
	EmailAlreadyBeUsed:                "email already be used",
	ActivateCodeWrong:                 "activate code wrong",
	ActivateCodeExpired:               "activate code expired",
	ActivateCodeNotExpired:            "activate code not expired",
	FileCanNotBeFound:                 "file can not be found",
	EmailSendError:                    "email send error",
	EmailNotFound:                     "email not found",
	ResetCodeExpiredTimeNotReach:      "reset code expired time not reach",
	RestCodeWrong:                     "reset code wrong",
	GroupNameAlreadyBeUsed:            "group name already be used",
	GroupNotFound:                     "group not found",
	GroupHasResourceHookIn:            "group has resource hook in",
	GroupHasUserHookIn:                "group has user hook in",
	ResourceCountNumNotRight:          "resource count not right",
	UploadFileError:                   "upload file err",
	UploadFileTypeNotPermit:           "upload file type not permit",
	UploadFileTooMaxLimit:             "upload file too max limit",
	ContentNodeSeoAlreadyBeUsed:       "content node seo already be used",
	ContentNodeNotFound:               "content node not found",
	ContentParentNodeNotFound:         "parent content node not found",
	ContentNodeSortConflict:           "content node sort conflict",
	ContentNodeHasChildren:            "content node has children",
	ContentNodeHasContentCanNotDelete: "content node has content can not delete",
	ContentNotFound:                   "content not found",
	ContentBanPermit:                  "content ban permit",
	ContentPasswordWrong:              "content password wrong",
	ContentSeoAlreadyBeUsed:           "content seo already be used",
	ContentInRubbish:                  "content in rubbish",
	ContentsAreInDifferentNode:        "contents are in different node",
	ContentHistoryNotFound:            "content history can not found",
	ContentCanNotDelete:               "content can not delete for content not in rubbish",
	SystemProblem:                     "system problem",
	DbNotFound:                        "db not found",
	DbRepeat:                          "db repeat data",
	DbHookIn:                          "db hook in",
	I500:                              "500 error",
	LazyError:                         "db not found or err",
}

// common response
type Resp struct {
	Flag  bool        `json:"flag"`
	Cid   string      `json:"cid,omitempty"`
	Error *ErrorResp  `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// inner error response
type ErrorResp struct {
	ErrorID  int    `json:"id"`
	ErrorMsg string `json:"msg"`
}

func (e ErrorResp) Error() string {
	return fmt.Sprintf("%d|%s", e.ErrorID, e.ErrorMsg)
}

func Error(code int, detail string) *ErrorResp {
	_, ok := ErrorMap[code]
	if !ok {
		code = Unknown
	}

	str := fmt.Sprintf("%s:%s", ErrorMap[code], detail)

	if detail == "" {
		str = fmt.Sprintf("%s", ErrorMap[code])
	}

	return &ErrorResp{
		ErrorID:  code,
		ErrorMsg: str,
	}
}

// list api page helper
type PageHelp struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Pages int `json:"total_pages"` // set by yourself outside
}

// page build helper
func (page *PageHelp) build(s *xorm.Session, sort []string, base []string) {
	Build(s, sort, base)

	if page.Page == 0 {
		page.Page = 1
	}

	if page.Limit <= 0 {
		page.Limit = 20
	}

	if page.Limit > 100 {
		page.Limit = 100
	}
	s.Limit(page.Limit, (page.Page-1)*page.Limit)
}

func Build(s *xorm.Session, sort []string, base []string) {
	nowSort := make([]string, 0, len(sort))
	for _, v := range sort {
		nowSort = append(nowSort, v)
	}

	dict := make(map[string]struct{}, 0)

	for _, v := range base {
		a := v[1:]
		dict[a] = struct{}{}

		// if default use base sort field
		useBase := true
		for _, vv := range sort {
			b := vv[1:]
			if a == b {
				useBase = false
			}
		}

		if useBase {
			nowSort = append(nowSort, v)
		}
	}

	for _, v := range nowSort {
		a := v[1:]
		if _, ok := dict[a]; ok {
			if strings.HasPrefix(v, "+") {
				s.Asc(a)
			} else if strings.HasPrefix(v, "-") {
				s.Desc(a)

			}
		}

	}
}
