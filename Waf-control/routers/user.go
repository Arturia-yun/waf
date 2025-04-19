package routers

import (

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"Waf-control/models"
)

// ListUser 显示用户列表
func ListUser(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login")
		return
	}

	users, err := models.ListUser()
	if err != nil {
		flash.Error("获取用户列表失败: " + err.Error())
	}

	ctx.Data["Flash"] = flash
	ctx.Data["users"] = users
	ctx.Data["Title"] = "用户列表"
	ctx.Data["ActiveMenu"] = "user"
	ctx.HTML(200, "user/list") 
}

// UserJSON 返回用户列表的JSON格式
func UserJSON(ctx *macaron.Context, sess session.Store) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login")
		return
	}

	users, err := models.ListUser()
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "获取用户列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, users)
}

// NewUser 显示新建用户页面
func NewUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login")
		return
	}

	ctx.Data["ActiveMenu"] = "user"
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.Data["Title"] = "新建用户"
	ctx.HTML(200, "user/new")
}

// DoNewUser 处理新建用户请求
func DoNewUser(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login")
		return
	}

	username := ctx.Req.Form.Get("username")
	password := ctx.Req.Form.Get("password")

	if username == "" || password == "" {
		flash.Error("用户名和密码不能为空")
		ctx.Redirect("/admin/user/new")
		return
	}

	err := models.NewUser(username, password)
	if err != nil {
		flash.Error("创建用户失败: " + err.Error())
		ctx.Redirect("/admin/user/new")
		return
	}

	flash.Success("用户创建成功")
	ctx.Redirect("/admin/user")
}

// EditUser 显示编辑用户页面
func EditUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login")
		return
	}

	id := ctx.ParamsInt64(":id")
	user, err := models.GetUserById(id)
	if err != nil {
		flash.Error("获取用户信息失败: " + err.Error())
		ctx.Redirect("/admin/user")
		return
	}

	if user == nil {
		flash.Error("用户不存在")
		ctx.Redirect("/admin/user")
		return
	}

	ctx.Data["csrf_token"] = x.GetToken()
	ctx.Data["user"] = user
	ctx.Data["Title"] = "编辑用户"
	ctx.Data["ActiveMenu"] = "user"
	ctx.HTML(200, "user/edit")
}

// DoEditUser 处理编辑用户请求
func DoEditUser(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login")
		return
	}

	id := ctx.ParamsInt64(":id")
	username := ctx.Req.Form.Get("username")
	password := ctx.Req.Form.Get("password")

	if username == "" {
		flash.Error("用户名不能为空")
		ctx.Redirect("/admin/user/edit/" + ctx.Params(":id"))
		return
	}

	err := models.UpdateUser(id, username, password)
	if err != nil {
		flash.Error("更新用户失败: " + err.Error())
		ctx.Redirect("/admin/user/edit/" + ctx.Params(":id"))
		return
	}

	flash.Success("用户更新成功")
	ctx.Redirect("/admin/user")
}

// DelUser 处理删除用户请求
func DelUser(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	if sess.Get("uid") == nil {
		ctx.Redirect("/login")
		return
	}

	id := ctx.ParamsInt64(":id")
	
	// 不允许删除当前登录用户
	if id == sess.Get("uid").(int64) {
		flash.Error("不能删除当前登录的用户")
		ctx.Redirect("/admin/user")
		return
	}

	err := models.DelUser(id)
	if err != nil {
		flash.Error("删除用户失败: " + err.Error())
	} else {
		flash.Success("用户删除成功")
	}

	ctx.Redirect("/admin/user")
}