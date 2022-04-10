package handler

import (
	"context"

	"ihome/service/getArea/model"
	getArea "ihome/service/getArea/proto/getArea"
	"ihome/service/getArea/utils"
)

type GetArea struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetArea) MicroGetArea(ctx context.Context, req *getArea.Request, rsp *getArea.Response) error {
	//获取数据并返回给调用者
	areas, err := model.GetArea()
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return err
	}
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	for _, v := range areas {
		var areaInfo getArea.AreaInfo
		areaInfo.Aid = int32(v.Id)
		areaInfo.Aname = v.Name

		rsp.Data = append(rsp.Data, &areaInfo)
	}

	return nil
}
