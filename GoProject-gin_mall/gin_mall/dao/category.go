package dao

import (
	"context"
	"example.com/unicorn-acc/model"
	"gorm.io/gorm"
)

type CategoryDao struct {
	*gorm.DB
}

func NewCategoryDao(ctx context.Context) *CategoryDao {
	return &CategoryDao{NewDBClient(ctx)}
}

// ListCategory 分类列表
func (dao *CategoryDao) ListCategory() (category []*model.Category, err error) {
	err = dao.DB.Model(&model.Category{}).Find(&category).Error
	return
}
