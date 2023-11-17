package models

import "time"

type UserCollection struct {
	UserID      uint64
	NewsID      uint64
	CreatedTime time.Time `gorm:"column:create_time;type:datetime(0);autoCreateTime;comment:创建时间" json:"created_time"`
}

func (UserCollection) TableName() string {
	return "info_user_collection"
}
