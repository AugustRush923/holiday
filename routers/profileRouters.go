package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/profile"
)

func ProfileRouterInit(r *gin.Engine) {
	profileRouters := r.Group("/profile")
	{
		profileRouters.GET("/info", profile.ProfileController{}.UserProfile)               // 用户中心
		profileRouters.GET("/user_info", profile.ProfileController{}.GetUserInfo)          // 查询用户信息
		profileRouters.PUT("/user_info", profile.ProfileController{}.UpdateUserInfo)       // 更新用户信息                                          // 更新用户信息
		profileRouters.PUT("/change_password", profile.ProfileController{}.ChangePassword) // 更换密码
		profileRouters.GET("/collection", profile.ProfileController{}.GetCollection)       // 收藏列表
		profileRouters.GET("/news_release", profile.ProfileController{}.GetNewsRelease)    // 新闻发布新闻分类下拉框列表
		profileRouters.POST("/news_release")
	}
}
