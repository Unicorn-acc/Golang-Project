package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
	"todoList.com/todoList/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 400 // 返回400， 因为传入的token不对
		} else {
			// 对token进行解析
			claim, err := util.ParseToken(token)
			if err != nil {
				code = 403 // forbid, token解析有误，无权限，是假的
			} else if time.Now().Unix() > claim.ExpiresAt {
				code = 401 // Token无效，已经过期了
			}
		}
		if code != 200 {
			c.JSON(400, gin.H{
				"status": code,
				"msg":    "Token解析错误",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
