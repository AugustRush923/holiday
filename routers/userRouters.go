package routers

import (
	"holiday/controllers/user"

	"github.com/gin-gonic/gin"
)

func UserRouterInit(r *gin.Engine) {
	userRouters := r.Group("/user")
	{
		userRouters.POST("/", user.UserController{}.CreateUser)       // 新建用户
		userRouters.GET("/:id", user.UserController{}.GetUserDetail)  // 用户详情
		userRouters.POST("/login", user.UserController{}.Login)       // 用户登录
		userRouters.POST("/logout", user.UserController{}.Logout)     // 用户退出
		userRouters.POST("/register", user.UserController{}.Register) // 用户注册
	}
}
