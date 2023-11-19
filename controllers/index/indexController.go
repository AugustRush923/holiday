package index

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"holiday/dao"
	"holiday/models"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type IndexController struct {
}

func (IndexController) GetIndex(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var user models.User
	userID := session.Get("user_id")
	if userID != nil {
		err := dao.DB.First(&user, userID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "数据不存在"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询错误"})
			zap.L().Error("查询错误: " + err.Error())
			return
		}
	}

	//	获取点击排行榜数据
	var newsList []models.News
	err := dao.DB.Order("clicks Desc").Limit(10).Find(&newsList).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询错误"})
		zap.L().Error("查询错误: " + err.Error())
		return
	}

	var clickNewsList []map[string]any
	for _, news := range newsList {
		clickNewsList = append(clickNewsList, news.ToBasicDict())
	}

	// 获取分类数据
	var categories []models.Category
	err = dao.DB.Find(&categories).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询错误"})
		zap.L().Error("查询错误: " + err.Error())
		return
	}
	var categoriesList []map[string]any
	for _, category := range categories {
		categoriesList = append(categoriesList, category.ToDict())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": map[string]any{
			"categories":      categoriesList,
			"click_news_list": clickNewsList,
			"user_info":       user.ToDict(),
		},
	})
}

func (IndexController) NewsList(ctx *gin.Context) {
	cid := ctx.Query("cid")
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "10")
	cidInt, err := strconv.Atoi(cid)
	if err != nil {
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		return
	}
	offset := (pageInt - 1) * perPageInt
	var (
		news      []models.News
		newsCount int64
	)
	query := []string{"status = 0"}
	if cidInt != 1 {
		query = append(query, fmt.Sprintf("category_id = %v", cidInt))
	}

	err = dao.DB.Where(strings.Join(query, " AND ")).Order("create_time Desc").Offset(offset).Limit(10).Find(&news).Count(&newsCount).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询错误"})
		zap.L().Error("查询错误: " + err.Error())
		return
	}

	var newsList []map[string]any
	for _, n := range news {
		newsList = append(newsList, n.ToBasicDict())
	}

	totalPage := 0
	if perPageInt != 0 {
		totalPage = int(math.Ceil(float64(newsCount) / float64(perPageInt)))
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true, "data": map[string]any{
		"news_list":    newsList,
		"count":        newsCount,
		"total_page":   totalPage,
		"current_page": pageInt,
	}})
}
