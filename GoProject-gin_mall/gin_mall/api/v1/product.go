package v1

import (
	"example.com/unicorn-acc/pkg/utils"
	"example.com/unicorn-acc/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 创建商品
func CreateProduct(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	var createProductService service.ProductService
	//c.SaveUploadedFile()
	if err := c.ShouldBind(&createProductService); err == nil {
		res := createProductService.Create(c.Request.Context(), claim.ID, files)
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func ListProducts(c *gin.Context) {
	var listPorductService service.ProductService
	if err := c.ShouldBind(&listPorductService); err == nil {
		res := listPorductService.List(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func SearchProducts(c *gin.Context) {
	var searchPorductService service.ProductService
	if err := c.ShouldBind(&searchPorductService); err == nil {
		res := searchPorductService.Search(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}

func ShowProduct(c *gin.Context) {
	var showProductService service.ProductService
	res := showProductService.Show(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, res)
}

func ListProductImg(c *gin.Context) {
	var listPorductImgService service.ListProductImgService
	if err := c.ShouldBind(&listPorductImgService); err == nil {
		res := listPorductImgService.List(c.Request.Context(), c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		utils.LogrusObj.Infoln(err)
	}
}
