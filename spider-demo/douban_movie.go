package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// 定义数据库连接的常量
const (
	USERNAME = "root"
	PASSWORD = "123456"
	HOST     = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "spider"
)

// 与数据库表字段相同的结构体
type MovieData struct {
	Title    string `json:"title"`
	Director string `json:"Director"`
	Picture  string `json:"Picture"`
	Actor    string `json:"Actor"`
	Year     string `json:"Year"`
	Score    string `json:"Score"`
	Quote    string `json:"Quote"`
}

var DB *sql.DB

// 初始化数据库连接
func InitDBBySQL() {
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(10)
	DB.SetMaxIdleConns(5)
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail")
		return
	}
	fmt.Println("connect success")
}

var gormdb *gorm.DB

func InitDBByGORM() {
	var err error
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	gormdb, err = gorm.Open("mysql", path)
	if err != nil {
		panic(err)
	}
	_ = gormdb.AutoMigrate(&NewD{})
	sqlDB := gormdb.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	fmt.Println("connect success")
}

func main() {
	//InitDBBySQL()
	InitDBByGORM()
	//豆瓣top250 ： https://movie.douban.com/top250
	for i := 0; i < 10; i++ {
		fmt.Printf("正在爬取第 %d 页的信息\n", i)
		Spider(strconv.Itoa(i * 25))
	}
}

// 爬虫爬取过程
func Spider(page string) {
	// todo 1. 发送请求
	// 构造客户端
	client := http.Client{}
	// 构造GET请求
	URL := "https://movie.douban.com/top250?start=" + page
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("构造Get请求失败： ", err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败: ", err)
	}
	defer resp.Body.Close()

	// todo 2. 解析网页
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("解析失败", err)
	}

	// todo 3. 获取节点信息
	// 举例：查看豆瓣电影的top1的title
	// 下面这句查看网页源码然后右键复制selector,使用Find()方法查看内容(记得后面有.text())
	// #content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
	//title1 := doc.Find("#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)").Text()
	//fmt.Println("title:", title1)
	////////////////////////////////////////////////////////////////////////////////////////
	// 循环列表，进行遍历
	doc.Find("#content > div > div.article > ol > li"). // 列表
								Each(func(i int, s *goquery.Selection) { // 在列表里面继续找
			var data MovieData
			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text() // 电影标题
			img := s.Find("div > div.pic > a > img")                                  // img图片，但是img标签在sec属性里面
			imgTmp, ok := img.Attr("src")                                             // 获得img中src中的值，不存在ok就为err
			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()         // 电影信息
			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text() // 电影评分
			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()        // 电影评论
			if ok {
				// todo 4. 保存信息
				director, actor, year := InfoSpite(info)
				data.Title = title
				data.Director = director
				data.Actor = actor
				data.Picture = imgTmp
				data.Year = year
				data.Score = score
				data.Quote = quote
				if InsertDataByGORM(data) {
					fmt.Printf("插入成功：%+v\n", data)
				} else {
					fmt.Printf("插入失败：%+v\n", data)
				}
			}
		})

}

// 将豆瓣info信息中的导演、演员等信息用 正则表达式 提取出来
func InfoSpite(info string) (director, actor, year string) {
	directorRe, _ := regexp.Compile(`导演:(.*)主演:`) // 正则 .*匹配一行所有
	director = string(directorRe.Find([]byte(info)))

	actorRe, _ := regexp.Compile(`主演:(.*)`) // 正则 .*匹配一行所有
	actor = string(actorRe.Find([]byte(info)))

	yearRe, _ := regexp.Compile(`(\d+)`) //正则表达式 \d+匹配数字 == 年份
	year = string(yearRe.Find([]byte(info)))
	return
}

// 插入数据库
func InsertDataBySQL(data MovieData) bool {
	tx, err := DB.Begin() // 开启数据库事务
	if err != nil {
		fmt.Println("开启数据库事务DB.Begin()失败:", err)
		return false
	}
	// 编写sql语句(需要提前把数据准备好)
	stmt, err := tx.Prepare("Insert INTO douban_movie(`Title`, `Director`, `Picture`, `Actor`, `Year`, `Score`, `Quote`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("数据准备tx.Prepare失败：", err)
		return false
	}
	_, err = stmt.Exec(data.Title, data.Director, data.Picture, data.Actor, data.Year, data.Score, data.Quote)
	if err != nil {
		fmt.Println("数据插入stmt.Exec失败：", err)
		return false
	}
	// 提交事务
	tx.Commit()
	return true
}

type NewD struct {
	gorm.Model
	Title    string `gorm:"type:varchar(255);not null;"`
	Director string `gorm:"type:varchar(256);not null;"`
	Picture  string `gorm:"type:varchar(256);not null;"`
	Actor    string `gorm:"type:varchar(256);not null;"`
	Year     string `gorm:"type:varchar(256);not null;"`
	Score    string `gorm:"type:varchar(256);not null;"`
	Quote    string `gorm:"type:varchar(256);not null;"`
}

func InsertDataByGORM(data MovieData) bool {
	NewA := NewD{
		Title:    data.Title,
		Director: data.Director,
		Picture:  data.Picture,
		Actor:    data.Actor,
		Year:     data.Year,
		Score:    data.Score,
		Quote:    data.Quote,
	}
	err := gormdb.Create(&NewA).Error
	if err != nil {
		return false
	}
	return true
}
