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
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.Section("mysql").Key("username").String(),
		config.Cfg.Section("mysql").Key("password").String(),
		config.Cfg.Section("mysql").Key("hostname").String(),
		config.Cfg.Section("mysql").Key("port").String(),
		config.Cfg.Section("mysql").Key("database").String(),
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("数据库连接失败：", err)
		os.Exit(0)
	}

}
