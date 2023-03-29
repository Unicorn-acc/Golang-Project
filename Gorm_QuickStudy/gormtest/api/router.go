package api

import "github.com/gin-gonic/gin"

// 在这里写路由
func RegisterRouter(r *gin.Engine) {
	r.GET("/save", SaveUser)

	r.GET("/get/:id", GetUserById)

	r.GET("/getAll", GetAll)

	r.GET("/update", UpdateUser)

	r.GET("/delete", DeleteUser)
}
