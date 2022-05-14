/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 23:08:59
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:39:01
 */
package controllers

import (
	"encoding/json"
	"fmt"
	"ihome/utils"
	"regexp"

	"ihome/models"

	"github.com/beego/beego"
	"github.com/beego/beego/logs"
	"github.com/beego/beego/orm"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Post() {
	params := make(map[string]interface{})
	resp := make(map[string]interface{})
	resp["errno"] = 4001
	resp["errmsg"] = "测试中"
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(params)
	rep := `/^1(3[0-9]|4[01456879]|5[0-35-9]|6[2567]|7[0-8]|8[0-9]|9[0-35-9])\d{8}$/`
	ok, _ := regexp.Match(rep, []byte(params["mobile"].(string)))
	if !ok {
		resp["errno"] = 4001
		resp["errmsg"] = "请输入正确的手机号"
		this.Data["json"] = &resp

		this.ServeJSON()

		return
	}
	user := models.User{}
	user.Name = params["mobile"].(string)
	user.Mobile = params["mobile"].(string)
	user.Password_hash = params["password"].(string)
	o := orm.NewOrm()
	_, err = o.Insert(&user)
	if err != nil {
		resp["errmsg"] = err
		this.Data["json"] = &resp
		this.ServeJSON()
		return
	} else {
		resp["errno"] = 0
		resp["errmsg"] = "成功"
	}
	this.SetSession("id", user.Id)
	this.SetSession("name", user.Name)
	this.Data["json"] = &resp
	this.ServeJSON()
}

func (this *UserController) Get() {
	resp := make(map[string]interface{})
	id := this.GetSession("id").(int)
	user := models.User{Id: id}
	o := orm.NewOrm()
	err := o.Read(&user)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		this.Data["json"] = resp
		this.ServeJSON()
	}
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	data := make(map[string]interface{})
	data["name"] = user.Name
	data["avatar"] = user.Avatar_url
	data["mobile"] = user.Mobile
	resp["data"] = data
	this.Data["json"] = resp
	fmt.Println(resp)
	this.ServeJSON()
}

//修改用户名
func (this *UserController) ModifyName() {
	params := make(map[string]interface{})
	resp := make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		this.Data["json"] = &resp
		this.ServeJSON()
		return
	}
	name := params["name"].(string)
	if name == "" {
		resp["errno"] = 4001
		resp["errmsg"] = "昵称不能为空"
		this.Data["json"] = &resp
		this.ServeJSON()
		return
	}
	id := this.GetSession("id").(int)
	o := orm.NewOrm()
	user := models.User{Id: id}
	err = o.Begin()
	if err != nil {
		logs.Error("start the transaction failed")
		return
	}
	o.Read(&user)
	user.Name = name
	_, err = o.Update(&user, "Name")
	if err != nil {
		o.Rollback()
		resp["errno"] = 4001
		resp["errmsg"] = err
		this.Data["json"] = &resp
		this.ServeJSON()
		return

	}
	o.Commit()
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *UserController) CheckOrders() {
	role := this.GetString("role")
	resp := make(map[string]interface{})
	id := this.GetSession("id").(int)
	user := models.User{Id: id}
	data := make(map[string]interface{})
	var order []map[string]interface{}
	o := orm.NewOrm()
	o.Read(&user)

	_, err := o.LoadRelated(&user, "Orders")
	if err != nil {
		fmt.Println(err)
	}
	o.LoadRelated(&user, "Houses")
	if role == "custom" {

		orders := user.Orders

		for _, v := range orders {
			_, err = o.LoadRelated(v, "House")
			if err != nil {
				fmt.Println(err)
			}
			order = append(order, v.To_order_info().(map[string]interface{}))
		}

	} else if role == "landlord" {

		for _, v := range user.Houses {
			o.LoadRelated(v, "Orders")
			for _, v2 := range v.Orders {
				o.LoadRelated(v2, "House")
				order = append(order, v2.To_order_info().(map[string]interface{}))
			}
		}
	}
	resp["errno"] = 0
	resp["errmsg"] = "成功"

	data["orders"] = order
	resp["data"] = data
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *UserController) retData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *UserController) CheckHouses() {
	resp := make(map[string]interface{})
	defer this.retData(resp)
	id := this.GetSession("id").(int)
	o := orm.NewOrm()
	user := models.User{Id: id}
	err := o.Read(&user)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return
	}
	var houses []map[string]interface{}
	qs := o.QueryTable("House")
	qs.Filter("user_id", user.Id).RelatedSel().All(&user.Houses)
	for _, v := range user.Houses {

		houses = append(houses, v.To_house_info().(map[string]interface{}))
	}
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	house := make(map[string]interface{})
	house["houses"] = houses
	resp["data"] = house

}
func (this *UserController) UploadAvatar() {
	resp := make(map[string]interface{})
	defer this.retData(resp)
	f, h, err := this.GetFile("avatar")
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		fmt.Println(err, 1)
		return
	}
	defer f.Close()
	o := orm.NewOrm()
	o.Begin()
	user := models.User{Id: this.GetSession("id").(int)}
	o.Read(&user)
	savePath := "static/upload/avatar/" + user.Name + h.Filename
	err = this.SaveToFile("avatar", savePath) // 保存位置在 static/upload
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		fmt.Println(err, 2)
		o.Rollback()
		return
	}
	user.Avatar_url = savePath
	_, err = o.Update(&user, "Avatar_url")
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		fmt.Println(err, 3)
		o.Rollback()
		return
	}
	o.Commit()
	resp["errno"] = 0
	data := make(map[string]interface{})
	data["avatar_url"] = utils.AddDomain2Url(user.Avatar_url)
	resp["data"] = data

}
