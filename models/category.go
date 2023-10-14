package models

type Category struct {
	BaseModel
	Name string
}

func (Category) TableName() string {
	return "info_category"
}
