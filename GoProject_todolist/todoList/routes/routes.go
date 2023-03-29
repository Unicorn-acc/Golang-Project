package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"todoList.com/todoList/api"
	"todoList.com/todoList/middleware"
)

// 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default() //生成了一个WSGI应用程序实例
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	// 编写分组路由
	v1 := r.Group("/api/v1")
	{
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		authed := v1.Group("/")
		authed.Use(middleware.JWT()) // 进行权限验证，验证通过才有权限进行访问下面的路由
		{
			authed.POST("task", api.CreateTask)
			authed.GET("task/:id", api.GetTaskById)
			authed.GET("tasks", api.GetAllTaskById)
			authed.PUT("task/:id", api.UpdateTaskById)
			authed.POST("search", api.SearchTask)
			authed.DELETE("task/:id", api.DeleteTask)
		}
	}
	return r
}
