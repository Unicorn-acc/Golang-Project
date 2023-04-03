package serializer

import (
	"example.com/unicorn-acc/model"
	"example.com/unicorn-acc/pkg/utils"
)

type Money struct {
	UserID    uint   `json:"user_id" form:"user_id"`
	UserName  string `json:"user_name" form:"user_name"`
	UserMoney string `json:"user_money" form:"user_money"`
}

func BuildMoney(item *model.User, key string) Money {
	utils.Encrypt.SetKey(key)
	return Money{
		UserID:    item.ID,
		UserName:  item.UserName,
		UserMoney: utils.Encrypt.AesDecoding(item.Money),
	}
}
