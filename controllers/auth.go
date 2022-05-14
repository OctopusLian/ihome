/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 23:10:23
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:32:59
 */
package controllers

import (
	"encoding/json"
	"regexp"

	"ihome/models"

	"github.com/beego/beego/orm"
	"github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	web.Controller
}

func (this *AuthController) Get() {
	id := this.GetSession("id").(int)
	user := models.User{Id: id}
	o := orm.NewOrm()
	o.Read(&user)
	resp := make(map[string]interface{})
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	data := make(map[string]interface{})
	data["real_name"] = user.Real_name
	data["id_card"] = user.Id_card
	resp["data"] = data
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *AuthController) retData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *AuthController) Post() {
	params := make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	resp := make(map[string]interface{})
	defer this.retData(resp)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return

	}

	name := params["real_name"].(string)
	idCard := params["id_card"].(string)
	if name == "" || idCard == "" {
		resp["errno"] = 4001
		resp["errmsg"] = "姓名和身份证号不能为空"
		return
	}
	pat := `^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`
	ok, _ := regexp.Match(pat, []byte(idCard))
	if !ok {
		resp["errno"] = 4001
		resp["errmsg"] = "请输入正确的身份证号"
		return
	}
	id := this.GetSession("id").(int)
	user := models.User{Id: id}
	o := orm.NewOrm()
	o.Begin()
	o.Read(&user)
	if user.Real_name != "" {
		resp["errno"] = 4001
		resp["errmsg"] = "已实名"
		o.Rollback()
		return
	}
	user.Id_card = idCard
	user.Real_name = name
	o.Update(&user, "Real_name", "Id_card")
	o.Commit()
	resp["errno"] = 0
	resp["errmsg"] = "成功"
}
