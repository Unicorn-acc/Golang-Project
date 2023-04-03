package service

import (
	"context"
	"example.com/unicorn-acc/conf"
	"example.com/unicorn-acc/dao"
	"example.com/unicorn-acc/model"
	"example.com/unicorn-acc/pkg/e"
	"example.com/unicorn-acc/pkg/utils"
	"example.com/unicorn-acc/serializer"
	"gopkg.in/mail.v2"
	"mime/multipart"
	"strings"
	"time"
)

// UserService 管理用户服务
type UserService struct {
	NickName string `form:"nick_name" json:"nick_name"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Key      string `form:"key" json:"key"` // 前端进行判断
}

type SendEmailService struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	//OpertionType 1:绑定邮箱 2：解绑邮箱 3：改密码
	OperationType uint `form:"operation_type" json:"operation_type"`
}
type ValidEmailService struct {
}

type ShowMoneyService struct {
	Key string `json:"key" form:"key"`
}

func (service UserService) Register(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS
	// service.Key是密钥
	if service.Key == "" || len(service.Key) != 16 {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密钥长度不足",
		}
	}
	// 进行对称加密  ==> 密文存储 对称加密操作
	utils.Encrypt.SetKey(service.Key)
	// todo 创建一个userdao对象，（将与数据库操作的部分都在DAO层进行处理）
	userDao := dao.NewUserDao(ctx)
	// 1. 判断当前用户名是否被注册了
	_, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//  注册的用户符合条件
	user = &model.User{
		NickName: service.NickName,
		UserName: service.UserName,
		Status:   model.Active,
		Avatar:   "avatar.JPG",                       // 默认头像
		Money:    utils.Encrypt.AesEncoding("10000"), // 初始金额的encoding

	}
	// 2. 加密密码
	if err = user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 3. 创建用户
	err = userDao.CreateUser(user)
	if err != nil {
		code = e.ERROR
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

func (service *UserService) Login(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	// 1.查看用户是否存在
	user, exist, err := userDao.ExistOrNotByUserName(service.UserName)
	if !exist { //如果查询不到，返回相应的错误
		code = e.ErrorUserNotFound
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//  2.用户存在，查看密码输入是否正确
	//  加密后进行对比
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 3. 认证通过，生成token（JWT签发token）
	token, err := utils.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 成功，返回token
	return serializer.Response{
		Status: code,
		Data:   serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:    e.GetMsg(code),
	}
}

// Update 修改用户信息
func (service *UserService) Update(ctx context.Context, uid uint) serializer.Response {
	var user *model.User
	var err error
	code := e.SUCCESS
	// 1. 查找用户
	userDao := dao.NewUserDao(ctx)
	user, err = userDao.GetUserById(uid)
	if service.NickName != "" {
		user.NickName = service.NickName
	}

	// 2.更新用户
	err = userDao.UpdateUserById(uid, user)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildUser(user),
		Msg:    e.GetMsg(code),
	}
}

// Post 头像更新
func (service *UserService) Post(ctx context.Context, uid uint, file multipart.File, filesize int64) serializer.Response {
	code := e.SUCCESS
	var user *model.User
	var err error
	userDao := dao.NewUserDao(ctx)
	// 1.查看用户是否存在
	user, err = userDao.GetUserById(uid)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// 2. 保存图片到本地路径，返回路径
	path, err := UploadAvatarToLocalStatic(file, uid, user.UserName)
	if err != nil {
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  path,
		}
	}
	user.Avatar = path
	err = userDao.UpdateUserById(uid, user)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildUser(user),
		Msg:    e.GetMsg(code),
	}
}

// Send 绑定邮箱
func (service *SendEmailService) Send(ctx context.Context, uid uint) serializer.Response {
	code := e.SUCCESS
	var address string
	var notice *model.Notice // ①绑定邮箱 ②修改密码，都需要模板通知
	// 1. 重新生成一个emailtoken
	token, err := utils.GenerateEmailToken(uid, service.OperationType, service.Email, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 2. 获取邮箱模板
	noticeDao := dao.NewNoticeDao(ctx)
	notice, err = noticeDao.GetNoticeById(service.OperationType)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// 3.创建一个新的邮件并发送
	address = conf.ValidEmail + token // 发送方
	// 对获取到的模板 进行内容的替换
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "Email", address, -1)
	// 导入依赖go get gopkg.in/mail.v2， 对邮件进行填充
	m := mail.NewMessage()
	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "subject")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
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

// Valid 验证邮箱
func (service *ValidEmailService) Valid(ctx context.Context, token string) serializer.Response {
	var userID uint
	var email string
	var password string
	var operationType uint
	code := e.SUCCESS

	// 1.验证token
	if token == "" {
		code = e.InvalidParams
	} else {
		claims, err := utils.ParseEmailToken(token)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		} else {
			userID = claims.UserID
			email = claims.Email
			password = claims.Password
			operationType = claims.OperationType
		}
	}
	// token有问题
	if code != e.SUCCESS {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 2. 用户token信息无误，进行对应操作
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(userID)
	if err != nil {
		return serializer.Result(e.ErrorDatabase)
	}

	if operationType == 1 {
		// 1 : 绑定邮箱
		user.Email = email
	} else if operationType == 2 {
		// 2 : 解绑邮箱
		user.Email = ""
	} else if operationType == 3 {
		// 3 : 修改密码
		err = user.SetPassword(password)
		if err != nil {
			return serializer.Result(e.ErrorDatabase)
		}
	}
	// 3.  保存用户信息
	err = userDao.UpdateUserById(userID, user)
	if err != nil {
		return serializer.Result(e.ErrorDatabase)
	}

	// 4.成功，返回用户信息
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}

}

// Show 查看用户金额
func (service *ShowMoneyService) Show(ctx context.Context, uid uint) serializer.Response {
	var code int = e.SUCCESS
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uid)
	if err != nil {
		return serializer.Result(e.ErrorDatabase)
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildMoney(user, service.Key),
		Msg:    e.GetMsg(code),
	}
}
