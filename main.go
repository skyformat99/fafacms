/*
	版权所有，侵权必究
	署名-非商业性使用-禁止演绎 4.0 国际
	警告： 以下的代码版权归属hunterhug，请不要传播或修改代码
	你可以在教育用途下使用该代码，但是禁止公司或个人用于商业用途(在未授权情况下不得用于盈利)
	商业授权请联系邮箱：gdccmcm14@live.com QQ:459527502

	All right reserved
	Attribution-NonCommercial-NoDerivatives 4.0 International
	Notice: The following code's copyright by hunterhug, Please do not spread and modify.
	You can use it for education only but can't make profits for any companies and individuals!
	For more information on commercial licensing please contact hunterhug.
	Ask for commercial licensing please contact Mail:gdccmcm14@live.com Or QQ:459527502

    2019-4-24：

	程序主入口
	花花CMS是一个内容管理系统，代码尽可能地补充必要注释，方便后人协作
**/
package main

import (
	"flag"
	"fmt"
	"github.com/hunterhug/fafacms/core/config"
	"github.com/hunterhug/fafacms/core/controllers"
	"github.com/hunterhug/fafacms/core/flog"
	"github.com/hunterhug/fafacms/core/model"
	"github.com/hunterhug/fafacms/core/router"
	"github.com/hunterhug/fafacms/core/server"
	"github.com/hunterhug/fafacms/core/session"
	"github.com/hunterhug/fafacms/core/util"
	"github.com/hunterhug/fafacms/core/util/mail"
	"time"
)

var (
	version = "1.0.0"

	// 全局配置文件路径
	configFile string

	// 是否创建数据库表
	createTable bool

	// 开发时每次都发邮件的形式不好，可以先调试模式
	mailDebug bool

	// 跳过授权，某些超级管理接口需要绑定组和路由，可以先开调试模式
	canSkipAuth bool

	// 是否内容刷新进历史表进行保存
	historyRecord bool
)

// 初始化时解析命令行，辅助程序
// 这些调试参数不置于文件配置中
func init() {
	// 默认读取本路径下 ./config.json 配置
	flag.StringVar(&configFile, "config", "./config.json", "config file")

	// 正式部署时，请全部设置为 false
	flag.BoolVar(&createTable, "init_db", true, "create db table")
	flag.BoolVar(&mailDebug, "email_debug", false, "Email debug")
	flag.BoolVar(&canSkipAuth, "auth_skip_debug", false, "Auth skip debug")
	flag.BoolVar(&historyRecord, "history_record", false, "Content history record")

	flag.Parse()
}

// 初始化URL资源
func initResource() (adminUrl map[string]int) {
	adminUrl = make(map[string]int)
	for url, handler := range router.V1Router {
		if !handler.Admin {
			continue
		}
		r := new(model.Resource)
		url1 := fmt.Sprintf("/v1%s", url)
		r.UrlHash, _ = util.Sha256([]byte(url1))
		r.Admin = true
		exist, err := r.GetRaw()
		if err != nil {
			panic(err)
		}

		if exist {
			adminUrl[url1] = r.Id
			continue
		} else {
			r := new(model.Resource)
			r.Url = url1
			r.UrlHash, _ = util.Sha256([]byte(url1))
			r.Name = handler.Name
			r.Describe = handler.Name
			r.Admin = handler.Admin
			r.CreateTime = time.Now().Unix()
			err := r.InsertOne()
			if err != nil {
				panic(err)
			}
			adminUrl[url1] = r.Id
		}
	}
	//fmt.Printf("admin url:%#v\n", adminUrl)
	return adminUrl
}

// 入口
// 欢迎查看优美代码，我是花花
func main() {

	// 将调试参数跨包注入
	mail.Debug = mailDebug
	controllers.AuthDebug = canSkipAuth
	model.HistoryRecord = historyRecord

	var err error

	// 初始化全局配置
	err = server.InitConfig(configFile)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	flog.InitLog(config.FafaConfig.DefaultConfig.LogPath)

	// 如果全局调试，那么所有DEBUG以上级别日志将会打印
	// 实际情况下，最好设置为 true，
	if config.FafaConfig.DefaultConfig.LogDebug {
		flog.SetLogLevel("DEBUG")
	}

	welcome()
	flog.Log.Debugf("Hi! Config is %#v", config.FafaConfig)

	// 初始化数据库连接
	err = server.InitRdb(config.FafaConfig.DbConfig)
	if err != nil {
		panic(err)
	}

	// 初始化网站Session存储
	err = session.InitSession(config.FafaConfig.SessionConfig)
	if err != nil {
		panic(err)
	}

	// 创建数据库表，需要先手动创建DB
	if createTable {
		model.CreateTable([]interface{}{
			model.User{},           // 用户表
			model.Group{},          // 用户组表，用户可以拥有一个组
			model.Resource{},       // 资源表，主要为需要管理员权限的路由服务
			model.GroupResource{},  // 组可以被分配资源
			model.Content{},        // 内容表
			model.ContentHistory{}, // 内容历史表
			model.ContentNode{},    // 内容节点表，内容必须拥有一个节点
			model.File{},           // 文件表
			//model.Comment{},        // 评论表
			//model.Log{},            // 日志表
		})
	}

	controllers.AdminUrl = initResource()

	// Server Run
	engine := server.Server()

	// Storage static API
	engine.Static("/storage", config.FafaConfig.DefaultConfig.StoragePath)
	engine.Static("/storage_x", config.FafaConfig.DefaultConfig.StoragePath+"_x")

	// Web welcome home!
	router.SetRouter(engine)

	// V1 API, will may be change to V2...
	v1 := engine.Group("/v1")
	v1.Use(controllers.AuthFilter)

	// Router Set
	router.SetAPIRouter(v1, router.V1Router)

	flog.Log.Noticef("Server run in %s", config.FafaConfig.DefaultConfig.WebPort)
	err = engine.Run(config.FafaConfig.DefaultConfig.WebPort)
	if err != nil {
		panic(err)
	}
}

func welcome() {
	flog.Log.Notice("Hi! FaFa CMS!")
	s := `
███████╗ █████╗ ███████╗ █████╗  ██████╗███╗   ███╗███████╗
██╔════╝██╔══██╗██╔════╝██╔══██╗██╔════╝████╗ ████║██╔════╝
█████╗  ███████║█████╗  ███████║██║     ██╔████╔██║███████╗
██╔══╝  ██╔══██║██╔══╝  ██╔══██║██║     ██║╚██╔╝██║╚════██║
██║     ██║  ██║██║     ██║  ██║╚██████╗██║ ╚═╝ ██║███████║
╚═╝     ╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚═╝     ╚═╝╚══════╝`

	flog.Log.Noticef("\n%s-v%s_%s\n", s, version, util.BuildTime())
}
