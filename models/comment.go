package models

import "time"

type Comment struct {
	BaseModel
	UserID    uint64
	NewsID    uint64
	Content   string
	ParentID  *uint64
	LikeCount uint64
}

func (Comment) TableName() string {
	return "info_comment"
}

func (comment Comment) ToDict() map[string]any {
	isEmpty := comment == Comment{}
	if isEmpty {
		return make(map[string]any)
	}
	return map[string]any{
		"comment_id":      comment.ID,
		"created_time":    comment.CreatedTime,
		"comment_content": comment.Content,
		"news_id":         comment.NewsID,
		"like_count":      comment.LikeCount,
	}
}

type CommentLike struct {
	CommentID   uint64
	UserID      uint64
	CreatedTime time.Time `gorm:"column:create_time;type:datetime(0);autoCreateTime;comment:创建时间" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:update_time;type:datetime(0);autoUpdateTime;comment:更新时间" json:"updated_time"`
}

func (CommentLike) TableName() string {
	return "info_comment_like"
}
