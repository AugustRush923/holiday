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
	"time"
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
	newsID := ctx.Param("news_id")
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

func (AdminController) NewsList(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		zap.L().Error("string转化int失败: " + err.Error())
		page = 1
	}
	offset := (page - 1) * 10

	news := make([]models.News, 0, 10)
	var (
		newsCount int64
		query     = ""
		newsList  []gin.H
	)
	if keyword != "" {
		query = fmt.Sprintf("Title like \"%%%v%%\"", keyword)
	}
	dao.DB.Where(query).Offset(offset).Limit(10).Order("create_time Desc").Find(&news).Count(&newsCount)

	totalPage := 0
	totalPage = int(math.Ceil(float64(newsCount) / float64(10)))

	for _, n := range news {
		newsList = append(newsList, n.ToBasicDict())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       true,
		"news_list":    newsList,
		"count":        newsCount,
		"current_page": page,
		"total_page":   totalPage,
	})
}

func (AdminController) CategoryList(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		zap.L().Error("string转化int失败: " + err.Error())
		page = 1
	}
	offset := (page - 1) * 10

	var (
		categoryCount int64
		category      = make([]models.Category, 0, 10)
		categoryList  []gin.H
	)
	dao.DB.Where("id != 1").Offset(offset).Limit(10).Find(&category).Count(&categoryCount)

	for _, c := range category {
		categoryList = append(categoryList, c.ToDict())
	}

	totalPage := int(math.Ceil(float64(categoryCount) / float64(10)))
	ctx.JSON(http.StatusOK, gin.H{
		"status":        true,
		"category_list": categoryList,
		"count":         categoryCount,
		"current_page":  page,
		"total_page":    totalPage,
	})
}

func (AdminController) EditNews(ctx *gin.Context) {
	newsID, err := strconv.Atoi(ctx.Param("news_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数错误"})
		zap.L().Error("string转int失败: " + err.Error())
		return
	}

	var news = models.News{BaseModel: models.BaseModel{ID: uint(newsID)}}
	err = ctx.BindJSON(&news)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数校验失败"})
		zap.L().Error("参数校验失败: " + err.Error())
		return
	}
	err = dao.DB.Updates(&news).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "更新失败"})
		zap.L().Error("更新失败: " + err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "保存成功"})
}

func (AdminController) EditCategory(ctx *gin.Context) {
	categoryID, err := strconv.Atoi(ctx.Param("category_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数错误"})
		zap.L().Error("string转int失败: " + err.Error())
		return
	}
	category := models.Category{
		BaseModel: models.BaseModel{ID: uint(categoryID)},
	}

	err = ctx.BindJSON(&category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数校验失败"})
		zap.L().Error("参数校验失败: " + err.Error())
		return
	}

	err = dao.DB.Model(&category).Update("name", category.Name).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "更新失败"})
		zap.L().Error("更新失败: " + err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "保存成功"})
}

func (AdminController) AddCategory(ctx *gin.Context) {
	var category models.Category
	err := ctx.BindJSON(&category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数校验失败"})
		zap.L().Error("参数校验失败: " + err.Error())
		return
	}
	err = dao.DB.Create(&category).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "创建失败"})
		zap.L().Error("创建失败: " + err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "创建成功"})
}

func (AdminController) UserCount(ctx *gin.Context) {
	// 查询总人数
	var totalCount int64
	dao.DB.Model(&models.User{}).Where("is_admin != 1").Count(&totalCount)
	// 查询月新增人数
	var monCount int64
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	dao.DB.Model(&models.User{}).Where("create_time >= ? AND is_admin = 0", firstDayOfMonth).Count(&monCount)

	// 查询日新增数
	var dayCount int64
	firstDayOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dao.DB.Model(&models.User{}).Where("create_time >= ? AND is_admin = 0", firstDayOfDay).Count(&dayCount)

	// 查询每日活跃人数
	var (
		activeDate  = make([]string, 0)
		activeCount = make([]int64, 0)
	)

	for i := 0; i < 31; i++ {
		var count int64
		beginDate := time.Date(now.Year(), now.Month(), now.Day()-i, 0, 0, 0, 0, now.Location())
		endDate := time.Date(now.Year(), now.Month(), now.Day()-i, 23, 59, 59, 99999999, now.Location())
		activeDate = append(activeDate, beginDate.Format("2006-01-02"))

		dao.DB.Model(&models.User{}).Where("is_admin = 0 AND last_login BETWEEN ? AND ?", beginDate, endDate).Count(&count)
		activeCount = append(activeCount, count)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       true,
		"total_count":  totalCount,
		"mon_count":    monCount,
		"day_count":    dayCount,
		"active_date":  activeDate,
		"active_count": activeCount,
	})
}
