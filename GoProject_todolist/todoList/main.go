package main

import (
	"todoList.com/todoList/config"
	"todoList.com/todoList/routes"
)

func main() {
	config.Init()
	r := routes.NewRouter()
	_ = r.Run(config.HttpPort) // 运行在 配置文件设置的端口上
}
