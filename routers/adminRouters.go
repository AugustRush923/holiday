package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/admin"
	"holiday/middlewares"
)

func AdminRouterInit(r *gin.Engine) {
	adminRouters := r.Group("/admin")
	{
		adminRouters.POST("/login", admin.AdminController{}.Login)                        // admin页面登录
		adminRouters.POST("/logout", admin.AdminController{}.Logout)                      // admin页面登出
		adminRouters.GET("/news_review", admin.AdminController{}.NewsReview)              // 待审核新闻列表
		adminRouters.GET("/news_review/detail", admin.AdminController{}.NewsReviewDetail) // 待审核新闻详情
		adminRouters.POST("/news_review/audit", admin.AdminController{}.NewsReviewAudit)  // 审核新闻
	}
	adminUseMiddlewareRouters := r.Group("/admin", middlewares.LoginRequireMiddleware)
	{
		adminUseMiddlewareRouters.GET("/index", admin.AdminController{}.Index) // admin页面首页
	}
}
