package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OctopusLian/ihome/model"
	"github.com/OctopusLian/ihome/utils"
	"github.com/asim/go-micro/plugins/store/redis"
	"github.com/gin-gonic/gin"
)

//路由过滤器
func Filter(ctx *gin.Context) {
	//登录校验
	session := sessions.Default(ctx)
	userName := session.Get("userName")
	resp := make(map[string]interface{})
	if userName == nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		ctx.JSON(http.StatusOK, resp)
		ctx.Abort()
		return
	}

	//计算这个业务耗时
	fmt.Println("next之前打印", time.Now())

	//执行函数
	ctx.Next()

	fmt.Println("next之后打印....")
}

//gin的核心知识点
func MiddleTest() func(*gin.Context) {
	return func(context *gin.Context) {

	}
}

func main() {
	//初始化路由
	router := gin.Default()
	//数据库处理
	model.InitRedis()
	err := model.InitDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	//初始化redis容器，存储session数据
	store, _ := redis.NewStore(20, "tcp", "192.168.137.81:6379", "", []byte("session"))
	//路由模块
	//router.Group()
	//展示静态页面
	//静态路由
	router.Static("/home", "view")
	/*
		router.Use()

		//初始化redis容器
		store,err := redis.NewStore(20,"tcp","192.168.137.81:6379","",[]byte("session"))
		if err != nil {
			fmt.Println("初始化session容器错误")
			return
		}

		store.Options(
			sessions.Options{
				MaxAge:0,
			},
		)

		//路由使用中间件 gin中的session默认是生效时间是一个月
		router.Use(sessions.Sessions("mysession",store))

		//使用路由的时候就可以使用session中间件了
		router.GET("/session", func(context *gin.Context) {
			//初始化session对象
			se := sessions.Default(context)
			//设置session的时候除了set函数之外,必须调用save
			se.Set("test","bj5q")
			se.Save()

			context.Writer.WriteString("设置session")
		})

		//获取session
		router.GET("/getSession", func(context *gin.Context) {
			//初始化session对象
			se := sessions.Default(context)
			//获取session
			result := se.Get("test")
			fmt.Println("得到的session数据为",result.(string))

			context.Writer.WriteString("获取session")
		})

		//测试
		router.GET("/test", func(context *gin.Context) {
			//设置cookie  cookie有两种,一种是有时间效应的,一种是临时cookie
	*/ /*context.SetCookie("myTest","bj5q",0,"","",false,true)
	context.Writer.WriteString("测试cookie")*/ /*


		})*/
	r1 := router.Group("api/v1.0")
	{
		// //路由规范
		// r1.GET("/areas", controller.GetArea)
		// //r1.GET("/session",controller.GetSession)
		// //传参方法,url传值,form表单传值,ajax传值,路径传值
		// r1.GET("/imagecode/:uuid", controller.GetImageCd)
		// r1.GET("/smscode/:mobile", controller.GetSmscd)
		// r1.POST("/users", controller.PostRet)

		// //登录业务   路由过滤器   中间件
		// r1.Use(sessions.Sessions("mysession", store))
		// r1.POST("/sessions", controller.PostLogin)
		// r1.GET("/session", controller.GetSession)
		// //路由过滤器   登录的情况下才能执行一下路由请求
		// r1.Use(Filter)
		// r1.DELETE("/session", controller.DeleteSession)
		// r1.GET("/user", controller.GetUserInfo)
		// r1.PUT("/user/name", controller.PutUserInfo)

		// r1.POST("/user/avatar", controller.PostAvatar)
		// r1.POST("/user/auth", controller.PutUserAuth)
		// r1.GET("/user/auth", controller.GetUserInfo)
		// //获取已发布房源信息
		// r1.GET("/user/houses", controller.GetUserHouses)
		// //发布房源
		// r1.POST("/houses", controller.PostHouses)
		// //添加房源图片
		// r1.POST("/houses/:id/images", controller.PostHousesImage)
		// //展示房屋详情
		// r1.GET("/houses/:id", controller.GetHouseInfo)
		// //展示首页轮播图
		// r1.GET("/house/index", controller.GetIndex)
		// //搜索房屋
		// r1.GET("/houses", controller.GetHouses)
		// //下订单
		// r1.POST("/orders", controller.PostOrders)
		// //获取订单
		// r1.GET("/user/orders", controller.GetUserOrder)
		// //同意/拒绝订单
		// r1.PUT("/orders/:id/status", controller.PutOrders)
	}
	router.Run(":8099")

}
