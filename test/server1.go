package main

import "github.com/gin-gonic/gin"

func main(){
	//初始化路由
	router := gin.Default()

	//做路由匹配
	router.GET("/", func(context *gin.Context) {
		context.Writer.WriteString("server1")
	})

	//跑起来
	router.Run(":8081")


}
