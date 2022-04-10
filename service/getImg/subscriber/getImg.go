/*
 * @Description:
 * @Author: neozhang
 * @Date: 2021-10-30 12:53:51
 * @LastEditors: neozhang
 * @LastEditTime: 2022-04-10 21:24:58
 */
package subscriber

import (
	"context"

	"github.com/micro/go-micro/util/log"

	getImg "ihome/service/getImg/proto/getImg"
)

type GetImg struct{}

func (e *GetImg) Handle(ctx context.Context, msg *getImg.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *getImg.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
