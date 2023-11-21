package profile

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"holiday/dao"
	"holiday/models"
	"holiday/utils"
	"net/http"
)

type ProfileController struct {
}

func (ProfileController) UserProfile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	user := models.User{}
	if userID != nil {
		err := dao.DB.First(&user).Error
		if err != nil {
			zap.L().Error("查询错误: " + err.Error())
		}
	} else {
		// 如果用户没登录跳转到首页
		ctx.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	data := map[string]any{
		"user_info": user.ToDict(),
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "data": data})
}

func (ProfileController) GetUserInfo(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	user := models.User{}
	if userID != nil {
		err := dao.DB.First(&user).Error
		if err != nil {
			zap.L().Error("查询错误: " + err.Error())
		}
	}

	data := map[string]any{
		"user_info": user.ToDict(),
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "data": data})
}

func (ProfileController) UpdateUserInfo(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	user := models.User{}
	if userID != nil {
		err := dao.DB.First(&user).Error
		if err != nil {
			zap.L().Error("查询错误: " + err.Error())
		}
	}
	var RequestJSON struct {
		NickName  string `json:"nick_name"`
		Signature string `json:"signature"`
		Gender    string `json:"gender"`
	}
	requestJSON := RequestJSON
	err := ctx.BindJSON(&requestJSON)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数不全"})
		return
	}

	err = dao.DB.Model(&user).Updates(models.User{
		NickName:  requestJSON.NickName,
		Signature: requestJSON.Signature,
		Gender:    requestJSON.Gender,
	}).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "更新失败"})
		zap.L().Error("更新失败: " + err.Error())
		return
	}

	session.Set("nick_name", requestJSON.NickName)
	err = session.Save()
	if err != nil {
		zap.L().Error("session更新失败: " + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "更新成功"})
}

func (ProfileController) ChangePassword(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	userID := session.Get("user_id")
	if userID != nil {
		err := dao.DB.First(&user, userID).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "用户查询错误"})
			zap.L().Error("用户查询错误: " + err.Error())
			return
		}
	}

	var RequestJSON struct {
		OldPasswd string `json:"old_passwd"`
		NewPasswd string `json:"new_passwd"`
	}
	var requestJSON = RequestJSON
	err := ctx.BindJSON(&requestJSON)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数错误"})
		return
	}

	if user.PasswordHash != utils.EncryptPasswd(requestJSON.OldPasswd) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "密码输入错误"})
		return
	}
	newPassword := utils.EncryptPasswd(requestJSON.NewPasswd)
	err = dao.DB.Model(&user).Update("password_hash", newPassword).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "更新密码失败"})
		zap.L().Error("更新密码失败: " + err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "更新密码成功"})
}
