/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 23:11:36
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:37:45
 */
package controllers

import (
	"encoding/json"
	"fmt"

	"ihome/models"

	"github.com/beego/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type SessionController struct {
	beego.Controller
}

func (this *SessionController) Get() {
	resp := make(map[string]interface{})
	if this.GetSession("name") != nil {
		resp["errno"] = 0
		data := make(map[string]interface{})
		data["name"] = this.GetSession("name")
		resp["data"] = data
	} else {
		resp["errno"] = 4001
	}
	this.Data["json"] = &resp
	this.ServeJSON()
}

func (this *SessionController) Post() {
	params := make(map[string]interface{})
	resp := make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	if err != nil {
		fmt.Println(err)
		resp["errno"] = 4001
		resp["errmsg"] = "参数解析失败"
		this.Data["json"] = &resp
		this.ServeJSON()
		return
	}
	mobile := params["mobile"].(string)
	password := params["password"].(string)
	if mobile == "" || password == "" {
		resp["errno"] = 4001
		resp["errmsg"] = "账号密码均不能为空"
		this.Data["json"] = &resp
		this.ServeJSON()
		return
	}
	user := models.User{Mobile: mobile}
	o := orm.NewOrm()

	o.Read(&user, "Mobile")
	if user.Password_hash == password {
		resp["errno"] = 0
		resp["errmsg"] = "成功"
		this.SetSession("id", user.Id)
		this.SetSession("name", user.Name)
	} else {
		fmt.Println(password, "s")
		fmt.Println(user.Password_hash)
		resp["errno"] = 4001
		resp["errmsg"] = "账号或密码错误"
	}
	this.Data["json"] = &resp
	this.ServeJSON()

}

func (this *SessionController) Delete() {
	resp := make(map[string]interface{})
	err := this.DestroySession()
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	this.Data["json"] = resp
	this.ServeJSON()
}
