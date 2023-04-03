package service

import (
	"errors"
	"example.com/unicorn-acc/dao"
	"example.com/unicorn-acc/model"
	"example.com/unicorn-acc/pkg/e"
	"example.com/unicorn-acc/pkg/utils"
	"example.com/unicorn-acc/serializer"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"strconv"
)

type OrderPay struct {
	OrderId   uint    `form:"order_id" json:"order_id"`
	Money     float64 `form:"money" json:"money"`
	OrderNo   string  `form:"orderNo" json:"orderNo"`
	ProductID int     `form:"product_id" json:"product_id"`
	PayTime   string  `form:"payTime" json:"payTime" `
	Sign      string  `form:"sign" json:"sign" `
	BossID    int     `form:"boss_id" json:"boss_id"`
	BossName  string  `form:"boss_name" json:"boss_name"`
	Num       int     `form:"num" json:"num"`
	Key       string  `form:"key" json:"key"`
}

func (service *OrderPay) PayDown(ctx context.Context, uid uint) serializer.Response {
	code := e.SUCCESS

	// 订单需要事务处理
	err := dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		utils.Encrypt.SetKey(service.Key)
		orderDao := dao.NewOrderDaoByDB(tx)
		// 1. 根据订单ID获取订单信息，计算订单金额
		order, err := orderDao.GetOrderById(service.OrderId)
		if err != nil {
			logging.Info(err)
			return err
		}
		money := order.Money
		num := order.Num
		money = money * float64(num)

		// 2. 查看下单用户存不存在
		userDao := dao.NewUserDao(ctx)
		user, err := userDao.GetUserById(uid)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		// 3.对用户的钱进行解密，减去订单 再加密保存
		moneyStr := utils.Encrypt.AesDecoding(user.Money)
		moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)
		if moneyFloat-money < 0.0 { // 用户余额不足，进行回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return errors.New("金币不足")
		}
		// 4.扣除用户的钱
		finMoney := fmt.Sprintf("%f", moneyFloat-money)
		user.Money = utils.Encrypt.AesEncoding(finMoney)
		err = userDao.UpdateUserById(uid, user)
		if err != nil { //  更新用户金额失败
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}
		// 5. 添加商品商家的金额
		boss := new(model.User)
		boss, err = userDao.GetUserById(uint(service.BossID))
		moneyStr = utils.Encrypt.AesDecoding(boss.Money)
		moneyFloat, _ = strconv.ParseFloat(moneyStr, 64)
		finMoney = fmt.Sprintf("%f", moneyFloat+money)
		boss.Money = utils.Encrypt.AesEncoding(finMoney)
		err = userDao.UpdateUserById(uint(service.BossID), boss)
		if err != nil { // 更新boss金额失败，回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}
		// 6. 更新产品数量
		product := new(model.Product)
		productDao := dao.NewProductDao(ctx)
		product, err = productDao.GetProductById(uint(service.ProductID))
		product.Num -= num
		err = productDao.UpdateProduct(uint(service.ProductID), product)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}
		// 7. 更新订单状态
		order.Type = 2
		err = orderDao.UpdateOrderById(service.OrderId, order)
		if err != nil { // 更新订单失败，回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}
		return nil
	})
	if err != nil {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
