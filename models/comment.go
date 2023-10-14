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

type CommentLike struct {
	BaseModel
	CommentID uint64
	UserID    uint64
}
