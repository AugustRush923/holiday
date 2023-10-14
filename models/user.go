package models

import (
	"time"
)

type User struct {
	BaseModel
	NickName     string    `json:"nickname" gorm:"column:nick_name;comment:昵称;"`
	PasswordHash string    `json:"password" gorm:"column:password_hash;comment:密码"`
	Mobile       string    `json:"mobile" gorm:"size:11;comment:手机号;unique"`
	AvatarUrl    string    `json:"avatar" gorm:"column:avatar_url;size:256;comment:头像"`
	LastLogin    time.Time `json:"last_login" gorm:"column:last_login;type:datetime(0);default:null"`
	IsAdmin      uint8     `json:"is_admin" gorm:"column:is_admin;default:0;comment:是否为管理页0-否,1-是"`
	Signature    string    `json:"signature" gorm:"size:512;comment:用户签名"`
	Gender       string    `json:"gender" gorm:"comment:性别:MAN WOMAN"`
}

func (User) TableName() string {
	return "info_user"
}

func (user User) Dict() (userDict map[string]any) {
	if user.IsAdmin == 1 {
		userDict = map[string]any{
			"id":         user.ID,
			"nickname":   user.NickName,
			"mobile":     user.Mobile,
			"avatar_url": user.AvatarUrl,
			"last_login": user.LastLogin,
			"gender":     user.Gender,
		}
	} else {
		userDict = map[string]any{
			"id":              user.ID,
			"nickname":        user.NickName,
			"mobile":          user.Mobile,
			"gender":          user.Gender,
			"signature":       user.Signature,
			"followers_count": 0,
			"news_count":      0,
		}
	}
	return userDict
}

func (user User) CheckPasswd(passwd string) bool {

	return true
}
