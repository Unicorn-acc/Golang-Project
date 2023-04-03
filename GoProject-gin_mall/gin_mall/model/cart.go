package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID    uint
	ProductID uint `gorm:"not null"`
	BossID    uint
	Num       uint
	MaxNum    uint
	Check     bool
}
