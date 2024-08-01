package models

import (
	"go-dress/database"
)

const IMAGE_BASE_URL = "https://image.han-dress.cn/"

type User struct {
	ID       uint   `gorm:"column:id;primary_key;auto_increment"`
	Avatar   string `gorm:"column:avatar"`
	AvatarID uint   `gorm:"column:avatar_id"`
}

func UpdateAvatarByUid(uid uint, avatar string, avatarId uint) error {
	return database.DB.Table("user").Where("id = ?", uid).Update("avatar", IMAGE_BASE_URL+avatar).Update("avatar_id", avatarId).Error
}

func UpdateUserById(uid uint, info map[string]any) bool {
	res := database.Table("user").Where("id = ?", uid).Updates(info)
	if res.Error != nil {
		panic(res.Error)
	}
	return true
}
