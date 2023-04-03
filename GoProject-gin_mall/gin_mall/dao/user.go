package dao

import (
	"context"
	"example.com/unicorn-acc/model"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// ExistOrNotByUserName 根据username判断是否存在该名字
func (dao *UserDao) ExistOrNotByUserName(username string) (user *model.User, exist bool, err error) {
	// 1. find查看一下是否存在该名字
	var count int64
	err = dao.DB.Model(&model.User{}).Where("user_name = ?", username).Find(&user).Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	return user, true, nil
}

// 创建用户
func (dao *UserDao) CreateUser(user *model.User) (err error) {
	return dao.DB.Model(&model.User{}).Create(&user).Error
}

func (dao *UserDao) GetUserById(id uint) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return
}

func (dao *UserDao) UpdateUserById(id uint, user *model.User) error {
	return dao.DB.Model(&model.User{}).Where("id = ?", id).Updates(&user).Error
}
