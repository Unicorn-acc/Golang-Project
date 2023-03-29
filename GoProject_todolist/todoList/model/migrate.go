package model

import "log"

func migration() {
	// 自动迁移模式
	err := DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&User{}, &Task{})
	if err != nil {
		log.Println("表迁移失败")
		return
	}
}
