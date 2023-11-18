package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/news"
)

func NewsRouterInit(r *gin.Engine) {
	userRouters := r.Group("/news")
	{
		userRouters.GET("/:news_id", news.NewsController{}.NewsInfoDetail)     // 新闻详情
		userRouters.POST("/news_collect", news.NewsController{}.NewsCollect)   // 新闻收藏/取消收藏
		userRouters.POST("/followed_user", news.NewsController{}.NewsFollowed) // 用户关注/取消关注
	}
}
