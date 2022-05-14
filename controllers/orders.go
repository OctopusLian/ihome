/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 23:13:38
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:13:43
 */
package controllers

import (
	"encoding/json"
	"ihome/models"
	"strconv"
	"time"

	"github.com/beego/beego/orm"
	"github.com/beego/beego/v2/server/web"
)

type OrdersControllers struct {
	web.Controller
}

func (this *OrdersControllers) retData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *OrdersControllers) Post() {
	resp := make(map[string]interface{})
	defer this.retData(resp)
	params := make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return
	}
	start_date := params["start_date"].(string)
	end_date := params["end_date"].(string)
	house_id := params["house_id"].(string)
	if start_date == "" || end_date == "" || house_id == "" {
		resp["errno"] = 4001
		resp["errmsg"] = "参数缺失"
		return
	}
	realStartDate, err := time.Parse("2006-01-02 15:04:05", start_date+" 00:00:00")
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return
	}
	realEndDate, err := time.Parse("2006-01-02 15:04:05", end_date+" 00:00:00")
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return
	}
	if realEndDate.Before(realStartDate) {
		resp["errno"] = 4001
		resp["errmsg"] = "结束时间在开始时间之前"
		return
	}
	days := realEndDate.Sub(realStartDate).Hours()/24 + 1
	o := orm.NewOrm()
	o.Begin()
	userId := this.GetSession("id").(int)
	user := models.User{Id: userId}
	order := models.OrderHouse{}
	Hid, err := strconv.Atoi(house_id)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house := models.House{Id: Hid}
	err = o.Read(&house)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	order.User = &user
	order.House = &house
	order.Begin_date = realStartDate
	order.End_date = realEndDate
	order.Days = int(days)
	order.House_price = house.Price
	order.Amount = house.Price * int(days)
	order.Status = models.ORDER_STATUS_WAIT_ACCEPT
	_, err = o.Insert(&order)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	resp["errno"] = 0
	o.Commit()
}
func (this *OrdersControllers) Status() {
	resp := make(map[string]interface{})
	defer this.retData(resp)
	params := make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	action := params["action"].(string)
	orderID := this.Ctx.Input.Param(":id")
	orderId, err := strconv.Atoi(orderID)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err

		return
	}
	order := models.OrderHouse{Id: orderId}
	o := orm.NewOrm()
	o.Begin()
	o.Read(&order)
	switch action {
	case "reject":
		reason := params["reason"].(string)

		order.Status = models.ORDER_STATUS_REJECTED
		order.Comment = reason
		break
	case "accept":
		order.Status = models.ORDER_STATUS_WAIT_COMMENT
		break
	case "comment":
		comment := params["comment"].(string)
		order.Status = models.ORDER_STATUS_COMPLETE
		order.Comment = comment
	}
	_, err = o.Update(&order, "Status", "Comment")
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	o.Commit()
	resp["errno"] = 0
	resp["errmsg"] = "成功"
}
