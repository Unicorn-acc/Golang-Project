package service

import (
	"example.com/unicorn-acc/conf"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strconv"
)

func UploadAvatarToLocalStatic(file multipart.File, userId uint, username string) (filepath string, err error) {
	bId := strconv.Itoa(int(userId)) // int => string
	basePath := "." + conf.AvatarPath + "user" + bId + "/"
	// 如果路径不存在，则创建路径
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + username + ".jpg"
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err // 读取失败
	}
	err = os.WriteFile(avatarPath, content, 0666)
	if err != nil {
		return "", err
	}
	return "user" + bId + "/" + username + ".jpg", err
}

// DirExistOrNot 判断文件是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 755)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func UploadProductToLocalStatic(file multipart.File, bossId uint, productName string) (filePath string, err error) {
	bId := strconv.Itoa(int(bossId))
	basePath := "." + conf.ProductPhotoPath + "boss" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	productPath := basePath + productName + ".jpg"
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(productPath, content, 0666)
	if err != nil {
		return "", err
	}
	return "boss" + bId + "/" + productName + ".jpg", err
}
