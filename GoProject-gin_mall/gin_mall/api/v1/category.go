package v1

import (
	"example.com/unicorn-acc/pkg/utils"
	"example.com/unicorn-acc/service"
	"github.com/gin-gonic/gin"
)

func ListCategories(c *gin.Context) {
	listCategoriesService := service.ListCategoriesService{}
	if err := c.ShouldBind(&listCategoriesService); err == nil {
		res := listCategoriesService.List(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}
