package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunterhug/fafacms/core/controllers"
)

type HttpHandle struct {
	Name   string
	Func   gin.HandlerFunc
	Method []string
	Admin  bool
}

var (
	POST = []string{"POST"}
	GET  = []string{"GET"}
	GP   = []string{"POST", "GET"}
)

// Router
var (
	HomeRouter = map[string]HttpHandle{
		// 前端路由
		// 需要考虑更友好的展示，反盗链，反爬虫等
		"/": {"Home", controllers.Home, GP, false},

		"/u":         {"List Peoples", controllers.Peoples, GP, false},         // 列出用户
		"/u/node":    {"List User Nodes One", controllers.NodeInfo, GP, false}, // 查找某一个节点
		"/u/nodes":   {"List User Nodes", controllers.NodesInfo, GP, false},    // 列出某用户下的节点
		"/u/info":    {"List User Info", controllers.UserInfo, GP, false},      // 获取某用户信息
		"/u/count":   {"Count User Content", controllers.UserCount, GP, false}, // 统计某用户文章情况（某用户可留空）
		"/u/content": {"List User Content", controllers.Contents, GP, false},   // 列出某用户下文章（某用户可留空）
		"/content":   {"Get Content", controllers.Content, GP, false},          // 获取文章

		// start at 2019/9
		"/content/comment": {"List Comment of Content", controllers.ListHomeComment, GP, false}, // 列出文章下的评论

		// 前端的用户授权路由，不需要登录即可操作
		// 已经Review 2019/5/12 chen
		"/user/token/get":       {"User Token get", controllers.Login, GP, false},
		"/user/token/refresh":   {"User Token refresh", controllers.Refresh, GP, false},
		"/user/token/delete":    {"User Token delete", controllers.Logout, GP, false},
		"/user/register":        {"User Register", controllers.RegisterUser, GP, false},
		"/user/activate":        {"User Verify Email To Activate", controllers.ActivateUser, GP, false},               // 用户自己激活
		"/user/activate/code":   {"User Resend Email Activate Code", controllers.ResendActivateCodeToUser, GP, false}, // 激活码过期重新获取
		"/user/password/forget": {"User Forget Password Gen Code", controllers.ForgetPasswordOfUser, GP, false},       // 忘记密码，验证码发往邮箱
		"/user/password/change": {"User Change Password", controllers.ChangePasswordOfUser, GP, false},                // 根据邮箱验证码修改密码
	}

	// /v1/user/create
	// need login group auth
	V1Router = map[string]HttpHandle{
		// 用户组操作
		// 已经Review 2019/5/12 chen
		"/group/create":        {"Create Group", controllers.CreateGroup, POST, true},
		"/group/update":        {"Update Group", controllers.UpdateGroup, POST, true},
		"/group/delete":        {"Delete Group", controllers.DeleteGroup, POST, true},
		"/group/take":          {"Take Group", controllers.TakeGroup, GP, true},
		"/group/list":          {"List Group", controllers.ListGroup, GP, true},
		"/group/user/list":     {"Group List User", controllers.ListGroupUser, GP, true},         // 超级管理员列出组下的用户
		"/group/resource/list": {"Group List Resource", controllers.ListGroupResource, GP, true}, // 超级管理员列出组下的资源

		// 用户操作
		// 已经Review 2019/5/12 chen
		"/user/list":         {"User List All", controllers.ListUser, GP, true},              // 超级管理员列出用户列表
		"/user/create":       {"User Create", controllers.CreateUser, GP, true},              // 超级管理员创建用户，默认激活
		"/user/assign":       {"User Assign Group", controllers.AssignGroupToUser, GP, true}, // 超级管理员给用户分配用户组
		"/user/update":       {"User Update Self", controllers.UpdateUser, GP, false},        // 更新自己的信息
		"/user/admin/update": {"User Update Admin", controllers.UpdateUserAdmin, GP, true},   // 管理员修改其他用户信息，可以修改用户密码，以及将用户加入黑名单，禁止使用等
		"/user/info":         {"User Info Self", controllers.TakeUser, GP, false},            // 获取自己的信息

		// 资源操作
		// 已经Review 2019/5/12 chen
		"/resource/list":   {"Resource List All", controllers.ListResource, GP, true},              // 列出资源
		"/resource/assign": {"Resource Assign Group", controllers.AssignResourceToGroup, GP, true}, // 资源分配给组

		// 文件操作
		// 已经Review 2019/5/12 chen
		"/file/upload":       {"File Upload", controllers.UploadFile, POST, false},
		"/file/list":         {"File List Self", controllers.ListFile, POST, false},
		"/file/admin/list":   {"File List All", controllers.ListFileAdmin, POST, true}, // 管理员查看所有文件
		"/file/update":       {"File Update Self", controllers.UpdateFile, POST, false},
		"/file/admin/update": {"File Update All", controllers.UpdateFileAdmin, POST, true}, // 管理员修改文件

		// 比较重要的, 节点和文章都应该支持拖曳，文章首页排序还是按照创建时间，但是后台使用排序字段
		// 需要参考简书
		// 内容节点操作
		// 已经Review 2019/5/13 chen
		"/node/create":        {"Create Node Self", controllers.CreateNode, POST, false},
		"/node/update/seo":    {"Update Node Self Seo", controllers.UpdateSeoOfNode, POST, false},          // 更新节点SEO
		"/node/update/info":   {"Update Node Self Info", controllers.UpdateInfoOfNode, POST, false},        // 更新节点名字和描述
		"/node/update/image":  {"Update Node Self Info Image", controllers.UpdateImageOfNode, POST, false}, // 更新图片地址
		"/node/update/status": {"Update Node Self Status", controllers.UpdateStatusOfNode, POST, false},    // 更新状态，可以设置隐藏
		"/node/update/parent": {"Update Node Self Parent", controllers.UpdateParentOfNode, POST, false},    // 这个接口不如下面这个全功能的接口
		"/node/sort":          {"Sort Node Self", controllers.SortNode, POST, false},                       // 拖曳超级函数
		"/node/delete":        {"Delete Node Self", controllers.DeleteNode, POST, false},

		// 已经Review 2019/5/14 chen
		"/node/take":       {"Take Node Self", controllers.TakeNode, GP, false}, //  和前端的那部分一毛一样
		"/node/list":       {"List Node Self", controllers.ListNode, GP, false},
		"/node/admin/list": {"List Node All", controllers.ListNodeAdmin, GP, true}, // 管理员查看其他用户节点

		// 内容操作
		// start review in 2019/5/15
		"/content/create":              {"Create Content Self", controllers.CreateContent, POST, false},                             // 创建文章内容(必须归属一个节点)
		"/content/update/seo":          {"Update Content Self Seo", controllers.UpdateSeoOfContent, POST, false},                    // 更新内容SEO
		"/content/update/image":        {"Update Content Self Image", controllers.UpdateImageOfContent, POST, false},                // 更新内容图片
		"/content/update/status":       {"Update Content Self Status", controllers.UpdateStatusOfContent, POST, false},              // 更新内容的状态，如设置隐藏
		"/content/admin/update/status": {"Update Content All Status", controllers.UpdateStatusOfContentAdmin, POST, true},           // 超级管理员修改文章，比如禁用或者逻辑删除/恢复文章
		"/content/update/node":         {"Update Content Self Node", controllers.UpdateNodeOfContent, POST, false},                  // 更改内容的节点，顺便需要重新排序
		"/content/update/top":          {"Update Content Self Top", controllers.UpdateTopOfContent, POST, false},                    // 设置内容的置顶与否
		"/content/update/comment":      {"Update Content Self Comment", controllers.UpdateCommentOfContent, POST, false},            // 设置内容可以评论与否
		"/content/update/password":     {"Update Content Self Password", controllers.UpdatePasswordOfContent, POST, false},          // 更改内容的密码保护
		"/content/update/info":         {"Update Content Self Info", controllers.UpdateInfoOfContent, POST, false},                  // 更新内容标题和内容
		"/content/sort":                {"Sort Content Self", controllers.SortContent, POST, false},                                 // 对内容进行拖曳排序
		"/content/publish":             {"Publish Content Self", controllers.PublishContent, POST, false},                           // 将预览刷进另外一个字段
		"/content/restore":             {"Restore Content Self", controllers.RestoreContent, POST, false},                           // 恢复历史，刷回来
		"/content/rubbish":             {"Sent Content Self To Rubbish", controllers.SentContentToRubbish, POST, false},             // 一般回收站
		"/content/recycle":             {"Sent Rubbish Content Self To Origin", controllers.ReCycleOfContentInRubbish, POST, false}, // 一般回收站恢复
		"/content/delete":              {"Delete Content Self Real", controllers.ReallyDeleteContent, POST, false},                  // 逻辑删除文章 已经修正为真删除

		// start review in 2019/5/16
		"/content/take":               {"Take Content Self", controllers.TakeContent, GP, false},                                 // 获取文章内容
		"/content/admin/take":         {"Take Content Admin", controllers.TakeContentAdmin, GP, true},                            // 管理员获取文章内容
		"/content/history/take":       {"Take Content History Self", controllers.TakeContentHistory, GP, false},                  // 获取文章历史内容
		"/content/history/admin/take": {"Take Content History Admin", controllers.TakeContentHistoryAdmin, GP, true},             // 管理员获取文章历史内容
		"/content/list":               {"List Content Self", controllers.ListContent, GP, false},                                 // 列出文章
		"/content/admin/list":         {"List Content All", controllers.ListContentAdmin, GP, true},                              // 管理员列出文章，什么类型都可以
		"/content/history/list":       {"List Content History Self", controllers.ListContentHistory, GP, false},                  // 列出文章的历史记录
		"/content/history/admin/list": {"List Content History All", controllers.ListContentHistoryAdmin, GP, true},               // 管理员列出文章的历史纪录
		"/content/history/delete":     {"Delete Content History Self Real", controllers.ReallyDeleteHistoryContent, POST, false}, // 真删除历史内容

		// start at 2019/9
		"/content/cool": {"Cool the Content Self", controllers.CoolContent, GP, false}, // 点赞内容
		"/content/bad":  {"Bad the Content Self", controllers.BadContent, GP, false},   // 举报内容

		"/comment/create": {"Create the Comment Self", controllers.CreateComment, POST, false}, // 创建评论
		"/comment/delete": {"Delete the Comment Self", controllers.DeleteComment, POST, false}, // 删除评论，逻辑删除
		"/comment/take":   {"Take the Comment Self", controllers.TakeComment, GP, false},       // 获取评论
		"/comment/cool":   {"Cool the Comment Self", controllers.CoolComment, GP, false},       // 点赞评论
		"/comment/bad":    {"Bad the Comment Self", controllers.BadComment, GP, false},         // 举报评论

		// admin url
		"/comment/admin/list":          {"List the Comment Admin", controllers.ListComment, GP, true},          // 管理员列出评论
		"/comment/admin/update/status": {"Update the Comment Status Admin", controllers.ListComment, GP, true}, // 管理员评论违禁处理
	}
)

// home end.
func SetRouter(router *gin.Engine) {
	for url, app := range HomeRouter {
		for _, method := range app.Method {
			router.Handle(method, url, app.Func)
		}
	}
}

func SetAPIRouter(router *gin.RouterGroup, handles map[string]HttpHandle) {
	for url, app := range handles {
		for _, method := range app.Method {
			router.Handle(method, url, app.Func)
		}
	}
}
