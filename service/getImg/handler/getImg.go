package handler

import (
	"context"

	"encoding/json"
	"github.com/afocus/captcha"
	"ihome/service/getImg/model"
	getImg "ihome/service/getImg/proto/getImg"
	"ihome/service/getImg/utils"
	"image/color"
)

type GetImg struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetImg) MicroGetImg(ctx context.Context, req *getImg.Request, rsp *getImg.Response) error {
	//生成验证码图片,存储图片验证码,返回图片数据
	cap := captcha.New()

	//设置字符集
	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}
	//设置验证码图片大小
	cap.SetSize(128, 64)
	//设置混淆程度
	cap.SetDisturbance(captcha.NORMAL)
	//设置字体颜色
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255}, color.RGBA{255, 0, 0, 255})
	//设置背景色  background
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	//生成验证码图片
	//rand.Seed(time.Now().UnixNano())
	img, rnd := cap.Create(4, captcha.NUM)

	//存储验证码   redis
	err := model.SaveImgRnd(req.Uuid, rnd)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return err
	}

	//传递图片信息给调用者
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	//json序列化
	imgJson, err := json.Marshal(img)
	rsp.Data = imgJson
	return nil
}
