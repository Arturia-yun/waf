package routers

import (
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"Waf-control/models"
)

// Login 登录页面
func Login(ctx *macaron.Context, sess session.Store, x csrf.CSRF, flash *session.Flash) {
    if sess.Get("uid") != nil {
        ctx.Redirect("/admin/index")
        return
    }
    
    ctx.Data["csrf_token"] = x.GetToken()
    ctx.Data["Title"] = "登录"
    // 删除 Layout 相关设置
    ctx.HTML(200, "user/login")
}

// DoLogin 处理登录请求
func DoLogin(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	username := ctx.Req.Form.Get("username")
	password := ctx.Req.Form.Get("password")
	
	if username == "" || password == "" {
		flash.Error("用户名和密码不能为空")
		ctx.Redirect("/login")  
		return
	}
	
	user, ok := models.ValidateUser(username, password)
	if !ok {
		flash.Error("用户名或密码错误")
		ctx.Redirect("/login") 
		return
	}
	
	// 设置会话
	sess.Set("uid", user.Id)
	sess.Set("username", user.Username)
	
	flash.Success("登录成功")
	ctx.Redirect("/admin/index")
}

// Logout 退出登录
func Logout(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	sess.Delete("uid")
	sess.Delete("username")
	flash.Success("已退出登录")
	ctx.Redirect("/login")  
}