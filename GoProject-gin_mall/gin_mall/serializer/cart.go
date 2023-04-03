package serializer

import (
	"example.com/unicorn-acc/conf"
	"example.com/unicorn-acc/dao"
	"example.com/unicorn-acc/model"
	"golang.org/x/net/context"
)

type Cart struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	ProductID     uint   `json:"product_id"`
	CreateAt      int64  `json:"create_at"`
	Num           uint   `json:"num"`
	MaxNum        uint   `json:"max_num"`
	Check         bool   `json:"check"`
	Name          string `json:"name"`
	ImgPath       string `json:"img_path"`
	DiscountPrice string `json:"discount_price"`
	BossId        uint   `json:"boss_id"`
	BossName      string `json:"boss_name"`
	Desc          string `json:"desc"`
}

func BuildCart(cart *model.Cart, product *model.Product, boss *model.User) Cart {
	return Cart{
		ID:            cart.ID,
		UserID:        cart.UserID,
		ProductID:     cart.ProductID,
		CreateAt:      cart.CreatedAt.Unix(),
		Num:           cart.Num,
		MaxNum:        cart.MaxNum,
		Check:         cart.Check,
		Name:          product.Name,
		ImgPath:       conf.PhotoHost + conf.HttpPort + conf.ProductPhotoPath + product.ImgPath,
		DiscountPrice: product.DiscountPrice,
		BossId:        boss.ID,
		BossName:      boss.UserName,
		Desc:          product.Info,
	}
}

func BuildCarts(items []*model.Cart) (carts []Cart) {
	for _, item1 := range items {
		product, err := dao.NewProductDao(context.Background()).
			GetProductById(item1.ProductID)
		if err != nil {
			continue
		}
		boss, err := dao.NewUserDao(context.Background()).
			GetUserById(item1.BossID)
		if err != nil {
			continue
		}
		cart := BuildCart(item1, product, boss)
		carts = append(carts, cart)
	}
	return carts
}
