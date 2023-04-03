package dao

import (
	"example.com/unicorn-acc/model"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type ProductImgDao struct {
	*gorm.DB
}

func NewProductImgDao(ctx context.Context) *ProductImgDao {
	return &ProductImgDao{NewDBClient(ctx)}
}

func NewProductImgDaoByDB(db *gorm.DB) *ProductImgDao {
	return &ProductImgDao{db}
}

// CreateProductImg 创建商品图片
func (dao *ProductImgDao) CreateProductImg(productImg *model.ProductImg) (err error) {
	err = dao.DB.Model(&model.ProductImg{}).Create(&productImg).Error
	return
}

func (dao *ProductImgDao) GetproductImgById(id uint) (products []*model.ProductImg, err error) {
	err = dao.DB.Model(&model.ProductImg{}).Where("product_id = ?", id).Find(&products).Error
	return
}
