package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	register "ihomebj5q/service/register/proto/register"
)

type Register struct{}

func (e *Register) Handle(ctx context.Context, msg *register.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *register.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
