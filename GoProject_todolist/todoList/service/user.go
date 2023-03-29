package service

import (
	"gorm.io/gorm"
	"todoList.com/todoList/model"
	"todoList.com/todoList/pkg/util"
	"todoList.com/todoList/serializer"
)

// UserService 用户注册服务
type UserService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" example:"FanOne"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=16" example:"FanOne666"`
}

// 返回值是json，所以调用了序列化包下的通用接口
func (service *UserService) Register() serializer.Response {
	var user model.User
	var count int64
	// 根据绑定的用户名查找数据库是否有相同的，取一个放到user中
	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).
		First(&user).Count(&count)
	if count == 1 {
		// 说明已经存在该用户名了,返回注册失败
		return serializer.Response{
			Status: 400,
			Msg:    "数据库已存在该用户名!",
		}
	}
	user.UserName = service.UserName
	// 对用户的密码进行加密
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{ // 用户密码加密错误，返回注册失败
			Status: 400,
			Msg:    "用户名密码加密错误" + err.Error(),
		}
	}
	// 往表中插入用户信息
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{ // 插入数据库错误，返回注册失败
			Status: 400,
			Msg:    "数据库操作错误" + err.Error(),
		}
	}
	// 返回注册成功响应
	return serializer.Response{
		Status: 200,
		Msg:    "用户注册成功!",
	}
}

func (service *UserService) Login() serializer.Response {
	var user model.User
	// 1. 在数据库中查看有没有用户
	if err := model.DB.Where("user_name = ?", service.UserName).
		First(&user).Error; err != nil {
		// 当 First、Last、Take 方法找不到记录时，GORM 会返回 ErrRecordNotFound 错误，
		// 可以通过对比gorm.ErrRecordNotFound进行判断，或者使用Find和Limit的组合进行查询。
		if err == gorm.ErrRecordNotFound {
			return serializer.Response{
				Status: 400,
				Msg:    "用户不存在，请先注册",
			}
		}
		// 如果不是用户不存在，是其他不可抗拒的因素导致的错误
		return serializer.Response{
			Status: 400,
			Msg:    "数据库错误",
		}
	}
	// 用户存在，验证密码是否想用
	if user.CheckPassword(service.Password) == false {
		return serializer.Response{
			Status: 400,
			Msg:    "密码错误",
		}
	}
	// 用户验证通过
	// 发一个token给浏览器，为了其他功能需要身份验证所给请前端存储的
	// 例如：创建一个备忘录，就需要知道这是哪一个用户创建的备忘录
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "Token签名分发错误",
		}
	}
	return serializer.Response{
		Status: 200,
		Data:   serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:    "登录成功",
	}

}
