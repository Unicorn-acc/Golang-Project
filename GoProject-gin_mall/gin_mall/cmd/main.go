package main

import (
	"example.com/unicorn-acc/conf"
	"example.com/unicorn-acc/routes"
)

func main() {
	conf.Init()
	r := routes.NewRouter()
	_ = r.Run(conf.HttpPort)
}
