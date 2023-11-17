package user

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"holiday/dao"
	"holiday/models"
	"holiday/utils"
	"net/http"
	"strconv"
	"time"
)

type UserController struct{}

func (UserController) CreateUser(ctx *gin.Context) {
	user := models.User{}
	bindErr := ctx.ShouldBind(&user)

	if bindErr != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "参数错误"})
		return
	}
	if user.AvatarUrl == "" {
		user.AvatarUrl = "默认头像地址"
	}
	// 密码加密存储
	encryptPasswd := utils.EncryptPasswd(user.PasswordHash)
	user.PasswordHash = encryptPasswd

	if user.NickName != "" {
		// 名字是否重复
		if dao.DB.Model(&models.User{}).Where("nick_name=?", user.NickName).First(&user).Error == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "名字重复"})
			return
		}
	}

	if user.Mobile != "" {
		// 手机号是否重复
		if dao.DB.Model(&models.User{}).Where("mobile=? ", user.Mobile).First(&user).Error == nil {
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
	err := dao.DB.First(&user, id).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "查询不存在"})
		return
	}
	var data map[string]any
	if user.IsAdmin == 1 {
		data = user.ToDict()
	} else {
		data = user.ToDict()
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "data": data})
}

func (UserController) Login(ctx *gin.Context) {
	// 初始化session对象
	session := sessions.Default(ctx)
	var body struct {
		Username string `json:"username"`
		Passwd   string `json:"passwd"`
	}
	//body.Username = ctx.Query("username")
	//body.Passwd = ctx.Query("passwd")
	user := models.User{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数不全"})
		return
	}
	if err := dao.DB.Model(&models.User{}).Where("nick_name=?", body.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "用户未找到"})
		return
	}
	fmt.Println(utils.EncryptPasswd(body.Passwd))
	if utils.EncryptPasswd(body.Passwd) != user.PasswordHash {
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "密码不匹配"})
		return
	}
	// 本地测试时如果只有localhost的话，没法设置上Cookie
	// 解决方案是 把运行的端口也带上 localhost:9000 或者 值为空
	ctx.SetCookie("user_id", strconv.Itoa(int(user.ID)), 315360000, "/", "", false, true)
	ctx.SetCookie("nickname", user.NickName, 315360000, "/", "", false, true)
	ctx.SetCookie("mobile", user.Mobile, 315360000, "/", "", false, true)
	session.Set("user_id", user.ID)
	session.Set("nick_name", user.NickName)
	session.Set("mobile", user.Mobile)
	err := session.Save()
	if err != nil {
		zap.L().Error("session设置错误: " + err.Error())
	}
	updateErr := dao.DB.Model(&user).Update("last_login", time.Now()).Error
	if updateErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "false", "message": "更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": false, "message": "登录成功"})
}

func (UserController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	ctx.SetCookie("user_id", "", -1, "/", "", false, true)
	ctx.SetCookie("nickname", "", -1, "/", "", false, true)
	ctx.SetCookie("mobile", "", -1, "/", "", false, true)
	session.Delete("user_id")
	session.Delete("nick_name")
	session.Delete("mobile")
	err := session.Save()
	if err != nil {
		zap.L().Error("session设置错误: " + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "退出成功"})
}

func (UserController) Register(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := models.User{}
	var body struct {
		Mobile string `json:"mobile"`
		Passwd string `json:"passwd"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数不全"})
		return
	}
	user.Mobile = body.Mobile
	user.NickName = body.Mobile
	user.PasswordHash = utils.EncryptPasswd(body.Passwd)
	user.AvatarUrl = "default avatar"
	err := dao.DB.Model(models.User{}).Select("mobile", "nick_name", "password_hash", "avatar_url").Create(&user).Error

	if err != nil {
		zap.L().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "创建失败"})
		return
	}

	ctx.SetCookie("user_id", strconv.Itoa(int(user.ID)), 315360000, "/", "", false, true)
	ctx.SetCookie("nickname", user.NickName, 315360000, "/", "", false, true)
	ctx.SetCookie("mobile", user.Mobile, 315360000, "/", "", false, true)
	session.Set("user_id", user.ID)
	session.Set("nick_name", user.NickName)
	session.Set("mobile", user.Mobile)
	err = session.Save()
	if err != nil {
		zap.L().Error("session设置错误: " + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "注册成功"})
}
