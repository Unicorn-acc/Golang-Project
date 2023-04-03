package v1

import (
	"example.com/unicorn-acc/pkg/utils"
	"example.com/unicorn-acc/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegister 这个相当于Controller层
func UserRegister(c *gin.Context) {
	var userRegisterService service.UserService
	if err := c.ShouldBind(&userRegisterService); err == nil {
		res := userRegisterService.Register(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func UserLogin(c *gin.Context) {
	var userLoginService service.UserService
	if err := c.ShouldBind(&userLoginService); err == nil {
		res := userLoginService.Login(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func UserUpdate(c *gin.Context) {
	var userUpdateService service.UserService
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdateService); err == nil {
		res := userUpdateService.Update(c.Request.Context(), claims.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func UploadAvatar(c *gin.Context) {
	// 从请求表单中获取头像
	file, fileHeader, _ := c.Request.FormFile("file")
	filesize := fileHeader.Size
	uploadAvatarService := service.UserService{}
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&uploadAvatarService); err == nil {
		res := uploadAvatarService.Post(c.Request.Context(), claims.ID, file, filesize)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}

}

func SendEmail(c *gin.Context) {
	var sendEmailService service.SendEmailService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&sendEmailService); err == nil {
		res := sendEmailService.Send(c.Request.Context(), claim.ID)
		c.JSON(200, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func ValidEmail(c *gin.Context) {
	var validEmailService service.ValidEmailService
	if err := c.ShouldBind(&validEmailService); err == nil {
		res := validEmailService.Valid(c.Request.Context(), c.GetHeader("Authorization"))
		c.JSON(200, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

func ShowMoney(c *gin.Context) {
	var showMoneyService service.ShowMoneyService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&showMoneyService); err == nil {
		res := showMoneyService.Show(c.Request.Context(), claim.ID)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
