package routes

import (
	api "example.com/unicorn-acc/api/v1"
	"example.com/unicorn-acc/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())                    // 使用自己写的跨域中间件
	r.StaticFS("/static", http.Dir("./static")) // 配置加载静态资源路径
	v1 := r.Group("api/v1")                     // 新建路由组
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "pong")
		})
		// 用户操作
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)

		// 商品操作
		v1.GET("carousels", api.ListCarousels) // 轮播图
		v1.GET("products", api.ListProducts)
		v1.GET("product/:id", api.ShowProduct)
		v1.POST("products", api.SearchProducts)
		v1.GET("imgs/:id", api.ListProductImg)   // 商品图片
		v1.GET("categories", api.ListCategories) // 商品分类

		authed := v1.Group("/") // 需要登陆保护
		authed.Use(middleware.JWT())
		{
			// 用户操作
			authed.PUT("user", api.UserUpdate)
			authed.POST("avatar", api.UploadAvatar)
			authed.POST("user/sending-email", api.SendEmail)
			authed.POST("user/valid-email", api.ValidEmail)

			// 显示金额
			authed.POST("money", api.ShowMoney)

			// 收藏夹
			authed.GET("favorites", api.ShowFavorites)
			authed.POST("favorites", api.CreateFavorite)
			authed.DELETE("favorites/:id", api.DeleteFavorite)

			// 收获地址操作
			authed.POST("addresses", api.CreateAddress)
			authed.GET("addresses/:id", api.GetAddress)
			authed.GET("addresses", api.ListAddress)
			authed.PUT("addresses/:id", api.UpdateAddress)
			authed.DELETE("addresses/:id", api.DeleteAddress)

			// 购物车
			authed.POST("carts", api.CreateCart)
			authed.GET("carts", api.ShowCarts)
			authed.PUT("carts/:id", api.UpdateCart) // 购物车id
			authed.DELETE("carts/:id", api.DeleteCart)

			// 订单操作
			authed.POST("orders", api.CreateOrder)
			authed.GET("orders", api.ListOrders)
			authed.GET("orders/:id", api.ShowOrder)
			authed.DELETE("orders/:id", api.DeleteOrder)

			// 商品操作
			authed.POST("product", api.CreateProduct)
			//authed.PUT("product/:id", api.UpdateProduct)
			//authed.DELETE("product/:id", api.DeleteProduct)

			// 支付功能
			authed.POST("paydown", api.OrderPay)
		}

	}

	return r
}
