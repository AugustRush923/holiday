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
	Status        int8
	Reason        string
}

func (News) TableName() string {
	return "info_news"
}

func (news News) ToDict() map[string]any {
	isEmpty := news == News{}
	if isEmpty {
		return make(map[string]any)
	}

	return map[string]any{
		"news_id":              news.ID,
		"news_title":           news.Title,
		"news_digest":          news.Digest,
		"news_content":         news.Content,
		"news_clicks":          news.Clicks,
		"news_index_image_url": news.IndexImageUrl,
		"news_category_id":     news.CategoryID,
		"news_status":          news.Status,
		"news_reason":          news.Reason,
		"news_created_time":    news.CreatedTime,
		"news_updated_time":    news.UpdatedTime,
	}
}

func (news News) ToBasicDict() map[string]any {
	isEmpty := news == News{}
	if isEmpty {
		return make(map[string]any)
	}

	return map[string]any{
		"news_id":              news.ID,
		"news_title":           news.Title,
		"news_source":          news.Source,
		"news_digest":          news.Digest,
		"news_clicks":          news.Clicks,
		"news_index_image_url": news.IndexImageUrl,
		"news_created_time":    news.CreatedTime.Format("2006-01-02 15:04:05"),
	}
}

func (news News) ToReviewDict() map[string]any {
	isEmpty := news == News{}
	if isEmpty {
		return make(map[string]any)
	}

	return map[string]any{
		"news_id":           news.ID,
		"news_title":        news.Title,
		"news_created_time": news.CreatedTime.Format("2006-01-02 15:04:05"),
		"news_status":       news.Status,
		"news_reason":       news.Reason,
	}
}
