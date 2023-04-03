package service

import (
	"example.com/unicorn-acc/dao"
	"example.com/unicorn-acc/model"
	"example.com/unicorn-acc/pkg/e"
	"example.com/unicorn-acc/serializer"
	"golang.org/x/net/context"
	"strconv"
)

type ListProductImgService struct {
}

func (service *ListProductImgService) List(ctx context.Context, id string) serializer.Response {
	var productImg []*model.ProductImg
	code := e.SUCCESS
	pId, _ := strconv.Atoi(id)
	productImgDao := dao.NewProductImgDao(ctx)
	productImg, err := productImgDao.GetproductImgById(uint(pId))
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.BuildListResponse(serializer.BuildProductImgs(productImg), uint(len(productImg)))

}
