package subscriber

import (
	"context"

	"github.com/micro/go-micro/util/log"

	house "ihome/service/house/proto/house"
)

type House struct{}

func (e *House) Handle(ctx context.Context, msg *house.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *house.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
