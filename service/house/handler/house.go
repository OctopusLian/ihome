package handler

import (
	"context"
	house "ihomebj5q/service/house/proto/house"
	"ihomebj5q/service/house/model"
	"ihomebj5q/utils"
	"strconv"
	"github.com/weilaihui/fdfs_client"
)

type House struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *House) PubHouse(ctx context.Context, req *house.Request, rsp *house.Response) error {
	//上传房屋业务  把获取到的房屋数据插入数据库
	houseId,err := model.AddHouse(req)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	var h house.HouseData
	h.HouseId = strconv.Itoa(houseId)
	rsp.Data = &h

	return nil
}

func (e*House) UploadHouseImg(ctx context.Context,req *house.ImgReq,resp *house.ImgResp)error{
	//把图片存储到fastdfs中
	//初始化fdfs的客户端
	fClient ,_:=fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	//上传图片到fdfs
	fdfsResp,err := fClient.UploadByBuffer(req.ImgData,req.FileExt[1:])
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	//把凭证存储到数据库中
	err = model.SaveHouseImg(req.HouseId,fdfsResp.RemoteFileId)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	var img house.ImgData
	img.Url = "http://192.168.137.81:8888/"+fdfsResp.RemoteFileId

	resp.Data = &img

	return nil
}

func (e*House) GetHouseInfo(ctx context.Context,req*house.GetReq,resp*house.GetResp)error{
	//根据用户名获取所有的房屋数据
	houseInfos,err := model.GetUserHouse(req.UserName)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	var getData house.GetData
	getData.Houses = houseInfos

	resp.Data = &getData

	return nil
}

func (e*House) GetHouseDetail(ctx context.Context,req*house.DetailReq,resp*house.DetailResp)error{
	//根据houseId获取所有的返回数据
	respData,err := model.GetHouseDetail(req.HouseId,req.UserName)
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	resp.Data = &respData

	return nil
}

func (e*House) GetIndexHouse(ctx context.Context,req*house.IndexReq,resp*house.GetResp)error{
	//获取房屋信息
	houseResp,err := model.GetIndexHouse()
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	resp.Data = &house.GetData{Houses:houseResp}

	return nil
}

func (e*House) SearchHouse(ctx context.Context ,req *house.SearchReq,resp*house.GetResp)error{
	//根据传入的参数,查询符合条件的房屋信息
	houseResp,err := model.SearchHouse(req.Aid,req.Sd,req.Ed,req.Sk)
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	resp.Data = &house.GetData{Houses:houseResp}
	return nil
}
