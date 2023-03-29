package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"test.cm/gormtest/dao"
	"time"
)

func SaveUser(c *gin.Context) {
	user := &dao.User{
		Username:   "张三",
		Password:   "123456",
		CreateTime: time.Now().UnixMilli(),
	}
	dao.SaveUser(user)
	c.JSON(200, user)
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	idx, _ := strconv.Atoi(id)
	var user dao.User
	user = dao.GetById(idx)
	c.JSON(200, user)
}

func GetAll(c *gin.Context) {
	user := dao.GetAll()
	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	dao.UpdateById(1)
	user := dao.GetById(1)
	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	dao.DeleteUser(1)
	user := dao.GetById(1)
	c.JSON(200, user)
}
