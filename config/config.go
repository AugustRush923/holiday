package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	Cfg *ini.File
	err error
)

func init() {
	// 读取配置信息
	Cfg, err = ini.Load("./settings.ini")
	if err != nil {
		fmt.Println("加载配置文件失败！", err)
		panic("读取配置信息失败")
	}
}
