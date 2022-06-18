package handler

import (
	"context"

	"fmt"
	"ihome/service/user/utils"
	"ihome/service/userOrder/model"
	userOrder "ihome/service/userOrder/proto/userOrder"
	"strconv"
)

type UserOrder struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *UserOrder) CreateOrder(ctx context.Context, req *userOrder.Request, rsp *userOrder.Response) error {
	//获取到相关数据,插入到数据库
	orderId, err := model.InsertOrder(req.HouseId, req.StartDate, req.EndDate, req.UserName)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var orderData userOrder.OrderData
	orderData.OrderId = strconv.Itoa(orderId)

	rsp.Data = &orderData

	return nil
}

func (e *UserOrder) GetOrderInfo(ctx context.Context, req *userOrder.GetReq, resp *userOrder.GetResp) error {
	//要根据传入数据获取订单信息   mysql
	respData, err := model.GetOrderInfo(req.UserName, req.Role)
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var getData userOrder.GetData
	getData.Orders = respData
	resp.Data = &getData

	return nil
}

func (e *UserOrder) UpdateStatus(ctx context.Context, req *userOrder.UpdateReq, resp *userOrder.UpdateResp) error {
	//根据传入数据,更新订单状态
	err := model.UpdateStatus(req.Action, req.Id, req.Reason)
	if err != nil {
		fmt.Println("更新订单装填错误", err)
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	return nil
}
