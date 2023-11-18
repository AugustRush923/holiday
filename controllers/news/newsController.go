package news

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"holiday/dao"
	"holiday/models"
	"holiday/utils"
	"net/http"
	"strconv"
)

type NewsController struct {
}

func (NewsController) NewsInfoDetail(ctx *gin.Context) {
	news := models.News{}
	newsId := ctx.Param("news_id")
	// 初始化session对象
	session := sessions.Default(ctx)
	//	获取用户信息
	user := models.User{}
	userId := session.Get("user_id")
	if userId != nil {
		dao.DB.Where("id=?", userId).Find(&user)
	}
	userIsEmpty := user == models.User{}
	// 查询新闻数据
	err := dao.DB.First(&news, newsId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "查询不存在"})
			return
		}
		zap.L().Error("查询错误:" + err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "查询错误"})
		return
	}
	isEmpty := news == models.News{}
	if !isEmpty {
		news.Clicks += 1
		dao.DB.Save(&news)
	}

	//获取点击排行数据
	var clickNews []models.News
	dao.DB.Order("clicks Desc").Limit(10).Find(&clickNews)

	//获取当前新闻的评论
	var comments []models.Comment
	dao.DB.Where("news_id=?", news.ID).Order("create_time Desc").Find(&comments)
	var commentsID []uint
	var commentLikes []models.CommentLike
	var commentLikesID []uint64
	if !userIsEmpty {
		for _, comment := range comments {
			commentsID = append(commentsID, comment.ID)
		}
		if len(commentsID) > 0 {
			//	取到当前用户在当前新闻的所有评论点赞的记录
			dao.DB.Where("comment_id in (?) AND user_id=?", commentsID, user.ID).Find(&commentLikes)
			//取出记录中所有的评论id
			for _, like := range commentLikes {
				commentLikesID = append(commentLikesID, like.CommentID)
			}
		}
	}
	var commentList []map[string]any
	for _, comment := range comments {
		commentDict := comment.ToDict()
		commentDict["is_like"] = false
		//判断用户是否点赞该评论
		inSlice := utils.IsUintInSlice(commentLikesID, uint64(comment.ID))
		if !userIsEmpty && inSlice {
			commentDict["is_like"] = true
		}
		commentList = append(commentList, commentDict)
	}
	//获取分类信息
	var categories []models.Category
	dao.DB.Find(&categories)

	//当前登录用户是否关注当前新闻作者
	isFollowed := false
	//判断是否收藏该新闻，默认值为 false
	isCollected := false

	if !userIsEmpty {
		userFans := models.UserFans{}
		dao.DB.Where("followed_id=? AND follower_id=?", user.ID, news.UserID).Find(&userFans)
		userFansIsEmpty := userFans == models.UserFans{}
		if !userFansIsEmpty {
			isFollowed = true
		}

		userCollection := models.UserCollection{}
		dao.DB.Where("news_id = ? AND user_id = ?", news.ID, user.ID).Find(&userCollection)
		userCollectionIsEmpty := userCollection == models.UserCollection{}
		if !userCollectionIsEmpty {
			isCollected = true
		}
	}
	responseData := map[string]any{
		"user_info":       user.ToDict(),
		"news":            news.ToDict(),
		"click_news_list": clickNews,
		"comments":        commentList,
		"categories":      categories,
		"is_collected":    isCollected,
		"is_followed":     isFollowed,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "data": responseData})
}

func (NewsController) NewsCollect(ctx *gin.Context) {
	session := sessions.Default(ctx)
	//	获取用户信息
	user := models.User{}
	userId := session.Get("user_id")
	if userId != nil {
		dao.DB.Where("id=?", userId).Find(&user)
	}
	userIsEmpty := user == models.User{}
	if userIsEmpty {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "用户未登录"})
		return
	}
	// 获取参数
	type RequestJson struct {
		Action string `json:"action"`
		NewsID string `json:"news_id"`
	}
	var requestJson RequestJson
	err := ctx.BindJSON(&requestJson)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数错误"})
		return
	}
	action := map[string]bool{"collect": true, "cancel": true}
	if !action[requestJson.Action] {
		ctx.JSON(http.StatusForbidden, gin.H{"status": false, "message": "不被允许的操作"})
		return
	}

	news := models.News{}
	err = dao.DB.First(&news, requestJson.NewsID).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "数据不存在"})
		return
	}

	if requestJson.Action == "collect" {
		var exist int64
		dao.DB.Where("user_id = ? AND news_id = ?", user.ID, news.ID).Find(&models.UserCollection{}).Count(&exist)
		if exist > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "已收藏"})
			return
		}
		userCollect := models.UserCollection{UserID: uint64(user.ID), NewsID: uint64(news.ID)}
		err = dao.DB.Create(&userCollect).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "收藏失败"})
			zap.L().Error("收藏失败: " + err.Error())
			return
		}
	} else {
		err = dao.DB.Where("user_id = ? AND news_id = ?", user.ID, news.ID).Delete(&models.UserCollection{}).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "取消收藏失败"})
			zap.L().Error("取消收藏失败: " + err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "操作成功"})
	return
}

func (NewsController) NewsFollowed(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	user := models.User{}
	if userID != nil {
		dao.DB.Find(&user, userID)
	}
	userIsEmpty := user == models.User{}
	if userIsEmpty {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "用户未登录"})
		return
	}
	var action = map[string]bool{"follow": true, "unfollowed": true}
	type RequestJSON struct {
		UserID string `json:"user_id"`
		Action string `json:"action"`
	}
	var requestJson RequestJSON
	err := ctx.BindJSON(&requestJson)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数不全"})
		return
	}
	if !action[requestJson.Action] {
		ctx.JSON(http.StatusForbidden, gin.H{"status": false, "message": "不被允许的操作"})
		return
	}

	//查询要关注的用户信息
	var targetUser models.User
	dao.DB.First(&targetUser, requestJson.UserID)
	targetUserIsEmpty := targetUser == models.User{}
	if targetUserIsEmpty {
		ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "关注的用户不存在"})
		return
	}

	if user.ID == targetUser.ID {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "不能对自己进行操作"})
		return
	}

	if requestJson.Action == "follow" { // 关注
		userFans := models.UserFans{}
		dao.DB.Where("follower_id = ? AND followed_id = ?", user.ID, targetUser.ID).Find(&userFans)
		userFansIsEmpty := userFans == models.UserFans{}
		if !userFansIsEmpty {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "已关注"})
			return
		}
		err = dao.DB.Create(&models.UserFans{
			FollowerID: uint64(user.ID),
			FollowedID: uint64(targetUser.ID),
		}).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "关注失败"})
			zap.L().Error("用户关注出错: " + err.Error())
			return
		}
	} else { // 取消关注
		err = dao.DB.Where("follower_id = ? AND followed_id = ?", user.ID, targetUser.ID).Delete(&models.UserFans{}).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "取消关注失败"})
			zap.L().Error("用户取消关注出错: " + err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "操作成功"})
}

func (NewsController) NewsComment(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	var user models.User
	if userID != nil {
		err := dao.DB.First(&user, userID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "用户不存在"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询失败"})
			zap.L().Error("查询失败: " + err.Error())
			return
		}
	}

	userIsEmpty := user == models.User{}
	if userIsEmpty {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "用户未登录"})
		return
	}

	type RequestJSON struct {
		NewsID   string `json:"news_id"`
		Comment  string `json:"comment"`
		ParentID string `json:"parent_id"`
	}
	var requestJSON RequestJSON
	err := ctx.BindJSON(&requestJSON)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数错误"})
		return
	}

	//	查询新闻信息
	var news models.News
	err = dao.DB.First(&news, requestJSON.NewsID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "数据不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询错误"})
		zap.L().Error("查询失败：" + err.Error())
		return
	}

	comment := models.Comment{UserID: uint64(user.ID), NewsID: uint64(news.ID), Content: requestJSON.Comment}
	if requestJSON.ParentID != "" {
		parentID, strconvErr := strconv.Atoi(requestJSON.ParentID)
		if strconvErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "ParentID转换失败"})
			zap.L().Error("ParentID转换失败" + err.Error())
			return
		}
		uintID := uint64(parentID)
		comment.ParentID = &uintID
	}

	err = dao.DB.Create(&comment).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "新增评论失败"})
		zap.L().Error("新增评论失败：" + err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "success", "data": comment.ToDict()})
}

func (NewsController) NewsCommentLike(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	var user models.User
	if userID != nil {
		err := dao.DB.First(&user, userID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "用户不存在"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询失败"})
			zap.L().Error("查询失败: " + err.Error())
			return
		}
	}
	userIsEmpty := user == models.User{}
	if userIsEmpty {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "用户未登录"})
		return
	}

	action := map[string]bool{"add": true, "remove": true}

	type RequestJSON struct {
		//NewsID    string `json:"news_id"`
		CommentID string `json:"comment_id"`
		Action    string `json:"action"`
	}
	var requestJSON RequestJSON
	err := ctx.BindJSON(&requestJSON)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "JSON数据绑定失败"})
		zap.L().Error("JSON数据绑定失败: " + err.Error())
		return
	}

	if !action[requestJSON.Action] {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "参数错误"})
		return
	}

	//	查询评论数据
	var comment models.Comment
	err = dao.DB.First(&comment, requestJSON.CommentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": false, "message": "用户不存在"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "查询失败"})
		zap.L().Error("查询失败: " + err.Error())
		return
	}

	var commentLike models.CommentLike
	err = dao.DB.Where("comment_id = ? AND user_id = ?", comment.ID, user.ID).First(&commentLike).Error
	if requestJSON.Action == "add" {
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				dao.DB.Create(&models.CommentLike{
					CommentID: uint64(comment.ID),
					UserID:    uint64(user.ID),
				})
				// 增加点赞条数
				comment.LikeCount += 1
				dao.DB.Model(&comment).Update("like_count", comment.LikeCount)
			} else {
				zap.L().Error("查询失败: " + err.Error())
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "已点过赞"})
			return
		}
	} else {
		// 删除点赞数据
		if err == nil {
			dao.DB.Where("comment_id = ? AND user_id = ?", comment.ID, user.ID).Delete(&commentLike)
			//	减少点赞条数
			comment.LikeCount -= 1
			err = dao.DB.Model(&comment).Update("like_count", comment.LikeCount).Error
			if err != nil {
				zap.L().Error("更新失败：" + err.Error())
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "未点过赞"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": true, "message": "操作成功"})
}
