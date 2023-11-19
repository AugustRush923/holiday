package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"holiday/config"
	_ "holiday/dao"
	"holiday/middlewares"
	"holiday/routers"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// session设置
	// 创建基于cookie的存储引擎，[]byte 参数是用于加密的密钥
	store := cookie.NewStore([]byte(config.Cfg.Section("app").Key("secret_key").String()))
	// 初始化日志
	config.InitLogger()
	// 创建服务
	// gin.SetMode("release")
	r := gin.New()
	r.Use(middlewares.GinZapLogger(), middlewares.GinRecovery(zap.L(), false))
	// 设置session中间件，holidaySession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	r.Use(sessions.Sessions("holidaySession", store))

	// 路由注册
	routers.UserRouterInit(r)
	routers.NewsRouterInit(r)
	routers.IndexRouterInit(r)

	err := r.Run(strings.Join([]string{config.Cfg.Section("app").Key("ip_address").String(), config.Cfg.Section("app").Key("port").String()}, ":"))
	if err != nil {
		zap.L().Error("服务启动失败...")
		os.Exit(1)
	}
}
