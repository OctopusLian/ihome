package main

import (
	"ihome/service/userOrder/handler"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	"ihome/service/userOrder/model"
	userOrder "ihome/service/userOrder/proto/userOrder"

	"github.com/micro/go-micro/registry/consul"
)

func main() {
	// New Service
	consulReg := consul.NewRegistry()

	service := micro.NewService(
		micro.Name("go.micro.srv.userOrder"),
		micro.Version("latest"),
		micro.Registry(consulReg),
		micro.Address(":9986"),
	)

	// Initialise service
	service.Init()
	model.InitDb()

	// Register Handler
	userOrder.RegisterUserOrderHandler(service.Server(), new(handler.UserOrder))


	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
