package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"path"

	getImg "github.com/OctopusLian/ihome/proto/getImg/proto/getImg"

	houseMicro "github.com/OctopusLian/ihome/proto/house/proto/house"
	orderMicro "github.com/OctopusLian/ihome/proto/order/proto/userOrder"
	register "github.com/OctopusLian/ihome/proto/register/proto/register"
	user "github.com/OctopusLian/ihome/proto/user/proto/user"
	"github.com/OctopusLian/ihome/utils"

	getArea "github.com/OctopusLian/ihome/proto/getArea/proto/getArea"

	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"
)

//获取所有地区信息
func GetArea(ctx *gin.Context) {
	/*resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK,resp)
	//获取数据库数据
	areas,err := model.GetArea()
	if err != nil {
		fmt.Println("获取所有地域信息错误")
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}
	//把数据返回给前端
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = areas*/

	//调用远程服务,获取所有地域信息
	//初始化客户端
	//从consul中获取服务
	consulRegistry := consul.NewRegistry()
	micService := micro.NewService(
		micro.Registry(consulRegistry),
	)

	microClient := getArea.NewGetAreaService("go.micro.srv.getArea", micService.Client())
	//调用远程服务
	resp, err := microClient.MicroGetArea(context.TODO(), &getArea.Request{})
	if err != nil {
		fmt.Println(err)
		/*ctx.JSON(http.StatusOK,resp)
		return */
	}

	//把int 的0值  json的特性,如果字段是零值,不对这个字段做序列化

	ctx.JSON(http.StatusOK, resp)
}

//写一个假的session请求返回  真的
func GetSession(ctx *gin.Context) {
	//构造未登录
	resp := make(map[string]interface{})

	//查询session数据,如果查询到了,返回数据
	//初始化session对象
	session := sessions.Default(ctx)

	//获取session数据
	userName := session.Get("userName")
	if userName == nil {
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		//可以是结构体,可以是map
		tempMap := make(map[string]interface{})
		tempMap["name"] = userName.(string)
		resp["data"] = tempMap
	}

	ctx.JSON(http.StatusOK, resp)

}

//获取验证码图片方法
func GetImageCd(ctx *gin.Context) {
	//获取数据
	uuid := ctx.Param("uuid")
	//校验数据
	if uuid == "" {
		fmt.Println("获取数据错误")
		return
	}

	//处理数据

	//调用远程服务
	//初始化客户端
	consulReg := consul.NewRegistry()
	microService := micro.NewService(
		micro.Registry(consulReg),
	)

	microClient := getImg.NewGetImgService("go.micro.srv.getImg", microService.Client())
	//调用远程服务
	resp, err := microClient.MicroGetImg(context.TODO(), &getImg.Request{Uuid: uuid})

	//获取数据
	if err != nil {
		fmt.Println("获取远端数据失败")
		ctx.JSON(http.StatusOK, resp)
		return
	}
	//返回json数据
	//反序列化拿到img数据
	var img captcha.Image
	json.Unmarshal(resp.Data, &img)

	png.Encode(ctx.Writer, img)
}

func GetSmscd(ctx *gin.Context) {
	//获取数据
	mobile := ctx.Param("mobile")
	//获取输入的图片验证码
	text := ctx.Query("text")
	//获取验证码图片的uuid
	uuid := ctx.Query("id")

	//校验数据
	if mobile == "" || text == "" || uuid == "" {
		fmt.Println("传入数据不完整")
		return
	}

	//处理数据  放在服务端处理
	//初始化客户端
	microClient := register.NewRegisterService("go.micro.srv.register", utils.GetMicroClient())
	//调用远程客户端
	resp, err := microClient.SmsCode(context.TODO(), &register.Request{
		Uuid:   uuid,
		Text:   text,
		Mobile: mobile,
	})

	if err != nil {
		fmt.Println("调用远程服务错误", err)
		/*ctx.JSON(http.StatusOK,resp)
		return*/
	}

	ctx.JSON(http.StatusOK, resp)
}

//注册方法
type RegStu struct {
	Mobile   string `json:"mobile"`
	PassWord string `json:"password"`
	SmsCode  string `json:"sms_code"`
}

//注册业务
func PostRet(ctx *gin.Context) {
	//获取数据
	/*mobile := ctx.PostForm("mobile")
	pwd := ctx.PostForm("password")
	smsCode := ctx.PostForm("sms_code")*/
	var reg RegStu
	err := ctx.Bind(&reg)

	//校验数据
	if err != nil {
		fmt.Println("获取前段传递数据失败")
		return
	}
	//处理数据  放在服务端处理
	//初始化客户端
	microClient := register.NewRegisterService("go.micro.srv.register", utils.GetMicroClient())
	//调用远程服务
	resp, err := microClient.Register(context.TODO(), &register.RegRequest{
		Mobile:   reg.Mobile,
		Password: reg.PassWord,
		SmsCode:  reg.SmsCode,
	})

	if err != nil {
		fmt.Println("调用远程服务错误", err)
	}

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

type LogStu struct {
	Mobile   string `json:"mobile"`
	PassWord string `json:"password"`
}

//登录业务
func PostLogin(ctx *gin.Context) {
	//获取数据
	var log LogStu
	err := ctx.Bind(&log)
	//校验数据
	if err != nil {
		fmt.Println("获取数据失败")
		return
	}
	//处理数据   把业务放在为服务中
	//初始化客户端
	microClient := register.NewRegisterService("go.micro.srv.register", utils.GetMicroClient())

	//调用远程服务
	resp, err := microClient.Login(context.TODO(), &register.RegRequest{Mobile: log.Mobile, Password: log.PassWord})
	defer ctx.JSON(http.StatusOK, resp)
	if err != nil {
		fmt.Println("调用login服务错误", err)
		return
	}

	//返回数据  存储session  并返回数据给web端
	session := sessions.Default(ctx)
	session.Set("userName", resp.Name)
	session.Save()
}

//退出登录
func DeleteSession(ctx *gin.Context) {
	//删除session
	session := sessions.Default(ctx)

	//删除session
	session.Delete("userName")
	err := session.Save()

	fmt.Println("控制器函数执行....")

	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)
	if err != nil {
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
}

//获取用户信息
func GetUserInfo(ctx *gin.Context) {
	//获取session数据
	session := sessions.Default(ctx)
	userName := session.Get("userName")

	//调用远程服务
	microClient := user.NewUserService("go.micro.srv.user", utils.GetMicroClient())
	//调用远程服务
	resp, err := microClient.MicroGetUser(context.TODO(), &user.Request{Name: userName.(string)})
	if err != nil {
		fmt.Println("调用远程user服务错误", err)
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	}

	ctx.JSON(http.StatusOK, resp)
}

type UpdateStu struct {
	Name string `json:"name"`
}

//更新用户名
func PutUserInfo(ctx *gin.Context) {
	//获取数据
	var nameData UpdateStu
	err := ctx.Bind(&nameData)
	//校验数据
	if err != nil {
		fmt.Println("获取数据错误")
		return
	}

	//从session中获取原来的用户名
	session := sessions.Default(ctx)
	userName := session.Get("userName")
	//处理数据
	microClient := user.NewUserService("go.micro.srv.user", utils.GetMicroClient())
	//调用远程服务
	resp, _ := microClient.UpdateUserName(context.TODO(), &user.UpdateReq{NewName: nameData.Name, OldName: userName.(string)})

	//更新session数据
	if resp.Errno == utils.RECODE_OK {
		//更新成功,session中的用户名也需要更新一下
		session.Set("userName", nameData.Name)
		session.Save()
	}

	ctx.JSON(http.StatusOK, resp)

}

//上传用户头像
func PostAvatar(ctx *gin.Context) {
	//获取数据  获取图片  文件流  文件头  err
	fileHeader, err := ctx.FormFile("avatar")

	//检验数据
	if err != nil {
		fmt.Println("文件上传失败")
		return
	}

	//三种校验 大小,类型,防止重名  fastdfs
	if fileHeader.Size > 50000000 {
		fmt.Println("文件过大,请重新选择")
		return
	}

	fileExt := path.Ext(fileHeader.Filename)
	if fileExt != ".png" && fileExt != ".jpg" {
		fmt.Println("文件类型错误,请重新选择")
		return
	}
	//只读的文件指针
	file, _ := fileHeader.Open()
	buf := make([]byte, fileHeader.Size)
	file.Read(buf)

	/*
		fdfsClient,_ := fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
		//fdfsClient.UploadByFilename()
		fdfsResp,_ := fdfsClient.UploadByBuffer(buf,fileExt[1:])
		fmt.Println("上传文件到fastdfs的组名为",fdfsResp.GroupName," 凭证为",fdfsResp.RemoteFileId)*/

	//获取用户名
	session := sessions.Default(ctx)
	userName := session.Get("userName")

	//处理数据
	//初始化客户端
	microClient := user.NewUserService("go.micro.srv.user", utils.GetMicroClient())
	//调用远程函数
	resp, _ := microClient.UploadAvatar(context.TODO(), &user.UploadReq{
		UserName: userName.(string),
		Avatar:   buf,
		FileExt:  fileExt,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

type AuthStu struct {
	IdCard   string `json:"id_card"`
	RealName string `json:"real_name"`
}

func PutUserAuth(ctx *gin.Context) {
	//获取数据
	var auth AuthStu
	err := ctx.Bind(&auth)
	//校验数据
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}

	session := sessions.Default(ctx)
	userName := session.Get("userName")

	//处理数据  微服务
	microClient := user.NewUserService("go.micro.srv.user", utils.GetMicroClient())
	//调用远程服务
	resp, _ := microClient.AuthUpdate(context.TODO(), &user.AuthReq{
		UserName: userName.(string),
		RealName: auth.RealName,
		IdCard:   auth.IdCard,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

//获取已发布房源信息  假数据
func GetUserHouses(ctx *gin.Context) {

	/*resp := make(map[string]interface{})

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = ""

	ctx.JSON(http.StatusOK,resp)*/

	//获取用户名
	userName := sessions.Default(ctx).Get("userName")

	/*//测试一对多查询
	//有用户名
	var userInfo model.User
	if err := model.GlobalDB.Where("name = ?",userName).Find(&userInfo).Error;err != nil {
		fmt.Println("获取当前用户信息错误",err)
	}
	//房源信息   一对多查询
	var houses []model.House

	model.GlobalDB.Model(&userInfo).Related(&houses)
	fmt.Println("11111111",houses)*/

	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//调用远程服务
	resp, _ := microClient.GetHouseInfo(context.TODO(), &houseMicro.GetReq{UserName: userName.(string)})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

type HouseStu struct {
	Acreage   string   `json:"acreage"`
	Address   string   `json:"address"`
	AreaId    string   `json:"area_id"`
	Beds      string   `json:"beds"`
	Capacity  string   `json:"capacity"`
	Deposit   string   `json:"deposit"`
	Facility  []string `json:"facility"`
	MaxDays   string   `json:"max_days"`
	MinDays   string   `json:"min_days"`
	Price     string   `json:"price"`
	RoomCount string   `json:"room_count"`
	Title     string   `json:"title"`
	Unit      string   `json:"unit"`
}

//发布房源
func PostHouses(ctx *gin.Context) {
	//获取数据   bind数据的时候不带自动转换   c.getInt()
	var house HouseStu
	err := ctx.Bind(&house)

	//校验数据
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}

	//获取用户名
	userName := sessions.Default(ctx).Get("userName")

	//处理数据  服务端处理
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//调用远程服务
	resp, _ := microClient.PubHouse(context.TODO(), &houseMicro.Request{
		Acreage:   house.Acreage,
		Address:   house.Address,
		AreaId:    house.AreaId,
		Beds:      house.Beds,
		Capacity:  house.Capacity,
		Deposit:   house.Deposit,
		Facility:  house.Facility,
		MaxDays:   house.MaxDays,
		MinDays:   house.MinDays,
		Price:     house.Price,
		RoomCount: house.RoomCount,
		Title:     house.Title,
		Unit:      house.Unit,
		UserName:  userName.(string),
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

//上传房屋图片
func PostHousesImage(ctx *gin.Context) {
	//获取数据
	houseId := ctx.Param("id")
	fileHeader, err := ctx.FormFile("house_image")
	//校验数据
	if houseId == "" || err != nil {
		fmt.Println("传入数据不完整", err)
		return
	}

	//三种校验 大小,类型,防止重名  fastdfs
	if fileHeader.Size > 50000000 {
		fmt.Println("文件过大,请重新选择")
		return
	}

	fileExt := path.Ext(fileHeader.Filename)
	if fileExt != ".png" && fileExt != ".jpg" {
		fmt.Println("文件类型错误,请重新选择")
		return
	}

	//获取文件字节切片
	file, _ := fileHeader.Open()
	buf := make([]byte, fileHeader.Size)
	file.Read(buf)

	//处理数据  服务中实现
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//调用服务
	resp, _ := microClient.UploadHouseImg(context.TODO(), &houseMicro.ImgReq{
		HouseId: houseId,
		ImgData: buf,
		FileExt: fileExt,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

//获取房屋详情
func GetHouseInfo(ctx *gin.Context) {
	//获取数据
	houseId := ctx.Param("id")
	//校验数据
	if houseId == "" {
		fmt.Println("获取数据错误")
		return
	}
	userName := sessions.Default(ctx).Get("userName")
	//处理数据
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//调用远程服务
	resp, _ := microClient.GetHouseDetail(context.TODO(), &houseMicro.DetailReq{
		HouseId:  houseId,
		UserName: userName.(string),
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

func GetIndex(ctx *gin.Context) {
	//处理数据
	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//调用服务
	resp, _ := microClient.GetIndexHouse(context.TODO(), &houseMicro.IndexReq{})

	ctx.JSON(http.StatusOK, resp)
}

//搜索房屋
func GetHouses(ctx *gin.Context) {
	//获取数据
	//areaId
	aid := ctx.Query("aid")
	//start day
	sd := ctx.Query("sd")
	//end day
	ed := ctx.Query("ed")
	//排序方式
	sk := ctx.Query("sk")
	//page  第几页
	//ctx.Query("p")
	//校验数据
	if aid == "" || sd == "" || ed == "" || sk == "" {
		fmt.Println("传入数据不完整")
		return
	}

	//处理数据   服务端  把字符串转换为时间格式,使用函数time.Parse()  第一个参数是转换模板,需要转换的二字符串,两者格式一致
	/*sdTime ,_:=time.Parse("2006-01-02 15:04:05",sd+" 00:00:00")
	edTime,_ := time.Parse("2006-01-02",ed)*/

	/*sdTime,_ :=time.Parse("2006-01-02",sd)
	edTime,_ := time.Parse("2006-01-02",ed)
	d := edTime.Sub(sdTime)
	fmt.Println(d.Hours())*/

	microClient := houseMicro.NewHouseService("go.micro.srv.house", utils.GetMicroClient())
	//调用远程服务
	resp, _ := microClient.SearchHouse(context.TODO(), &houseMicro.SearchReq{
		Aid: aid,
		Sd:  sd,
		Ed:  ed,
		Sk:  sk,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)

}

type OrderStu struct {
	EndDate   string `json:"end_date"`
	HouseId   string `json:"house_id"`
	StartDate string `json:"start_date"`
}

//下订单
func PostOrders(ctx *gin.Context) {
	//获取数据
	var order OrderStu
	err := ctx.Bind(&order)

	//校验数据
	if err != nil {
		fmt.Println("获取数据错误", err)
		return
	}
	//获取用户名
	userName := sessions.Default(ctx).Get("userName")

	//处理数据  服务端处理业务
	microClient := orderMicro.NewUserOrderService("go.micro.srv.userOrder", utils.GetMicroClient())
	//调用服务
	resp, _ := microClient.CreateOrder(context.TODO(), &orderMicro.Request{
		StartDate: order.StartDate,
		EndDate:   order.EndDate,
		HouseId:   order.HouseId,
		UserName:  userName.(string),
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

//获取订单信息
func GetUserOrder(ctx *gin.Context) {
	//获取get请求传参
	role := ctx.Query("role")
	//校验数据
	if role == "" {
		fmt.Println("获取数据失败")
		return
	}

	//处理数据  服务端
	microClient := orderMicro.NewUserOrderService("go.micro.srv.userOrder", utils.GetMicroClient())
	//调用远程服务
	resp, _ := microClient.GetOrderInfo(context.TODO(), &orderMicro.GetReq{
		Role:     role,
		UserName: sessions.Default(ctx).Get("userName").(string),
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}

type StatusStu struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

//更新订单状态
func PutOrders(ctx *gin.Context) {
	//获取数据
	id := ctx.Param("id")
	var statusStu StatusStu
	err := ctx.Bind(&statusStu)

	//校验数据
	if err != nil || id == "" {
		fmt.Println("获取数据错误", err)
		return
	}

	//处理数据   更新订单状态
	microClient := orderMicro.NewUserOrderService("go.micro.srv.userOrder", utils.GetMicroClient())
	//调用元和产能服务
	resp, _ := microClient.UpdateStatus(context.TODO(), &orderMicro.UpdateReq{
		Action: statusStu.Action,
		Reason: statusStu.Reason,
		Id:     id,
	})

	//返回数据
	ctx.JSON(http.StatusOK, resp)
}
