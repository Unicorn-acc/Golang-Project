package v1

import (
	"example.com/unicorn-acc/service"
	"github.com/gin-gonic/gin"
)

func ListCarousels(c *gin.Context) {
	var listCarouselService service.ListCarouselService
	if err := c.ShouldBind(&listCarouselService); err == nil {
		res := listCarouselService.Show(c)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
