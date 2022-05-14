/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 16:24:08
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-15 00:04:25
 */
package utils

import (
	"github.com/astaxie/beego/config"
)

var (
	G_image_addr   string
	G_server_name  string //项目名称
	G_server_addr  string //服务器ip地址
	G_server_port  string
	G_redis_addr   string
	G_redis_port   string
	G_redis_dbnum  string
	G_mysql_addr   string
	G_mysql_port   string
	G_mysql_dbname string
	G_fastdfs_addr string //fastdfs ip
	G_fastdfs_port string //fastdfs 端口
)

func InitConfig() {
	//从配置文件读取配置信息
	//env := os.Getenv("ENV_CLUSTER")
	env := "beta"
	ConfPath := ""
	if env == "online" {
		ConfPath = "./conf/online.conf"
	} else if env == "beta" {
		ConfPath = "./conf/app.conf"
	} else {
		ConfPath = "./conf/dev.conf"
	}
	appconf, _ := config.NewConfig("ini", ConfPath)
	// if err != nil {
	// 	beego.Debug(err)
	// 	return
	// }
	G_image_addr = appconf.String("imageaddr")
	//G_server_addr = appconf.String("httpaddr")
	G_server_port = appconf.String("httpport")
	G_redis_addr = appconf.String("redisaddr")
	G_redis_port = appconf.String("redisport")
	G_redis_dbnum = appconf.String("redisdbnum")
	G_mysql_addr = appconf.String("mysqladdr")
	G_mysql_port = appconf.String("mysqlport")
	G_mysql_dbname = appconf.String("mysqldbname")

	return
}

func init() {
	InitConfig()
}
