package models

type News struct {
	BaseModel
	Title         string
	Source        string
	Digest        string
	Content       string
	Clicks        uint64
	IndexImageUrl string
	CategoryID    uint64
	UserID        uint64
	Status        uint8
	reason        string
}

func (News) TableName() string {
	return "info_news"
}
