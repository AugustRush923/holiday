package models

import (
	"holiday/dao"
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

func (user User) ToDict() (userDict map[string]any) {
	isEmpty := user == User{}
	if isEmpty {
		return make(map[string]any)
	}
	var followersCount int64
	dao.DB.Where("followed_id = ?", user.ID).Find(&UserFans{}).Count(&followersCount)
	var newsCount int64
	dao.DB.Where("user_id = ?", user.ID).Find(&News{}).Count(&newsCount)
	userDict = map[string]any{
		"id":              user.ID,
		"nickname":        user.NickName,
		"mobile":          user.Mobile,
		"gender":          user.Gender,
		"signature":       user.Signature,
		"followers_count": followersCount,
		"news_count":      newsCount,
	}
	return userDict
}

func (user User) ToAdminDict() (userDict map[string]any) {
	isEmpty := user == User{}
	if isEmpty {
		return make(map[string]any)
	}
	userDict = map[string]any{
		"id":         user.ID,
		"nickname":   user.NickName,
		"mobile":     user.Mobile,
		"avatar_url": user.AvatarUrl,
		"last_login": user.LastLogin,
		"gender":     user.Gender,
	}
	return userDict
}
