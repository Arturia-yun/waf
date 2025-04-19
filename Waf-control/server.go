package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
	"Waf-control/models"
	"Waf-control/routers"
	"Waf-control/setting"
)

func main() {
	// 加载配置
	setting.Load()

	// 初始化数据库
	models.Init()

	// 创建Web服务
	m := macaron.Classic()

	// 设置模板
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Directory:       "templates",
		IndentJSON:      true,
		HTMLContentType: "text/html",
		Delims:          macaron.Delims{Left: "{{", Right: "}}"},
		Charset:         "UTF-8",
		Extensions:      []string{".tmpl", ".html"}, // 支持的扩展名
	}))

	// 设置会话
	m.Use(session.Sessioner())

	// 设置CSRF保护
	m.Use(csrf.Csrfer())

	// 设置静态文件
	m.Use(macaron.Static("templates/public/css", macaron.StaticOptions{
		Prefix: "/css",
	}))

	m.Use(macaron.Static("templates/public/js", macaron.StaticOptions{
		Prefix: "/js",
	}))

	// 设置路由
	// 在设置路由之前添加
	m.Use(func(ctx *macaron.Context, flash *session.Flash) {
		ctx.Data["Flash"] = flash
	})

	// 添加全局CSRF令牌中间件
	m.Use(func(ctx *macaron.Context, x csrf.CSRF) {
		ctx.Data["csrf_token"] = x.GetToken()
	})

	// 登录路由
	m.Get("/", func(ctx *macaron.Context) {
		ctx.Redirect("/login")
	})
	m.Get("/login", routers.Login)
	m.Post("/login", csrf.Validate, routers.DoLogin)
	m.Get("/logout", routers.Logout)

	m.Group("/admin", func() {
		m.Get("/", func(ctx *macaron.Context) {
			ctx.Redirect("/admin/index")
		})
		m.Get("/index", routers.Admin)

		// 规则管理路由
		m.Group("/rule", func() {
			m.Get("", routers.ListRules)
			m.Get("/list", routers.ListRules)
			m.Get("/new/:type", routers.NewRule)
			m.Post("/new/:type", csrf.Validate, routers.DoNewRule) // 这里必须是 /new/:type
			m.Get("/edit/:id", routers.EditRule)
			m.Post("/edit/:id", csrf.Validate, routers.DoEditRule)
			m.Get("/del/:id", routers.DelRule)
			m.Get("/sync", routers.SyncRule)
		})

		// 站点管理路由
		m.Group("/site", func() {
			m.Get("", routers.Admin)
			m.Get("/list", routers.Admin)
			m.Get("/new", routers.NewSite)
			m.Post("/new", csrf.Validate, routers.DoNewSite)
			m.Get("/edit/:id", routers.EditSite)
			m.Post("/edit/:id", csrf.Validate, routers.DoEditSite)
			m.Get("/del/:id", routers.DelSite)
			m.Get("/sync", routers.SyncSite)
			m.Get("/sync/:id", routers.SyncSiteById)
			m.Get("/json", routers.SiteJSON)
		})

		// 用户管理路由
		m.Group("/user", func() {
			m.Get("", routers.ListUser)
			m.Get("/list", routers.ListUser)
			m.Get("/new", routers.NewUser)
			m.Post("/new", csrf.Validate, routers.DoNewUser)
			m.Get("/edit/:id", routers.EditUser)
			m.Post("/edit/:id", csrf.Validate, routers.DoEditUser)
			m.Get("/del/:id", routers.DelUser)
			m.Get("/json", routers.UserJSON)
		})
	})

	// API路由
	m.Group("/api", func() {
		m.Get("/site/sync", routers.SyncSiteAPI)
		m.Get("/rule/sync", routers.SyncRuleAPI)
	})

	// 启动服务
	log.Printf("服务已启动，监听地址: %s:%d\n", setting.HTTPAddr, setting.HTTPPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", setting.HTTPAddr, setting.HTTPPort), m))
}
