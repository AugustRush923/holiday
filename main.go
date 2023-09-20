package main

import (
	"holiday/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建服务
	r := gin.Default()
	// // 创建日志
	// logger, _ := zap.NewDevelopment()
	// r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	// r.Use(ginzap.RecoveryWithZap(logger, true))

	r.Run(strings.Join([]string{config.Cfg.Section("app").Key("ip_address").String(), config.Cfg.Section("app").Key("port").String()}, ":"))
}
