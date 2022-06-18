package handler

import (
	"context"

	"fmt"
	"github.com/weilaihui/fdfs_client"
	"ihome/service/user/model"
	user "ihome/service/user/proto/user"
	"ihome/service/user/utils"
)

type User struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) MicroGetUser(ctx context.Context, req *user.Request, rsp *user.Response) error {
	//根据用户名获取用户信息 在mysql数据库中
	myUser, err := model.GetUserInfo(req.Name)
	if err != nil {
		rsp.Errno = utils.RECODE_USERERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_USERERR)
		return err
	}
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	//获取一个结构体对象
	var userInfo user.UserInfo
	userInfo.UserId = int32(myUser.ID)
	userInfo.Name = myUser.Name
	userInfo.Mobile = myUser.Mobile
	userInfo.RealName = myUser.Real_name
	userInfo.IdCard = myUser.Id_card
	userInfo.AvatarUrl = "http://192.168.137.81:8888/" + myUser.Avatar_url

	rsp.Data = &userInfo

	return nil
}

func (e *User) UpdateUserName(ctx context.Context, req *user.UpdateReq, resp *user.UpdateResp) error {
	//根据传递过来的用户名更新数据中新的用户名
	err := model.UpdateUserName(req.OldName, req.NewName)
	if err != nil {
		fmt.Println("更新失败", err)
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		//micro规定如果有错误,服务端只给客户端返回错误信息,不返回resp,如果没有错误,就返回resp
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var nameData user.NameData
	nameData.Name = req.NewName

	resp.Data = &nameData

	return nil
}

func (e *User) UploadAvatar(ctx context.Context, req *user.UploadReq, resp *user.UploadResp) error {
	//存入到fastdfs中
	fClient, _ := fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	//上传文件到fdfs
	fdfsResp, err := fClient.UploadByBuffer(req.Avatar, req.FileExt[1:])
	if err != nil {
		fmt.Println("上传文件错误", err)
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		return nil
	}

	//把存储凭证写入数据库
	err = model.SaveUserAvatar(req.UserName, fdfsResp.RemoteFileId)
	if err != nil {
		fmt.Println("存储用户头像错误", err)
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	var uploadData user.UploadData
	uploadData.AvatarUrl = "http://192.168.137.81:8888/" + fdfsResp.RemoteFileId
	resp.Data = &uploadData
	return nil
}

func (e *User) AuthUpdate(ctx context.Context, req *user.AuthReq, resp *user.AuthResp) error {
	//调用借口校验realName和idcard是否匹配

	//存储真实姓名和真是身份证号  数据库
	err := model.SaveRealName(req.UserName, req.RealName, req.IdCard)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	return nil
}
