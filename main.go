/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 19:10:13
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 22:59:18
 */
package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	//创建1个新的web服务
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":10086"),
	)
	//服务初始化
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	//使用路由中间件来映射页面
	router := httprouter.New()
	router.NotFound = http.FileServer(http.Dir("html"))
}
