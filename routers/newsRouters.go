package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/news"
)

func NewsRouterInit(r *gin.Engine) {
	newsRouters := r.Group("/news")
	{
		newsRouters.GET("/:news_id", news.NewsController{}.NewsInfoDetail)       // 新闻详情
		newsRouters.POST("/news_collect", news.NewsController{}.NewsCollect)     // 新闻收藏/取消收藏
		newsRouters.POST("/followed_user", news.NewsController{}.NewsFollowed)   // 用户关注/取消关注
		newsRouters.POST("/news_comment", news.NewsController{}.NewsComment)     // 新闻评论
		newsRouters.POST("/comment_like", news.NewsController{}.NewsCommentLike) // 新闻评论点赞/取消点赞
	}
}
