package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model            // Model会一共ID,CreateAt, UpdateAt, DeleteAt
	UserName       string `gorm:"unique"`
	PasswordDigest string // 存储的是密文，也就是加密后的密码
}

// 设置密码
func (user *User) SetPassword(password string) error {
	// 传入类型是[]byte, 和一个加密难度
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
