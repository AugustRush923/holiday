package main

import (
	"holiday/config"
	_ "holiday/dao"
	"holiday/middlewares"
	"holiday/routers"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	config.InitLogger()
	// 创建服务
	// gin.SetMode("release")
	r := gin.New()
	r.Use(middlewares.GinZapLogger(), middlewares.GinRecovery(zap.L(), false))

	// 路由注册
	routers.UserRouterInit(r)

	r.Run(strings.Join([]string{config.Cfg.Section("app").Key("ip_address").String(), config.Cfg.Section("app").Key("port").String()}, ":"))
}
