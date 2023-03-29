package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var DB *gorm.DB

// 根据配置文件进行Mysql连接
func DataBase(conn string) {
	fmt.Println(conn)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conn,  // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		// 打印日志
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表明不加s
		},
	})
	if err != nil {
		panic("数据库连接错误")
	}
	log.Println("数据库连接成功")
	// 设置数据库参数
	mysqldb, _ := db.DB()
	mysqldb.SetMaxIdleConns(20)                  //  设置连接池
	mysqldb.SetMaxOpenConns(100)                 //最大连接数
	mysqldb.SetConnMaxLifetime(time.Second * 30) // 设置最大连接时间
	DB = db
	migration()
}
