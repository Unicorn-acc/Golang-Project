package model

import "gorm.io/gorm"

// 任务模型
type Task struct {
	gorm.Model
	User      User   `gorm:"ForeignKey:Uid"`
	Uid       uint   `gorm:"not null"`
	Title     string `gorm:"index;not null"`
	Status    int    `gorm:"default:0"`     // 备忘录状态:0 未完成 ； 1 已完成
	Content   string `gorm:"type:longtext"` // 内容
	StartTime int64  // 备忘录开始时间
	EndTime   int64  `gorm:"default:0"` // 备忘录完成时间
}
