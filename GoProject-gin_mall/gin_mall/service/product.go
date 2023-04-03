package service

import (
	"context"
	"example.com/unicorn-acc/dao"
	"example.com/unicorn-acc/model"
	"example.com/unicorn-acc/pkg/e"
	"example.com/unicorn-acc/pkg/utils"
	"example.com/unicorn-acc/serializer"
	logging "github.com/sirupsen/logrus"
	"mime/multipart"
	"strconv"
	"sync"
)

// ProductService 商品创建、更新、列表 使用同一个服务
type ProductService struct {
	ID             uint   `form:"id" json:"id"`
	Name           string `form:"name" json:"name"`
	CategoryID     int    `form:"category_id" json:"category_id"`
	Title          string `form:"title" json:"title" `
	Info           string `form:"info" json:"info" `
	ImgPath        string `form:"img_path" json:"img_path"`
	Price          string `form:"price" json:"price"`
	DiscountPrice  string `form:"discount_price" json:"discount_price"`
	OnSale         bool   `form:"on_sale" json:"on_sale"`
	Num            int    `form:"num" json:"num"`
	model.BasePage        // 这个是model上的分页功能
}

func (service *ProductService) Create(ctx context.Context, uId uint, files []*multipart.FileHeader) serializer.Response {
	var boss *model.User
	var err error
	code := e.SUCCESS
	// 1. 获取当前创建商品的用户
	userDao := dao.NewUserDao(ctx)
	boss, _ = userDao.GetUserById(uId)
	// 2. 以第一张作为封面图，获得图片的路径
	tmp, _ := files[0].Open()
	path, err := UploadProductToLocalStatic(tmp, uId, service.Name)
	if err != nil {
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  path,
		}
	}
	// 3. 将商品添加进数据库
	product := &model.Product{
		Name:          service.Name,
		CategoryID:    uint(service.CategoryID),
		Title:         service.Title,
		Info:          service.Info,
		ImgPath:       path,
		Price:         service.Price,
		DiscountPrice: service.DiscountPrice,
		Num:           service.Num,
		OnSale:        true,
		BossID:        uId,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}
	productDao := dao.NewProductDao(ctx)
	err = productDao.CreateProduct(product)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// 4. todo 并发保存商品的图片
	wg := new(sync.WaitGroup)
	wg.Add(len(files))
	for index, file := range files {
		num := strconv.Itoa(index)
		productImgDao := dao.NewProductImgDaoByDB(productDao.DB)
		// 4.1 上传商品图片，获得路径
		tmp, _ = file.Open()
		path, err = UploadProductToLocalStatic(tmp, uId, service.Name+num)
		if err != nil {
			code = e.ErrorUploadFile
			return serializer.Response{
				Status: code,
				Data:   e.GetMsg(code),
				Error:  path,
			}
		}
		// 4.2 上传商品图片
		productImg := &model.ProductImg{
			ProductID: product.ID,
			ImgPath:   path,
		}
		err = productImgDao.CreateProductImg(productImg)
		if err != nil {
			code = e.ERROR
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		wg.Done()
	}
	wg.Wait()
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildProduct(product),
		Msg:    e.GetMsg(code),
	}
}

func (service *ProductService) List(ctx context.Context) serializer.Response {
	var products []*model.Product
	var total int64
	code := e.SUCCESS
	// 1. 配置查询参数
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	// condition ： 查询的条件
	condition := make(map[string]interface{})
	if service.CategoryID != 0 {
		condition["category_id"] = service.CategoryID // 某类商品还是全部商品
	}
	// 3.根据查询条件查询商品数量
	productDao := dao.NewProductDao(ctx)
	total, err := productDao.CountProductByCondition(condition)
	if err != nil {
		code = e.ErrorDatabase
		utils.LogrusObj.Infoln(err)
		return serializer.Result(code)
	}
	// 4. 并发获取
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		productDao = dao.NewProductDaoByDB(productDao.DB)
		products, _ = productDao.ListProductByCondition(condition, service.BasePage)
		wg.Done()
	}()
	wg.Wait()

	// 5. 进行序列化并返回
	return serializer.BuildListResponse(serializer.BuildProducts(products), uint(total))
}

func (service *ProductService) Search(ctx context.Context) serializer.Response {
	var code int = e.SUCCESS
	if service.PageSize == 0 {
		service.PageSize = 15
	}

	productDao := dao.NewProductDao(ctx)
	products, err := productDao.SearchProduct(service.Info, service.BasePage)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.BuildListResponse(serializer.BuildProducts(products), uint(len(products)))

}

func (service *ProductService) Show(ctx context.Context, id string) serializer.Response {
	code := e.SUCCESS
	pId, _ := strconv.Atoi(id)
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(uint(pId))
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildProduct(product),
		Msg:    e.GetMsg(code),
	}
}
