package dao

import (
	"testing"
	"time"
)

func TestSaveGoods(t *testing.T) {
	goods := Goods{
		Title:      "苹果",
		Price:      5,
		Stock:      100,
		Type:       0,
		CreateTime: time.Now(),
	}
	SaveGoods(goods)
}

func TestUpdateGoods(t *testing.T) {
	UpdateGoods()
}

func TestFindGood(t *testing.T) {
	FindGood()
}
