package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"ihomebj5q/service/getArea/handler"

	getArea "ihomebj5q/service/getArea/proto/getArea"
	"ihomebj5q/service/getArea/model"
	"github.com/micro/go-micro/registry/consul"
)

func main() {

	model.InitDb()
	model.InitRedis()
	// New Service
	consulRegistry := consul.NewRegistry()

	service := micro.NewService(
		micro.Name("go.micro.srv.getArea"),
		micro.Version("latest"),
		micro.Registry(consulRegistry),
	)

	// Initialise service
	service.Init()

	// Register Handler
	getArea.RegisterGetAreaHandler(service.Server(), new(handler.GetArea))


	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
