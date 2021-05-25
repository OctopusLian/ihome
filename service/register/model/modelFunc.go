package model

import "github.com/gomodule/redigo/redis"

var RedisPool redis.Pool

//连接池


//获取图片验证码
func GetImgCode(uuid string)(string,error){
	//获取redis链接
	conn := RedisPool.Get()
	//获取数据
	return redis.String(conn.Do("get",uuid))
}

//存短信验证码
func SaveSmsCode(phone,vcode string)error{
	//获取redis链接
	conn := RedisPool.Get()
	//存储验证码
	_,err := conn.Do("setex",phone+"_code",60 * 5,vcode)
	return  err
}

//存储用户名和密码  mysql
func SaveUser(mobile,password_hash string)error{
	//链接数据库  gorm插入数据
	var user User
	user.Mobile = mobile
	user.Password_hash = password_hash
	user.Name = mobile

	return GlobalDB.Create(&user).Error
}

//校验短信验证码是否正确
func GetSmsCode(phone string)(string,error){
	//获取redis链接
	conn := RedisPool.Get()
	//获取数据
	return redis.String(conn.Do("get",phone+"_code"))
}

//校验登录信息
func CheckUser(mobile,pwd_hash string)(User,error){
	//连接数据库

	var user User
	err := GlobalDB.Where("mobile = ?",mobile).Where("password_hash = ?",pwd_hash).Find(&user).Error

	return user,err
}

