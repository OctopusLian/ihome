package subscriber

import (
	"context"

	"github.com/micro/go-micro/util/log"

	getArea "ihome/service/getArea/proto/getArea"
)

type GetArea struct{}

func (e *GetArea) Handle(ctx context.Context, msg *getArea.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *getArea.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
