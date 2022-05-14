/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 23:14:50
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 23:42:05
 */
package routers

import (
	"ihome/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetArea") //获取地区
	beego.Router("/api/v1.0/user", &controllers.UserController{})
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")  //上传用户头像
	beego.Router("/api/v1.0/session", &controllers.SessionController{})
	beego.Router("/api/v1.0/user/name", &controllers.UserController{}, "put:ModifyName")  //更新用户名
	beego.Router("/api/v1.0/user/orders", &controllers.UserController{}, "get:CheckOrders")
	beego.Router("/api/v1.0/user/auth", &controllers.AuthController{})
	beego.Router("/api/v1.0/user/houses", &controllers.UserController{}, "get:CheckHouses")
	beego.Router("/api/v1.0/houses", &controllers.HousesController{}, "get:GetHouseByFilter;post:Post")
	beego.Router("/api/v1.0/houses/:id:int", &controllers.HousesController{})
	beego.Router("/api/v1.0/houses/index", &controllers.HousesController{}, "get:Index")
	beego.Router("/api/v1.0/houses/:id:int/images", &controllers.HousesController{}, "post:UploadHouseImage")  //上传房屋照片
	beego.Router("/api/v1.0/orders", &controllers.OrdersControllers{})
	beego.Router("/api/v1.0/orders/:id:int/status", &controllers.OrdersControllers{}, "put:Status") //修改订单状态
}
