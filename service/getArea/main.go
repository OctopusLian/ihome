package main

import (
	"ihome/service/getArea/handler"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	"ihome/service/getArea/model"
	getArea "ihome/service/getArea/proto/getArea"

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
