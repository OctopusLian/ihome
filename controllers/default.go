/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 23:11:14
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:54:42
 */
package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
