package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"todoList.com/todoList/pkg/util"
	"todoList.com/todoList/service"
)

func CreateTask(c *gin.Context) {
	var createTaskservice service.CreateTaskService
	// 进行身份的验证
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&createTaskservice); err == nil {
		res := createTaskservice.Create(claim.Id)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
		log.Println("createTask err : ", err)
	}

}

func GetTaskById(c *gin.Context) {
	var GetTaskservice service.ShowTaskService
	// 进行身份的验证
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))

	if err := c.ShouldBind(&GetTaskservice); err == nil {
		// claim.Id是token中取到的用户id， c.Param("id")是前端传递的参数：备忘录id
		res := GetTaskservice.GetById(claim.Id, c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
		log.Println("GetTaskById err : ", err)
	}
}

func GetAllTaskById(c *gin.Context) {
	var getListTaskService service.GetListTaskService
	// 进行身份的验证
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&getListTaskService); err == nil {
		// claim.Id是token中取到的用户id， c.Param("id")是前端传递的参数：备忘录id
		res := getListTaskService.GetAllById(claim.Id)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
		log.Println("GetAllTaskById err : ", err)
	}
}

func UpdateTaskById(c *gin.Context) {
	var updateTaskservice service.UpdateTaskService
	if err := c.ShouldBind(&updateTaskservice); err == nil {
		res := updateTaskservice.UpdateById(c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
		log.Println("UpdateTask err : ", err)
	}
}

func SearchTask(c *gin.Context) {
	var searchTaskservice service.SearchTaskService
	// 进行身份的验证
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))

	if err := c.ShouldBind(&searchTaskservice); err == nil {
		// claim.Id是token中取到的用户id， c.Param("id")是前端传递的参数：备忘录id
		res := searchTaskservice.Search(claim.Id)
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
		log.Println("SearchTask err : ", err)
	}
}

func DeleteTask(c *gin.Context) {
	var deleteTaskservice service.DeleteTaskService

	if err := c.ShouldBind(&deleteTaskservice); err == nil {
		res := deleteTaskservice.Delete(c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
		log.Println("DeleteTask err : ", err)
	}

}
