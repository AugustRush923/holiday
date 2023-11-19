package routers

import (
	"github.com/gin-gonic/gin"
	"holiday/controllers/profile"
)

func ProfileRouterInit(r *gin.Engine) {
	profileRouters := r.Group("/profile")
	{
		profileRouters.GET("/info", profile.ProfileController{}.UserProfile)         // 用户中心
		profileRouters.GET("/user_info", profile.ProfileController{}.GetUserInfo)    // 查询用户信息
		profileRouters.PUT("/user_info", profile.ProfileController{}.UpdateUserInfo) // 更新用户信息                                          // 更新用户信息
		profileRouters.GET("/change_password")
		profileRouters.POST("/change_password")
		profileRouters.GET("/collection")
		profileRouters.GET("/news_release")
		profileRouters.POST("/news_release")
	}
}
