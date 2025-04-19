package routers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-macaron/csrf"  
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"Waf-control/models"
	"Waf-control/modules/util"
	"Waf-control/setting"
)

// Admin 管理首页
func Admin(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	sites, err := models.ListSite()
	if err != nil {
		flash.Error("获取站点列表失败: " + err.Error())
	}

	// 确保 Flash 消息正确传递给模板
	ctx.Data["Flash"] = flash
	ctx.Data["Sites"] = sites
	ctx.Data["Title"] = "站点列表"
	ctx.HTML(200, "site/list")
}

// NewSite 新建站点页面
func NewSite(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	// 不需要再设置CSRF令牌，因为全局中间件已经设置
	ctx.Data["Title"] = "新建站点"
	ctx.HTML(200, "site/new")
}

// DoNewSite 处理新建站点
func DoNewSite(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	// 移除 x csrf.CSRF 参数
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}


	siteName := ctx.Req.Form.Get("site_name")
	portStr := ctx.Req.Form.Get("port")
	backendAddrStr := ctx.Req.Form.Get("backend_addr")
	unrealAddrStr := ctx.Req.Form.Get("unreal_addr")
	ssl := ctx.Req.Form.Get("ssl")
	debugLevel := ctx.Req.Form.Get("debug_level")

	if siteName == "" || portStr == "" || backendAddrStr == "" {
		flash.Error("站点名称、端口和后端地址不能为空")
		ctx.Redirect("/admin/site/new")
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		flash.Error("端口必须是数字")
		ctx.Redirect("/admin/site/new")
		return
	}

	backendAddr := strings.Split(backendAddrStr, ",")
	var unrealAddr []string
	if unrealAddrStr != "" {
		unrealAddr = strings.Split(unrealAddrStr, ",")
	}

	err = models.NewSite(siteName, port, backendAddr, unrealAddr, ssl, debugLevel)
	if err != nil {
		flash.Error("创建站点失败: " + err.Error())
		ctx.Redirect("/admin/site/new")
		return
	}

	flash.Success("站点创建成功")
	ctx.Redirect("/admin/site")
}

// EditSite 编辑站点页面
func EditSite(ctx *macaron.Context, sess session.Store, x csrf.CSRF, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		flash.Error("无效的站点ID")
		ctx.Redirect("/admin/site")
		return
	}

	site, err := models.GetSiteById(id)
	if err != nil {
		flash.Error("获取站点失败: " + err.Error())
		ctx.Redirect("/admin/site")
		return
	}

	if site == nil {
		flash.Error("站点不存在")
		ctx.Redirect("/admin/site")
		return
	}

	ctx.Data["csrf_token"] = x.GetToken()  // 添加这一行设置CSRF令牌
	ctx.Data["Site"] = site
	ctx.Data["BackendAddr"] = strings.Join(site.BackendAddr, ",")
	ctx.Data["UnrealAddr"] = strings.Join(site.UnrealAddr, ",")
	ctx.Data["Title"] = "编辑站点"
	ctx.HTML(200, "site/edit")
}

// DoEditSite 处理编辑站点
func DoEditSite(ctx *macaron.Context, sess session.Store, x csrf.CSRF, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		flash.Error("无效的站点ID")
		ctx.Redirect("/admin/site")
		return
	}

	siteName := ctx.Req.Form.Get("site_name")
	portStr := ctx.Req.Form.Get("port")
	backendAddrStr := ctx.Req.Form.Get("backend_addr")
	unrealAddrStr := ctx.Req.Form.Get("unreal_addr")
	ssl := ctx.Req.Form.Get("ssl")
	debugLevel := ctx.Req.Form.Get("debug_level")

	if siteName == "" || portStr == "" || backendAddrStr == "" {
		flash.Error("站点名称、端口和后端地址不能为空")
		ctx.Redirect("/admin/site/edit/" + idStr)
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		flash.Error("端口必须是数字")
		ctx.Redirect("/admin/site/edit/" + idStr)
		return
	}

	backendAddr := strings.Split(backendAddrStr, ",")
	var unrealAddr []string
	if unrealAddrStr != "" {
		unrealAddr = strings.Split(unrealAddrStr, ",")
	}

	err = models.UpdateSite(id, siteName, port, backendAddr, unrealAddr, ssl, debugLevel)
	if err != nil {
		flash.Error("更新站点失败: " + err.Error())
		ctx.Redirect("/admin/site/edit/" + idStr)
		return
	}

	flash.Success("站点更新成功")
	ctx.Redirect("/admin/site")
}

// DelSite 删除站点
func DelSite(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		flash.Error("无效的站点ID")
		ctx.Redirect("/admin/site")
		return
	}

	err = models.DelSite(id)
	if err != nil {
		flash.Error("删除站点失败: " + err.Error())
		ctx.Redirect("/admin/site")
		return
	}

	flash.Success("站点删除成功")
	ctx.Redirect("/admin/site")
}

// SyncSite 同步所有站点配置
func SyncSite(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	timestamp := time.Now().Unix()
	hash := util.MakeMd5(setting.AppKey + util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey)))

	log.Println("开始同步站点配置到 WAF 服务器...")
	log.Println("API 服务器列表:", setting.APIServers)

	for _, server := range setting.APIServers {
		server = strings.TrimSpace(server)
		
		url := fmt.Sprintf("http://%s:%v/api/site/sync/?hash=%v&timestamp=%v", 
			server, setting.HTTPPort, hash, timestamp)
		log.Println("尝试连接到 WAF 服务器:", url)
		
		resp, err := http.Get(url)
		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			
			log.Println("WAF 服务器响应:", string(body))
			if err != nil {
				log.Println("读取响应失败:", err.Error())
				flash.Error("读取响应失败: " + err.Error())
			} else {
				flash.Success("同步成功: " + string(body))
			}
		} else {
			log.Println("连接 WAF 服务器失败:", err.Error())
			flash.Error("连接 WAF 服务器失败: " + err.Error())
		}
	}

	ctx.Redirect("/admin/site")
}

// SyncSiteById 同步指定站点配置
func SyncSiteById(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		flash.Error("无效的站点ID")
		ctx.Redirect("/admin/site")
		return
	}

	timestamp := time.Now().Unix()
	hash := util.MakeMd5(setting.AppKey + util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey)))

	log.Println("开始同步站点 ID:", id, "到 WAF 服务器...")
	log.Println("API 服务器列表:", setting.APIServers)

	for _, server := range setting.APIServers {
		server = strings.TrimSpace(server)
		url := fmt.Sprintf("http://%s:%v/api/site/sync/?id=%v&hash=%v&timestamp=%v", 
			server, setting.HTTPPort, id, hash, timestamp)
		log.Println("尝试连接到 WAF 服务器:", url)

		resp, err := http.Get(url)
		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			log.Println("WAF 服务器响应:", string(body), err)
			flash.Success(string(body))
		} else {
			log.Println("连接 WAF 服务器失败:", err.Error())
			flash.Error(err.Error())
		}
	}

	ctx.Redirect("/admin/site")
}

// SiteJSON 获取站点JSON配置
func SiteJSON(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") == nil {
		ctx.JSON(403, map[string]interface{}{
			"status":  0,
			"message": "未授权访问",
		})
		return
	}

	sites, err := models.ListSite()
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"status":  0,
			"message": "获取站点列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"status":  1,
		"message": "获取站点列表成功",
		"data":    sites,
	})
}
