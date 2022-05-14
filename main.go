/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 19:10:13
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:18:10
 */
package main

import (
	_ "ihome/models"
	_ "ihome/routers"
	_ "ihome/utils"
	"net/http"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	_ "github.com/beego/beego/v2/server/web/session/redis"
)

func main() {
	env := "beta"
	//env := os.Getenv("ENV_CLUSTER")
	ConfPath := ""
	if env == "online" {
		//线上，正式
		ConfPath = "./conf/online.conf"
	} else if env == "beta" {
		//测试
		ConfPath = "./conf/app.conf"
	} else {
		//开发
		ConfPath = "./conf/dev.conf"
	}
	beego.LoadAppConfig("ini", ConfPath)
	ignoreStaticPath()
	beego.Run()
}

func ignoreStaticPath() {
	beego.InsertFilter("/", beego.BeforeRouter, TransportentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransportentStatic)
	beego.InsertFilter("/api/*", beego.AfterExec, setSession)
	beego.InsertFilter("/api/v1.0/user/*", beego.BeforeExec, checkLogin)
	beego.InsertFilter("/api/v1.0/user", beego.BeforeExec, checkLogin)
}

func TransportentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
}
func setSession(ctx *context.Context) {
	id, ok := ctx.Input.Session("id").(int)
	if ok {
		ctx.Output.Session("id", id)
	}
}
func checkLogin(ctx *context.Context) {
	// fmt.Println(ctx.Request.Method, ctx.Request.URL.Path)
	if ctx.Request.Method == "POST" && ctx.Request.URL.Path == "/api/v1.0/user" {
		return
	}
	_, ok := ctx.Input.Session("id").(int)
	if !ok {
		resp := make(map[string]interface{})
		resp["errno"] = 4101
		resp["errmsg"] = "请先登录"
		ctx.Output.JSON(resp, false, false)

	}
}
