package service

import (
	"context"
	"example.com/unicorn-acc/dao"
	"example.com/unicorn-acc/model"
	"example.com/unicorn-acc/pkg/e"
	"example.com/unicorn-acc/serializer"
	logging "github.com/sirupsen/logrus"
)

type FavoritesService struct {
	ProductId  uint `form:"product_id" json:"product_id"`
	BossId     uint `form:"boss_id" json:"boss_id"`
	FavoriteId uint `form:"favorite_id" json:"favorite_id"`
	PageNum    int  `form:"pageNum"`
	PageSize   int  `form:"pageSize"`
}

func (service *FavoritesService) Create(ctx context.Context, uid uint) serializer.Response {
	code := e.SUCCESS
	favoriteDao := dao.NewFavoritesDao(ctx)
	// 1. 查看一下收藏夹是否存在该商品
	exist, _ := favoriteDao.FavoriteExistOrNot(service.ProductId)
	if exist { // 存在就不用存了
		code = e.ErrorExistFavorite
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 2. 根据用户id获取一下用户信息
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uid)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 3. 根据商品的归属id，获取一下boss信息
	bossDao := dao.NewUserDaoByDB(userDao.DB)
	boss, err := bossDao.GetUserById(service.BossId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 4. 获取一下商品信息
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(service.ProductId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	favorite := &model.Favorite{
		UserID:    uid,
		User:      *user, // 存在外键约束
		ProductID: service.ProductId,
		Product:   *product,
		BossID:    service.BossId,
		Boss:      *boss,
	}
	favoriteDao = dao.NewFavoritesDaoByDB(favoriteDao.DB)
	err = favoriteDao.CreateFavorite(favorite)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}

// Show 商品收藏夹
func (service *FavoritesService) Show(ctx context.Context, uId uint) serializer.Response {
	favoritesDao := dao.NewFavoritesDao(ctx)
	code := e.SUCCESS
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	// 查询商品收藏夹
	favorites, total, err := favoritesDao.ListFavoriteByUserId(uId, service.PageSize, service.PageNum)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.BuildFavorites(ctx, favorites), uint(total))
}

// Delete 删除收藏夹
func (service *FavoritesService) Delete(ctx context.Context) serializer.Response {
	code := e.SUCCESS

	favoriteDao := dao.NewFavoritesDao(ctx)
	err := favoriteDao.DeleteFavoriteById(service.FavoriteId)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data:   e.GetMsg(code),
	}
}
