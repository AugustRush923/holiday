package routers

import (
	"holiday/controllers/user"

	"github.com/gin-gonic/gin"
)

func UserRouterInit(r *gin.Engine) {
	userRouters := r.Group("/user")
	{
		userRouters.POST("/", user.UserController{}.CreateUser)
		userRouters.GET("/:id", user.UserController{}.GetUserDetail)
	}
}
