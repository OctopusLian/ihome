package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"ihomebj5q/service/register/handler"

	register "ihomebj5q/service/register/proto/register"
	"github.com/micro/go-micro/registry/consul"
	"ihomebj5q/service/register/model"
)

func main() {
	//服务发现用consul
	consulReg := consul.NewRegistry()

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.register"),
		micro.Version("latest"),
		micro.Registry(consulReg),
		micro.Address(":9982"),
	)

	// Initialise service
	service.Init()
	model.InitRedis()
	model.InitDb()

	// Register Handler
	register.RegisterRegisterHandler(service.Server(), new(handler.Register))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
