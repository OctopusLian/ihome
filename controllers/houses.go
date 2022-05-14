/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 23:12:24
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:36:09
 */
package controllers

import (
	"encoding/json"
	"fmt"
	"ihome/models"
	"ihome/utils"
	"strconv"

	"github.com/beego/beego/orm"
	"github.com/beego/beego/v2/server/web"
)

type HousesController struct {
	web.Controller
}

func (this *HousesController) retData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}
func (this *HousesController) Post() {
	resp := make(map[string]interface{})
	defer this.retData(resp)
	params := make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		this.retData(resp)
		return
	}
	o := orm.NewOrm()
	o.Begin()
	house := models.House{}
	house.Acreage, err = strconv.Atoi(params["acreage"].(string))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Address = params["address"].(string)
	house.Beds = params["beds"].(string)
	house.Capacity, err = strconv.Atoi(params["capacity"].(string))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Deposit, err = strconv.Atoi(params["deposit"].(string))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Max_days, err = strconv.Atoi(params["max_days"].(string))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Min_days, err = strconv.Atoi(params["min_days"].(string))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Price, err = strconv.Atoi(params["price"].(string))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Room_count, err = strconv.Atoi(params["min_days"].(string))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Title = params["title"].(string)
	house.Unit = params["title"].(string)
	areaId, err := strconv.Atoi(params["area_id"].(string))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Area = &models.Area{Id: areaId}
	house.User = &models.User{Id: this.GetSession("id").(int)}
	houseId, err := o.Insert(&house)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	//多对多 m2m插入,将facilities 一起关联插入到表中
	facilities := []*models.Facility{}
	for _, fid := range params["facility"].([]interface{}) {
		id, _ := strconv.Atoi(fid.(string))
		facility := &models.Facility{Id: id}
		facilities = append(facilities, facility)
	}

	// 第一个参数的对象，主键必须有值
	// 第二个参数为对象需要操作的 M2M 字段
	// QueryM2Mer 的 api 将作用于 Id 为 1 的 House
	m2mhouse_facility := o.QueryM2M(&house, "Facilities")

	_, err = m2mhouse_facility.Add(facilities)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	data := make(map[string]interface{})
	data["house_id"] = houseId
	resp["data"] = data
	o.Commit()
}
func (this *HousesController) Get() {
	resp := make(map[string]interface{})
	defer this.retData(resp)
	houseId, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return
	}
	o := orm.NewOrm()
	house := models.House{Id: houseId}

	err = o.Read(&house)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return
	}
	o.LoadRelated(&house, "Area")
	o.LoadRelated(&house, "User")
	o.LoadRelated(&house, "Images")
	o.LoadRelated(&house, "Facilities")
	data := make(map[string]interface{})
	data["house"] = house.To_one_house_desc()
	data["user_id"] = this.GetSession("id")
	resp["data"] = data
	resp["errno"] = 0
	resp["errmsg"] = "成功"

}
func (this *HousesController) UploadHouseImage() {
	resp := make(map[string]interface{})
	defer this.retData(resp)
	houseID, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return
	}
	o := orm.NewOrm()
	o.Begin()
	house := models.House{Id: houseID}
	err = o.Read(&house)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	f, h, err := this.GetFile("house_image")
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		fmt.Println(err, 1)
		o.Rollback()
		return
	}
	defer f.Close()
	savePath := "static/upload/house/" + strconv.Itoa(houseID) + h.Filename
	err = this.SaveToFile("house_image", savePath) // 保存位置在 static/upload
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		fmt.Println(err, 2)
		return
	}
	houseimage := models.HouseImage{Url: savePath, House: &house}
	_, err = o.Insert(&houseimage)
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	house.Index_image_url = savePath
	_, err = o.Update(&house, "Index_image_url")
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		o.Rollback()
		return
	}
	o.Commit()
	resp["errno"] = 0
	data := make(map[string]interface{})
	data["url"] = utils.AddDomain2Url(savePath)
	resp["data"] = data

}

func (this *HousesController) Index() {
	resp := make(map[string]interface{})
	defer this.retData(resp)
	o := orm.NewOrm()
	qs := o.QueryTable("House")
	var houses []models.House
	qs.OrderBy("Order_count").Limit(models.HOME_PAGE_MAX_HOUSES).All(&houses)
	var retData []map[string]interface{}
	for _, v := range houses {
		o.LoadRelated(&v, "Area")
		o.LoadRelated(&v, "User")
		retData = append(retData, v.To_house_info().(map[string]interface{}))
	}
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	resp["data"] = retData

}
func (this *HousesController) GetHouseByFilter() {
	resp := make(map[string]interface{})
	defer this.retData(resp)

	aid := this.GetString("aid")

	sd := this.GetString("sd")
	ed := this.GetString("ed")
	sk := this.GetString("sk")
	p, err := this.GetInt("p")
	if err != nil {
		resp["errno"] = 4001
		resp["errmsg"] = err
		return
	}
	o := orm.NewOrm()
	qs := o.QueryTable("House")
	var houses []models.House
	if aid != "" {
		areaId, err := strconv.Atoi(aid)
		if err != nil {
			resp["errno"] = 4001
			resp["errmsg"] = err
			return
		}
		qs = qs.Filter("Area__Id__exact", areaId)

	}
	if sd != "" {
		qs = qs.Filter("Ctime__gte", sd)
	}
	if ed != "" {
		qs = qs.Filter("Ctime__lte", ed)
	}
	orderby := ""
	switch sk {
	case "new":
		orderby = "-Ctime"
		break
	default:
		orderby = "-Ctime"
		break
	}
	qs = qs.OrderBy(orderby)
	qs = qs.Limit(models.HOUSE_LIST_PAGE_CAPACITY).Offset((p - 1) * models.HOUSE_LIST_PAGE_CAPACITY)
	qs.All(&houses)
	var retHouse []map[string]interface{}
	for _, v := range houses {

		o.LoadRelated(&v, "User")

		retHouse = append(retHouse, v.To_house_info().(map[string]interface{}))
	}
	data := make(map[string]interface{})
	data["houses"] = retHouse
	resp["errno"] = 0
	resp["errmsg"] = "成功"
	resp["data"] = data
}
