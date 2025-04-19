package routers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/macaron.v1"

	"Waf-control/models"
	"Waf-control/modules/util"
	"Waf-control/setting"
)

// APIResponse API响应结构
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SyncRuleAPI 处理规则同步API请求
func SyncRuleAPI(ctx *macaron.Context) {
	resp := &APIResponse{
		Status:  0,
		Message: "同步失败",
	}

	// 验证请求
	hash := ctx.Query("hash")
	timestampStr := ctx.Query("timestamp")
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		resp.Message = "无效的时间戳"
		ctx.JSON(200, resp)
		return
	}

	// 检查时间戳是否过期（5分钟内有效）
	now := time.Now().Unix()
	if now-timestamp > 300 {
		resp.Message = "请求已过期"
		ctx.JSON(200, resp)
		return
	}

	// 验证哈希
	expectedHash := util.MakeMd5(setting.AppKey + util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey)))
	if hash != expectedHash {
		resp.Message = "无效的请求签名"
		ctx.JSON(200, resp)
		return
	}

	// 获取所有规则类型
	ruleTypes := models.GetAllRuleTypes()

	// 创建规则目录
	err = os.MkdirAll(setting.RulePath, 0755)
	if err != nil {
		resp.Message = "创建规则目录失败: " + err.Error()
		ctx.JSON(200, resp)
		return
	}

	// 为每种规则类型生成JSON文件
	for _, ruleType := range ruleTypes {
		rules, err := models.GetRulesByType(ruleType)
		if err != nil {
			resp.Message = "获取规则失败: " + err.Error()
			ctx.JSON(200, resp)
			return
		}

		// 生成JSON文件
		jsonData, err := json.Marshal(rules)
		if err != nil {
			resp.Message = "序列化规则失败: " + err.Error()
			ctx.JSON(200, resp)
			return
		}

		// 写入文件
		ruleFile := filepath.Join(setting.RulePath, ruleType+".rule")
		err = ioutil.WriteFile(ruleFile, jsonData, 0644)
		if err != nil {
			resp.Message = "写入规则文件失败: " + err.Error()
			ctx.JSON(200, resp)
			return
		}
	}

	// 重新加载Nginx配置
	err = util.ReloadNginx()
	if err != nil {
		resp.Message = "重新加载Nginx失败: " + err.Error()
		ctx.JSON(200, resp)
		return
	}

	resp.Status = 1
	resp.Message = "规则同步成功"
	ctx.JSON(200, resp)
}

// SyncSiteAPI 处理站点同步API请求
func SyncSiteAPI(ctx *macaron.Context) {
	log.Println("收到站点同步API请求")

	// 在开始处理前，确保使用正确的路径
	if setting.IsWindows {
		log.Println("当前是 Windows 环境，使用 Windows 路径格式")
		setting.NginxVhosts = filepath.FromSlash(setting.NginxVhosts)
		setting.RulePath = filepath.FromSlash(setting.RulePath)
	}

	// 获取请求参数
	timestamp := ctx.Query("timestamp")
	hash := ctx.Query("hash")
	id := ctx.Query("id")

	log.Println("请求参数:", "hash=", hash, "timestamp=", timestamp, "id=", id)

	// 验证哈希
	if util.MakeMd5(setting.AppKey+util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey))) == hash {
		// 根据ID获取站点
		var sites []models.Site
		var err error

		if id != "" {
			Id, err := strconv.Atoi(id)
			log.Println("同步指定站点，ID:", Id, err)
			if err == nil {
				sites, _ = models.ListSiteById(int64(Id))
			} else {
				log.Println("ID转换错误:", err)
				sites, _ = models.ListSite()
			}
		} else {
			log.Println("同步所有站点")
			sites, err = models.ListSite()
		}

		log.Println("获取站点列表:", sites, err)
		log.Println("虚拟主机目录:", setting.NginxVhosts)

		// 在生成配置前添加操作系统类型
		ctx.Data["IsWindows"] = setting.IsWindows

		// 为每个站点生成配置
		for _, site := range sites {
			ctx.Data["site"] = site
			proxyConfig, err := ctx.HTMLString("proxy", ctx.Data)
			log.Println("生成配置:", err)
			if err != nil {
				log.Println("生成配置错误:", err)
				ret := util.Message{Status: 1, Message: "生成配置失败: " + err.Error()}
				ctx.JSON(200, &ret)
				return
			}

			log.Println("生成的配置长度:", len(proxyConfig))

			err = util.WriteNginxConf(proxyConfig, site.SiteName, setting.NginxVhosts)
			if err != nil {
				log.Println("写入配置错误:", err)
				ret := util.Message{Status: 1, Message: "写入配置失败: " + err.Error()}
				ctx.JSON(200, &ret)
				return
			}
		}

		// 重新加载Nginx配置
		log.Println("开始重新加载Nginx配置")
		if util.ReloadNginx() == nil {
			ret := util.Message{Status: 1, Message: "站点同步成功"}
			ctx.JSON(200, &ret)
		} else {
			ret := util.Message{Status: 0, Message: "重载Nginx配置失败"}
			ctx.JSON(200, &ret)
		}
	} else {
		log.Println("验证失败")
		ret := util.Message{Status: 2, Message: "无效的请求签名"}
		ctx.JSON(200, &ret)
	}
}

// SyncSiteByIdAPI 处理指定站点同步API请求
func SyncSiteByIdAPI(ctx *macaron.Context) {
	resp := &APIResponse{
		Status:  0,
		Message: "同步失败",
	}

	// 在开始处理前，确保使用正确的路径
	if setting.IsWindows {
		log.Println("当前是 Windows 环境，使用 Windows 路径格式")
		setting.NginxVhosts = filepath.FromSlash(setting.NginxVhosts)
		setting.RulePath = filepath.FromSlash(setting.RulePath)
	}

	// 验证请求
	hash := ctx.Query("hash")
	timestampStr := ctx.Query("timestamp")
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		resp.Message = "无效的时间戳"
		ctx.JSON(200, resp)
		return
	}

	// 检查时间戳是否过期（5分钟内有效）
	now := time.Now().Unix()
	if now-timestamp > 300 {
		resp.Message = "请求已过期"
		ctx.JSON(200, resp)
		return
	}

	// 验证哈希
	expectedHash := util.MakeMd5(setting.AppKey + util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey)))
	if hash != expectedHash {
		resp.Message = "无效的请求签名"
		ctx.JSON(200, resp)
		return
	}

	// 获取站点ID
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Message = "无效的站点ID"
		ctx.JSON(200, resp)
		return
	}

	// 获取指定站点
	site, err := models.GetSiteById(id)
	if err != nil {
		resp.Message = "获取站点失败: " + err.Error()
		ctx.JSON(200, resp)
		return
	}

	if site == nil {
		resp.Message = "站点不存在"
		ctx.JSON(200, resp)
		return
	}

	// 生成站点配置
	ctx.Data["site"] = site
	proxyConfig, err := ctx.HTMLString("proxy", ctx.Data)
	if err != nil {
		resp.Message = "生成站点配置失败: " + err.Error()
		ctx.JSON(200, resp)
		return
	}

	err = util.WriteNginxConf(proxyConfig, site.SiteName, setting.NginxVhosts)
	if err != nil {
		resp.Message = "写入站点配置失败: " + err.Error()
		ctx.JSON(200, resp)
		return
	}

	// 重新加载Nginx配置
	err = util.ReloadNginx()
	if err != nil {
		resp.Message = "重新加载Nginx失败: " + err.Error()
		ctx.JSON(200, resp)
		return
	}

	resp.Status = 1
	resp.Message = "站点同步成功"
	ctx.JSON(200, resp)
}
