package v1

import (
	"example.com/unicorn-acc/pkg/utils"
	"example.com/unicorn-acc/service"
	"github.com/gin-gonic/gin"
)

func OrderPay(c *gin.Context) {
	orderPayService := service.OrderPay{}
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&orderPayService); err == nil {
		res := orderPayService.PayDown(c.Request.Context(), claim.ID)
		c.JSON(200, res)
	} else {
		utils.LogrusObj.Infoln(err)
		c.JSON(400, ErrorResponse(err))
	}
}
