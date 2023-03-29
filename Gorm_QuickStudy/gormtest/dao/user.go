package dao

import "log"

// 定义User模型，绑定users表，ORM库操作数据库，需要定义一个struct类型和MYSQL表进行绑定或者叫映射，struct字段和MYSQL表字段一一对应
type User struct {
	ID int64 // 主键
	//通过在字段后面的标签说明，定义golang字段和表字段的关系
	//例如 `gorm:"column:username"` 标签说明含义是: Mysql表的列名（字段名)为username
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	//创建时间，时间戳
	CreateTime int64 `gorm:"column:createtime"`
}

// 告诉框架，我们在操作User结构体时实际上在操作users这个表
func (u User) TableName() string {
	//绑定MYSQL表名为users
	return "users"
}

func SaveUser(user *User) {
	// 数据库的操作，连接数据库
	err := DB.Create(user).Error
	if err != nil {
		log.Println("insert fail : ", err)
	}
}

func GetById(id int) User {
	var user User
	// .First(&user) 查询第一条数据  == limit 0,1 ; 查询结果赋值到user中
	err := DB.Where("id=?", id).First(&user).Error
	if err != nil {
		log.Println("find user error : ", err)
	}
	return user
}

func GetAll() []User {
	var user []User
	// .First(&user) 查询第一条数据  == limit 0,1 ; 查询结果赋值到user中
	err := DB.Find(&user).Error
	if err != nil {
		log.Println("find user error : ", err)
	}
	return user
}

func UpdateById(id int) {
	// .Model(&User{}) 让找到这个user表
	err := DB.Model(&User{}).Where("id=?", id).Update("username", "lisi")
	if err != nil {
		log.Println("update users  fail : ", err)
	}
}

func DeleteUser(id int64) {
	err := DB.Where("id=?", id).Delete(&User{})
	if err != nil {
		log.Println("delete user fail : ", err)
	}
}
