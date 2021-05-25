package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	userOrder "ihomebj5q/service/userOrder/proto/userOrder"
)

type UserOrder struct{}

func (e *UserOrder) Handle(ctx context.Context, msg *userOrder.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *userOrder.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
