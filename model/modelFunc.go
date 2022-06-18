package model

import (
	"github.com/gomodule/redigo/redis"
	"encoding/json"
	"fmt"
)

//获取所有地域信息
func GetArea()([]Area,error){
	//连接数据库
	var areas []Area

	//从缓存中获取数据  从redis中获取数据
	//conn,err := redis.Dial("tcp","192.168.137.81:6379")
	conn := GlobalRedis.Get()
	//关闭,释放资源
	areaByte,_:= redis.Bytes(conn.Do("get","areaData"))
	if len(areaByte) == 0{
		//从mysql中获取数据
		if err := GlobalDB.Find(&areas).Error;err != nil {
			return areas,err
		}

		//序列化数据,存入redis中
		//把数据存入redis中
		areaJson,err := json.Marshal(areas)
		if err != nil {
			return nil,err
		}
		_,err = conn.Do("set","areaData",areaJson)
		fmt.Println(err)

		fmt.Println("从mysql中获取数据")
	}else {
		json.Unmarshal(areaByte,&areas)
		fmt.Println("从redis中获取数据")
	}
	return areas,nil
}
