package models

type Category struct {
	BaseModel
	Name string
}

func (Category) TableName() string {
	return "info_category"
}

func (c Category) ToDict() map[string]any {
	isEmpty := c == Category{}
	if isEmpty {
		return map[string]any{}
	}
	return map[string]any{
		"category_id":   c.ID,
		"category_name": c.Name,
	}
}
