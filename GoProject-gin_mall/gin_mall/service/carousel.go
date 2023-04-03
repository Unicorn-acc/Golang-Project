package service

import (
	"context"
	"example.com/unicorn-acc/dao"
	"example.com/unicorn-acc/pkg/e"
	"example.com/unicorn-acc/serializer"
)

type ListCarouselService struct {
}

func (service *ListCarouselService) Show(ctx context.Context) serializer.Response {
	cauouselsdao := dao.NewCarouselDao(ctx)
	carousels, err := cauouselsdao.List()
	if err != nil {
		return serializer.Result(e.ErrorDatabase)
	}
	return serializer.BuildListResponse(serializer.BuildCarousels(carousels), uint(len(carousels)))
}
