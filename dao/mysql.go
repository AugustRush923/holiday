package dao

import (
	"fmt"
	"holiday/config"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

func init() {
	username := config.Cfg.Section("mysql").Key("username").String()
	passwd := config.Cfg.Section("mysql").Key("password").String()
	hostname := config.Cfg.Section("mysql").Key("hostname").String()
	port := config.Cfg.Section("mysql").Key("port").String()
	database := config.Cfg.Section("mysql").Key("database").String()

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", username, passwd, hostname, port, database)
	fmt.Println(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		os.Exit(0)
	}
}
