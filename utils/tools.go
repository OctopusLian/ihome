package utils

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry/consul"
)

func GetMicroClient() client.Client {
	consulReg := consul.NewRegistry()
	microService := micro.NewService(
		micro.Registry(consulReg),
	)
	return microService.Client()
}
