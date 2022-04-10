package main

import (
	"ihome/service/house/handler"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	"ihome/service/house/model"
	house "ihome/service/house/proto/house"

	"github.com/micro/go-micro/registry/consul"
)

func main() {
	// New Service
	consulReg:=consul.NewRegistry()

	service := micro.NewService(
		micro.Name("go.micro.srv.house"),
		micro.Version("latest"),
		micro.Address(":9985"),
		micro.Registry(consulReg),
	)

	// Initialise service
	service.Init()
	model.InitDb()

	// Register Handler
	house.RegisterHouseHandler(service.Server(), new(handler.House))


	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
