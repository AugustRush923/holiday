package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/news"
)

func NewsRouterInit(r *gin.Engine) {
	userRouters := r.Group("/news")
	{
		userRouters.GET("/:news_id", news.NewsController{}.NewsInfoDetail)
	}
}
