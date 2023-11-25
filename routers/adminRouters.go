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
		adminRouters.PUT("/news_review/audit", admin.AdminController{}.NewsReviewAudit)   // 审核新闻
		adminRouters.GET("/news_list", admin.AdminController{}.NewsList)                  // 新闻列表
		adminRouters.GET("/news/:news_id", admin.AdminController{}.NewsReviewDetail)      // 新闻详情
		adminRouters.PUT("/news/:news_id", admin.AdminController{}.EditNews)              // 新闻编辑
		adminRouters.GET("/categories", admin.AdminController{}.CategoryList)             // 分类列表
		adminRouters.PUT("/category/:category_id", admin.AdminController{}.EditCategory)  // 分类编辑
		adminRouters.POST("/category", admin.AdminController{}.AddCategory)               // 分类新增
	}
	adminUseMiddlewareRouters := r.Group("/admin", middlewares.LoginRequireMiddleware)
	{
		adminUseMiddlewareRouters.GET("/index", admin.AdminController{}.Index) // admin页面首页
	}
}
