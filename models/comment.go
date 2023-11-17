package models

type Comment struct {
	BaseModel
	UserID    uint64
	NewsID    uint64
	Content   string
	ParentID  uint64
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
	BaseModel
	CommentID uint64
	UserID    uint64
}
