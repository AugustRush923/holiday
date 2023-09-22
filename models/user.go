package models

import "time"

type User struct {
	BaseModel
	UserName  string    `json:"username" gorm:"size:30;not null;comment:用户名;unique"`
	Password  string    `json:"password" gorm:"size:20;not null;comment:密码"`
	NickName  string    `json:"nickname" gorm:"size:30;not null;comment:昵称"`
	Mobile    string    `json:"mobile" gorm:"size:11;comment:手机号;unique"`
	Gender    uint8     `json:"gender" gorm:"default:0;comment:性别.0-男,1-女"`
	Email     string    `json:"email" gorm:"size:40;comment:邮件;unique"`
	Avatar    string    `json:"avatar" gorm:"size:256;comment:头像"`
	LastLogin time.Time `json:"last_login" gorm:"column:last_login;type:datetime(0);default:null"`
}

func (User) TableName() string {
	return "user"
}
