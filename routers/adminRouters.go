package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/admin"
)

func AdminRouterInit(r *gin.Engine) {
	adminRouters := r.Group("/admin")
	{
		adminRouters.POST("/login", admin.AdminController{}.Login)           // admin页面登录
		adminRouters.POST("/logout", admin.AdminController{}.Logout)         // admin页面登出
		adminRouters.GET("/index", admin.AdminController{}.Index)            // admin页面首页
		adminRouters.GET("/news_review", admin.AdminController{}.NewsReview) // 待审核新闻列表
	}
}
