package models

type UserFans struct {
	FollowerID uint64
	FollowedID uint64
}

func (UserFans) TableName() string {
	return "info_user_fans"
}
