package router

import (
	"github.com/gin-gonic/gin"
	"test.cm/gormtest/api"
)

func InitRouter(r *gin.Engine) {
	api.RegisterRouter(r)
}
