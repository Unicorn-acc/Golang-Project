package dao

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Goods struct {
	Id         int
	Title      string
	Price      float64
	Stock      int
	Type       int
	CreateTime time.Time
}

func (v Goods) TableName() string {
	return "goods"
}

func (*Goods) BeforeSave(tx *gorm.DB) (err error) {
	log.Println("before create ....")
	return nil
}

func (*Goods) AfterSave(tx *gorm.DB) (err error) {
	log.Println("after create ....")
	return nil
}

func SaveGoods(goods Goods) {
	DB.Create(&goods)
}

func UpdateGoods() {
	goods := Goods{}
	DB.Where("id = ?", 1).Take(&goods)

	goods.Price = 35
	//UPDATE `goods` SET `title`='毛巾',`price`=100.000000,`stock`=100,`type`=0,`create_time  `='2022-11-25 13:03:48' WHERE `id` = 1
	DB.Save(&goods)
}

func FindGood() {
	var goods Goods
	err := DB.Where("id=?", 2).Limit(1).Find(&goods).Error
	fmt.Println(err)
}
