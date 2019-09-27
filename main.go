/*
	版权所有，侵权必究
	署名-非商业性使用-禁止演绎 4.0 国际
	你可以在教育用途下使用该代码，但是禁止公司或个人用于商业用途!

	All right reserved
	Attribution-NonCommercial-NoDerivatives 4.0 International
	You can use it for education only but can't make profits for any companies and individuals!

	The Door of the program

	FaFa awesome!
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
	// Global path of config file
	configFile string

	// Auto create database and tables when set true
	createTable bool

	// Email debug will not send email
	mailDebug bool

	// Skip some admin Auth when debug
	canSkipAuth bool

	// Record the history of content when edit or publish
	historyRecord bool

	// Time zone offset the utc, default is 8,Beijing
	timeZone int64

	// Auto ban the content or comment
	autoBan bool

	// Beyond and autoBan is true will Ban it!
	banTime int64

	// Single Login
	singleLogin bool

	// Login Session expire time
	sessionExpireTime int64
)

// Parse flag when init
// Those variables will not config in file
func init() {
	// Default read ./config.json
	flag.StringVar(&configFile, "config", "./config.json", "Config file")

	// Auto init db, the second time can set false
	flag.BoolVar(&createTable, "init_db", true, "Init create db table")

	flag.Int64Var(&timeZone, "time_zone", 8, "Time zone offset the utc")
	flag.BoolVar(&autoBan, "auto_ban", false, "Auto ban the content or comment")
	flag.Int64Var(&banTime, "ban_time", 10, "Content or comment will be ban in how much bad's time")
	flag.BoolVar(&historyRecord, "history_record", true, "Content history can be record")
	flag.BoolVar(&singleLogin, "single_login", false, "User can only single point login")
	flag.Int64Var(&sessionExpireTime, "session_expire_time", 7*3600*24, "Login session expire second time, token will destroy after this time")

	// When in production, please set to all false
	flag.BoolVar(&mailDebug, "email_debug", false, "Email debug")
	flag.BoolVar(&canSkipAuth, "auth_skip_debug", false, "Auth skip debug")

	flag.Parse()
}

// Init the URL resource, some admin url put inside a map will save a lot of time
func initResource() (adminUrl map[string]int64) {
	adminUrl = make(map[string]int64)
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

		// Exist will put in map, otherwise save in db then put in map
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
	return adminUrl
}

// The Beauty Main
// I'm FaFa
func main() {

	// Package var init
	mail.Debug = mailDebug
	controllers.AuthDebug = canSkipAuth
	controllers.TimeZone = timeZone
	controllers.BadTime = banTime
	controllers.AutoBan = autoBan
	controllers.SingleLogin = singleLogin
	controllers.SessionExpireTime = sessionExpireTime
	model.HistoryRecord = historyRecord

	var err error

	// Init global config
	err = server.InitConfig(configFile)
	if err != nil {
		panic(err)
	}

	// Init log
	flog.InitLog(config.FaFaConfig.DefaultConfig.LogPath)

	// Log can set to debug, i suggest not to change it
	if config.FaFaConfig.DefaultConfig.LogDebug {
		flog.SetLogLevel("DEBUG")
	}

	welcome()
	flog.Log.Debugf("Hi! Config is %#v", config.FaFaConfig)

	// Init db
	err = server.InitRdb(config.FaFaConfig.DbConfig)
	if err != nil {
		panic(err)
	}

	// Init session
	err = session.InitSession(config.FaFaConfig.SessionConfig)
	if err != nil {
		panic(err)
	}

	// Auto create db table
	if createTable {
		model.CreateTable([]interface{}{
			model.User{},           // User Table
			model.Group{},          // User Group, every user can assign a group
			model.Resource{},       // Url Resource, if user not own those will be refuse to auth
			model.GroupResource{},  // Resource will be assign to group
			model.Content{},        // Content Table, very import
			model.ContentCool{},    // Content Cool, user can cool your content
			model.ContentBad{},     // Content Bad, user can bad your content and if auto ban, your content will be ban
			model.ContentHistory{}, // Content History, when publish or edit a content, and you set history record, emm, save it
			model.ContentNode{},    // Contents' Node, every content must belong to a node
			model.File{},           // File Table, your picture file and some will save in.
			model.Comment{},        // Comment Table, comment for content, comment for comment
			model.CommentCool{},    // Like the Content Cool
			model.CommentBad{},     // Like the Content Bad
			model.Relation{},       // Who follow who
			model.Message{},        // Message inside
			//model.Log{},            // Log Table, not use
		})
	}

	controllers.AdminUrl = initResource()

	// Server Run
	engine := server.Server()

	// Storage static API
	engine.Static("/storage", config.FaFaConfig.DefaultConfig.StoragePath)
	engine.Static("/storage_x", config.FaFaConfig.DefaultConfig.StoragePath+"_x")

	// Web welcome home!
	router.SetRouter(engine)

	// V1 API, will may be change to V2...
	v1 := engine.Group("/v1")
	v1.Use(controllers.AuthFilter)

	// Router Set
	router.SetAPIRouter(v1, router.V1Router)

	flog.Log.Noticef("Server run in %s", config.FaFaConfig.DefaultConfig.WebPort)
	err = engine.Run(config.FaFaConfig.DefaultConfig.WebPort)
	if err != nil {
		panic(err)
	}
}

func welcome() {
	flog.Log.Notice("Hi! FaFa CMS! A Nice CMS.")
	s := `
███████╗ █████╗ ███████╗ █████╗  ██████╗███╗   ███╗███████╗
██╔════╝██╔══██╗██╔════╝██╔══██╗██╔════╝████╗ ████║██╔════╝
█████╗  ███████║█████╗  ███████║██║     ██╔████╔██║███████╗
██╔══╝  ██╔══██║██╔══╝  ██╔══██║██║     ██║╚██╔╝██║╚════██║
██║     ██║  ██║██║     ██║  ██║╚██████╗██║ ╚═╝ ██║███████║
╚═╝     ╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚═╝     ╚═╝╚══════╝`
	flog.Log.Noticef("\n%s-v%s_%s\n", s, config.Version, util.BuildTime())
}
