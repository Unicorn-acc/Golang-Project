package service

import (
	"time"
	"todoList.com/todoList/model"
	"todoList.com/todoList/serializer"
)

type CreateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status" form:"status"` // 0未做 1已做
}

type ShowTaskService struct {
	// Get请求，因此这个服务内容是空的
}

type GetListTaskService struct {
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
}

type UpdateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status" form:"status"` // 0未做 1已做
}

type SearchTaskService struct {
	Info     string `json:"info" form:"info"`
	PageNum  int    `json:"page_num" form:"page_num"`
	PageSize int    `json:"page_size" form:"page_size"`
}

type DeleteTaskService struct {
}

// 新增一条备忘录
func (service *CreateTaskService) Create(id uint) serializer.Response {
	var user model.User
	model.DB.First(&user, id)
	task := model.Task{
		User:      user,
		Uid:       user.ID,
		Title:     service.Title,
		Content:   service.Content,
		Status:    0,
		StartTime: time.Now().Unix(),
	}
	code := 200
	err := model.DB.Create(&task).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "创建备忘录失败",
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task),
		Msg:    "创建备忘录成功",
	}
}

// 展示用户的一条备忘录
func (service *ShowTaskService) GetById(uid uint, tid string) serializer.Response {
	// 1.查询到备忘录信息
	var task model.Task
	code := 200
	err := model.DB.First(&task, tid).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "查询备忘录失败",
		}
	}
	// 如果查询的备忘录不是当前用户的，返回查询失败
	if task.Uid != uid {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "未查询到当前用户的备忘录信息",
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task),
		Msg:    "查询成功",
	}
}

// 展示用户的所有备忘录
func (service *GetListTaskService) GetAllById(uid uint) serializer.Response {
	var tasks []model.Task
	var count int64 = 0
	if service.PageSize == 0 {
		service.PageSize = 15 // 如果传来的页面大小是0，默认为15
	}
	// 多表查询
	// 先找到是哪一个user，聚合函数查看一共多少条，然后进行分页操作
	model.DB.Model(&model.Task{}).Preload("User").Where("uid = ?", uid).Count(&count).
		Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&tasks)
	// 返回带数量的列表响应
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(count))
}

func (service *UpdateTaskService) UpdateById(id string) serializer.Response {
	var task model.Task
	model.DB.First(&task, id)
	task.Content = service.Content
	task.Title = service.Title
	task.Status = service.Status
	err := model.DB.Save(&task).Error
	if err != nil {
		return serializer.Response{
			Status: 400,
			Msg:    "更新失败",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "更新成功",
	}

}

func (service *SearchTaskService) Search(id uint) serializer.Response {
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	var tasks []model.Task
	var count int64
	//1. 先预加载用户
	model.DB.Model(&model.Task{}).Preload("User").Where("uid = ?", id).
		// 在用户的基础上搜索内容或者标题
		Where("title Like ? OR content like ?", "%"+service.Info+"%", "%"+service.Info+"%").
		// 计数 + 分页后赋值到tasks里
		Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&tasks)
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), uint(count))
}

func (service *DeleteTaskService) Delete(id string) serializer.Response {
	var task model.Task
	err := model.DB.Delete(&task, id).Error
	if err != nil {
		return serializer.Response{
			Status: 400,
			Msg:    "删除失败",
		}
	}
	return serializer.Response{
		Status: 200,
		Msg:    "删除成功",
	}
}
