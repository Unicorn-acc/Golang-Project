package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.cm/gormtest/router"
)

func main() {
	r := gin.Default()
	router.InitRouter(r)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}