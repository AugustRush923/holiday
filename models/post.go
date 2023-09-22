package models

type Post struct {
	BaseModel
	Title      string `json:"title" gorm:"size:128;not null;comment:文章名称"`
	Desc       string `json:"desc" gorm:"size:512;not null;comment:文章简介"`
	Content    string `json:"content" gorm:"type:text;not null;comment:文章内容"`
	Status     uint8  `json:"status" gorm:"default:0;comment:状态.0-正常,1-下架"`
	PV         uint64 `json:"pv" gorm:"default:0;comment:pv"`
	UV         uint64 `json:"uv" gorm:"default:0;comment:uv"`
	UserId     uint64 `json:"user_id" gorm:"not null;comment:用户ID"`
	IsTop      uint8  `json:"is_top" gorm:"default:0;comment:是否指定.0-否,1-是"`
	CategoryId uint64
	Tags       []Tags `gorm:"many2many:post_tags_rel"`
}

func (Post) TableName() string {
	return "post"
}

type Category struct {
	BaseModel
	CategoryName string `gorm:"size:64;not null;comment:分类名称"`
	Status       uint8  `json:"status" gorm:"default:0;comment:状态.0-正常,1-下架"`
	IsNav        uint8  `json:"is_nav" gorm:"default:0;comment:是否为导航.0-导航,1-不导航"`
	Post         []Post `gorm:"foreignKey:CategoryId"`
}

func (Category) TableName() string {
	return "category"
}

type Tags struct {
	BaseModel
	TagName string `gorm:"size:64;nut null;comment:标签名称"`
	Status  uint8  `json:"status" gorm:"default:0;comment:状态.0-正常,1-下架"`
	Post    []Post `gorm:"many2many:post_tags_rel"`
}

func (Tags) TableName() string {
	return "tags"
}

type PostTagsRel struct {
	BaseModel
	PostId uint64
	TagsId uint64
}

func (PostTagsRel) TableName() string {
	return "post_tags_rel"
}
