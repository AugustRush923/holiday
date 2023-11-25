package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"holiday/dao"
	"holiday/models"
	"net/http"
)

func LoginRequireMiddleware(ctx *gin.Context) {
	// 获取session对象
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	if userID == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "未登录"})
	}
	user := models.User{}
	err := dao.DB.First(&user, userID).Error
	if err != nil {
		zap.L().Error("查询失败: " + err.Error())
	}

	isUserEmpty := user == models.User{}
	if isUserEmpty {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "未登录"})
	}
	// 设置全局变量
	ctx.Set("user", user)

	ctx.Next()
}
