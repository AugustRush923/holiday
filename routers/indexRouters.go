package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/index"
)

func IndexRouterInit(r *gin.Engine) {
	indexRouters := r.Group("/")
	{
		indexRouters.GET("", index.IndexController{}.GetIndex)           // 首页
		indexRouters.GET("/news_list", index.IndexController{}.NewsList) // 新闻列表
	}
}
