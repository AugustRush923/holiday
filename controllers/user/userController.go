package user

import (
	"holiday/dao"
	"holiday/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (UserController) CreateUser(ctx *gin.Context) {
	user := models.User{}
	ctx.ShouldBind(&user)
	if user.Avatar == "" {
		user.Avatar = "默认头像地址"
	}
	// TODO: 密码加密存储
	if user.UserName != "" {
		// 名字是否重复
		if dao.DB.Model(&models.User{}).Where("user_name=? AND is_deleted=0", user.UserName).First(&user).Error == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "名字重复"})
			return
		}
	}

	if user.Email != "" {
		// 邮件是否重复
		if dao.DB.Model(&models.User{}).Where("email=? AND is_deleted=0", user.Email).First(&user).Error == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "邮件重复"})
			return
		}
	}

	if user.Mobile != "" {
		// 手机号是否重复
		if dao.DB.Model(&models.User{}).Where("mobile=? AND is_deleted = 0", user.Mobile).First(&user).Error == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "手机号重复"})
			return
		}
	}
	err := dao.DB.Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "id": user.ID})
}

func (UserController) GetUserDetail(ctx *gin.Context) {
	user := models.User{}
	id := ctx.Param("id")
	err := dao.DB.Where("is_deleted=0").First(&user, id).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "查询不存在"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "data": user})
}
