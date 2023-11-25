package admin

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"holiday/dao"
	"holiday/models"
	"holiday/utils"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type AdminController struct {
}

func (AdminController) Login(ctx *gin.Context) {
	var requestJSON struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	err := ctx.BindJSON(&requestJSON)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "传参错误"})
		zap.L().Error("传参错误: " + err.Error())
		return
	}

	var user models.User
	err = dao.DB.Where("mobile = ?", requestJSON.UserName).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "数据查询失败"})
		zap.L().Error("数据查询失败: " + err.Error())
		return
	}

	if user.PasswordHash != utils.EncryptPasswd(requestJSON.Password) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "密码错误"})
		return
	}
	if user.IsAdmin != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "权限不够"})
		return
	}
	session := sessions.Default(ctx)
	session.Set("user_id", user.ID)
	session.Set("nick_name", user.NickName)
	session.Set("mobile", user.Mobile)
	session.Set("is_admin", true)
	err = session.Save()
	if err != nil {
		zap.L().Error("session存储失败: " + err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "登录成功"})
}

func (AdminController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Set("user_id", nil)
	session.Set("nick_name", nil)
	session.Set("mobile", nil)
	session.Set("is_admin", nil)
	err := session.Save()
	if err != nil {
		zap.L().Error("session更新失败: " + err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "登出成功"})
}

func (AdminController) Index(ctx *gin.Context) {
	userInfo, _ := ctx.Get("user")

	user, ok := userInfo.(models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "断言错误"})
		zap.L().Error("断言错误")
		return
	}

	if user.IsAdmin != 1 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "权限不够"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "user": user.ToAdminDict()})
}

func (AdminController) NewsReview(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		zap.L().Error("string转换int失败: " + err.Error())
		page = 1
	}
	perPage, err := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	if err != nil {
		zap.L().Error("string转换int失败: " + err.Error())
		page = 1
	}
	keyword := ctx.Query("keyword")
	offset := (page - 1) * 10

	var (
		news      = make([]models.News, 0, 10)
		newsCount int64
	)
	query := "status != 0"
	if keyword != "" {
		query = strings.Join([]string{query, fmt.Sprintf("title like '%%%v%%'", keyword)}, " AND ")
	}

	err = dao.DB.Where(query).Offset(offset).Limit(10).Find(&news).Count(&newsCount).Error
	if err != nil {
		zap.L().Error("查询失败: " + err.Error())
	}

	newsList := make([]gin.H, 0, 10)
	for _, n := range news {
		newsList = append(newsList, n.ToReviewDict())
	}

	totalPage := 0
	if perPage != 0 {
		totalPage = int(math.Ceil(float64(newsCount) / float64(perPage)))
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "data": gin.H{
		"total_page":   totalPage,
		"current_page": page,
		"news_list":    newsList,
		"count":        newsCount,
	}})
}

func (AdminController) NewsReviewDetail(ctx *gin.Context) {
	newsID := ctx.Query("news_id")
	if newsID == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数错误"})
		return
	}
	var news models.News
	err := dao.DB.First(&news, newsID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "数据不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询错误"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "news_list": news.ToDict()})
}

func (AdminController) NewsReviewAudit(ctx *gin.Context) {
	var requestJSON struct {
		NewsID string `json:"news_id"`
		Action string `json:"action"`
		Reason string `json:"reason"`
	}
	var actionDict = map[string]bool{"accept": true, "reject": true}
	err := ctx.BindJSON(&requestJSON)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "传参错误"})
		zap.L().Error("传参错误: " + err.Error())
		return
	}
	if !actionDict[requestJSON.Action] {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "传参错误"})
		return
	}

	var news models.News
	err = dao.DB.First(&news, requestJSON.NewsID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "数据不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询错误"})
		return
	}

	if requestJSON.Action == "accept" {
		err = dao.DB.Model(&news).Update("status", 0).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "更新错误"})
			zap.L().Error("更新错误: " + err.Error())
			return
		}
	} else {
		if requestJSON.Reason == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数缺失"})
			return
		}
		err = dao.DB.Model(&news).Updates(models.News{Status: -1, Reason: requestJSON.Reason}).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "更新错误"})
			zap.L().Error("更新错误: " + err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "审核成功"})
}
