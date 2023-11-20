package repositories

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MeetingDB *gorm.DB

func Init() {
	username := "root"
	password := "Mcyouyou1@"
	hostname := "127.0.0.1"
	port := 3306
	database := "glimmermeeting"

	// 创建MySQL数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, hostname, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect database: " + err.Error())
	}

	MeetingDB = db
}
