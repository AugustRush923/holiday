package models

import "time"

type BaseModel struct {
	ID          uint      `json:"id" gorm:"primarykey;comment:id"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime(0);autoCreateTime;comment:创建时间" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime(0);autoUpdateTime;comment:更新时间" json:"updated_time"`
	IsDeleted   uint8     `gorm:"column:is_deleted;default:0;comment:是否删除:0-否,1-是" json:"is_deleted"`
}
