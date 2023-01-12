package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func init() {
	username := "chatgpt-dev"
	password := "CEtRH3nySTzZyiKK"
	host := "42.192.42.152"
	port := 3306
	dbname := "chatgpt-dev"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password,
		host, port, dbname)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("mysql connect error, error:" + err.Error())
	}

	// 设置数据库连接池参数
	sqlDB, _ := db.DB()
	// 连接池最大连接数
	sqlDB.SetMaxIdleConns(100)
	// 连接池最大允许空闲连接
	sqlDB.SetMaxIdleConns(20)
}

func GetDB() *gorm.DB {
	return db
}
