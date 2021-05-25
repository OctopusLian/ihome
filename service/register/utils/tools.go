package utils

import (
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
)

func GetMicroClient()client.Client{
	consulReg := consul.NewRegistry()
	microService := micro.NewService(
		micro.Registry(consulReg),
	)
	return microService.Client()
}