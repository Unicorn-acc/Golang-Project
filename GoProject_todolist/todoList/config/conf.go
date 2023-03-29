package config

import (
	"gopkg.in/ini.v1"
	"log"
	"strings"
	"todoList.com/todoList/model"
)

// 将配置文件内容读取出来
var (
	AppMode  string
	HttpPort string

	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

func Init() {
	// go get gopkg.in/ini.v1 ； 记得要在todolist包下
	// 加载文件
	file, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Println("配置文件读取错误，请检查文件路径", err)
	}
	LoadServer(file)
	LoadMysql(file) // 读取Mysql配置文件
	// 将Mysql配置文件传给model.init， 让它去做数据库的连接
	// 导入gorm：go get github.com/jinzhu/gorm
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	model.DataBase(path)
}

// 从file中加载文件信息
func LoadServer(file *ini.File) {
	/*	config.ini:
		# debug开发模式,release生产模式
		[service]
		AppMode = debug
		HttpPort = :3000
	*/
	// 这句话的意思是选中[service]中的AppMode信息，然后转为string类型赋值给变量
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}
