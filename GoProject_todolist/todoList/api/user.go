package api

import (
	"github.com/gin-gonic/gin"
	"todoList.com/todoList/service"
)

// 接收的是gin的上下文：*gin.Context
func UserRegister(ctx *gin.Context) {
	//相当于创建了一个UserRegisterService对象，调用这个对象中的Register方法。
	var userRegisterService service.UserService
	// shouldBind：上下文对UserService对象进行了绑定
	if err := ctx.ShouldBind(&userRegisterService); err == nil {
		// 绑定没有失败的话，调用service中的方法进行注册
		res := userRegisterService.Register()
		ctx.JSON(200, res)
	} else {
		ctx.JSON(400, err)
	}
}

func UserLogin(ctx *gin.Context) {
	var userLogin service.UserService
	if err := ctx.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login()
		ctx.JSON(200, res)
	} else {
		ctx.JSON(400, err)
	}
}
