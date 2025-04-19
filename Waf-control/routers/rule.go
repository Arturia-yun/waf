package routers

import (
	"Waf-control/models"
	"Waf-control/modules/util"
	"Waf-control/setting"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

// ListRules 规则列表页面
func ListRules(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	rules, err := models.GetRules()
	if err != nil {
		flash.Error("获取规则列表失败: " + err.Error())
	}

	// 获取所有规则类型
	ruleTypes := models.GetAllRuleTypes()

	ctx.Data["Flash"] = flash
	ctx.Data["Rules"] = rules
	ctx.Data["RuleTypes"] = ruleTypes
	ctx.Data["Title"] = "规则列表"
	ctx.Data["ActiveMenu"] = "rule"
	ctx.Data["BreadcrumbTitle"] = "规则管理"
	ctx.HTML(200, "rule/list")
}

// NewRule 创建新规则页面
func NewRule(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	ruleType := ctx.Params("type")
	if ruleType == "" {
		flash.Error("规则类型不能为空")
		ctx.Redirect("/admin/rule/")
		return
	}

	ctx.Data["RuleType"] = ruleType
	ctx.Data["Title"] = "新建规则"
	ctx.HTML(200, "rule/new")
}

// DoNewRule 处理新规则提交
func DoNewRule(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	ruleType := ctx.Params("type")
	ruleItem := ctx.Req.Form.Get("rule_item")

	if ruleType == "" || ruleItem == "" {
		flash.Error("规则类型和内容不能为空")
		ctx.Redirect("/admin/rule/")
		return
	}

	rule := &models.Rule{
		RuleType: ruleType,
		RuleItem: ruleItem,
	}

	err := models.CreateRule(rule)
	if err != nil {
		flash.Error("创建规则失败: " + err.Error())
		ctx.Redirect("/admin/rule/new/" + ruleType)
		return
	}

	flash.Success("规则创建成功")
	ctx.Redirect("/admin/rule/")
}

// EditRule 编辑规则页面
func EditRule(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		flash.Error("无效的规则ID")
		ctx.Redirect("/admin/rule/")
		return
	}

	rule, err := models.GetRuleById(id)
	if err != nil {
		flash.Error("获取规则失败: " + err.Error())
		ctx.Redirect("/admin/rule/")
		return
	}

	ctx.Data["Rule"] = rule
	ctx.Data["Title"] = "编辑规则"
	ctx.HTML(200, "rule/edit")
}

// DoEditRule 处理编辑规则提交
func DoEditRule(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		flash.Error("无效的规则ID")
		ctx.Redirect("/admin/rule/")
		return
	}

	ruleItem := ctx.Req.Form.Get("rule_item")
	if ruleItem == "" {
		flash.Error("规则内容不能为空")
		ctx.Redirect("/admin/rule/edit/" + idStr)
		return
	}

	rule, err := models.GetRuleById(id)
	if err != nil {
		flash.Error("获取规则失败: " + err.Error())
		ctx.Redirect("/admin/rule/")
		return
	}

	rule.RuleItem = ruleItem
	err = models.UpdateRule(rule)
	if err != nil {
		flash.Error("更新规则失败: " + err.Error())
		ctx.Redirect("/admin/rule/edit/" + idStr)
		return
	}

	flash.Success("规则更新成功")
	ctx.Redirect("/admin/rule/")
}

// DelRule 删除规则
func DelRule(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login/")
		return
	}

	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		flash.Error("无效的规则ID")
		ctx.Redirect("/admin/rule/")
		return
	}

	err = models.DeleteRule(id)
	if err != nil {
		flash.Error("删除规则失败: " + err.Error())
		ctx.Redirect("/admin/rule/")
		return
	}

	flash.Success("规则删除成功")
	ctx.Redirect("/admin/rule/")
}

// SyncRule 同步规则到所有WAF服务器
func SyncRule(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") != nil {
		timestamp := time.Now().Unix()
		hash := util.MakeMd5(setting.AppKey + util.MakeMd5(fmt.Sprintf("%v%v", timestamp, setting.AppKey)))
		
		for _, server := range setting.APIServers {
			server = strings.TrimSpace(server)
			url := fmt.Sprintf("http://%s:%v/api/rule/sync/?hash=%v&timestamp=%v",
				server, setting.HTTPPort, hash, timestamp)
			log.Println(url)
			resp, err := http.Get(url)
			if err == nil {
				body, err := ioutil.ReadAll(resp.Body)
				log.Println(string(body), err)
				flash.Success(string(body))
			} else {
				flash.Error(err.Error())
			}
		}
		
		ctx.Redirect("/admin/rule/")
	} else {
		ctx.Redirect("/login/")
	}
}
