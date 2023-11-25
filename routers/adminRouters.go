package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/admin"
	"holiday/middlewares"
)

func AdminRouterInit(r *gin.Engine) {
	adminRouters := r.Group("/admin")
	{
		adminRouters.POST("/login", admin.AdminController{}.Login) // admin页面登录

	}
	adminLoginRequireRouters := r.Group("/admin", middlewares.LoginRequireMiddleware)
	{
		adminLoginRequireRouters.GET("/index", admin.AdminController{}.Index)                         // admin页面首页
		adminLoginRequireRouters.POST("/logout", admin.AdminController{}.Logout)                      // admin页面登出
		adminLoginRequireRouters.GET("/news_review", admin.AdminController{}.NewsReview)              // 待审核新闻列表
		adminLoginRequireRouters.GET("/news_review/detail", admin.AdminController{}.NewsReviewDetail) // 待审核新闻详情
		adminLoginRequireRouters.PUT("/news_review/audit", admin.AdminController{}.NewsReviewAudit)   // 审核新闻
		adminLoginRequireRouters.GET("/news_list", admin.AdminController{}.NewsList)                  // 新闻列表
		adminLoginRequireRouters.GET("/news/:news_id", admin.AdminController{}.NewsReviewDetail)      // 新闻详情
		adminLoginRequireRouters.PUT("/news/:news_id", admin.AdminController{}.EditNews)              // 新闻编辑
		adminLoginRequireRouters.GET("/categories", admin.AdminController{}.CategoryList)             // 分类列表
		adminLoginRequireRouters.PUT("/category/:category_id", admin.AdminController{}.EditCategory)  // 分类编辑
		adminLoginRequireRouters.POST("/category", admin.AdminController{}.AddCategory)               // 分类新增
		adminLoginRequireRouters.GET("/user_count", admin.AdminController{}.UserCount)                // 用户统计
		adminLoginRequireRouters.GET("/user_list", admin.AdminController{}.UserList)                  // 用户管理列表
	}
}
