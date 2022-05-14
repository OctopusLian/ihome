/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 23:06:19
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:31:33
 */
package controllers

import (
	"ihome/models"

	"github.com/beego/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type AreaController struct {
	beego.Controller
}

func (this *AreaController) GetArea() {
	resp := make(map[string]interface{})
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	o := orm.NewOrm()
	var areas []*models.Area
	num, err := o.QueryTable("Area").All(&areas)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		this.Data["json"] = &resp
		this.ServeJSON()
		return
	}
	if num == 0 {
		resp["errno"] = 4001
		resp["errmsg"] = "没有查询到数据"
		this.Data["json"] = &resp
		this.ServeJSON()
		return
	} else {
		resp["data"] = areas
	}
	this.Data["json"] = &resp
	this.ServeJSON()
}
