package models

import (
	"go-dress/database"
	"time"
)

type Cast struct {
	ID       uint   `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	PlanId   uint   `gorm:"column:plan_id;default:0"`
	RoleNo   uint8  `gorm:"column:role_no;default:0"`
	VideoId  uint   `gorm:"column:video_id;default:0"`
	ActorUid uint   `gorm:"column:actor_uid;default:0"`
	RoleId   uint   `gorm:"column:role_id;default:0"`
	RoleName string `gorm:"column:role_name;default:''"`
	Avatar   string `gorm:"column:avatar;default:'';"`
	Profile  string `gorm:"column:profile;default:'';"`
	Fee      uint   `gorm:"column:fee;default:0"`
	SignTime int64  `gorm:"column:sign_time;default:0"`
}

func AddCast(row *Cast) uint {
	row.SignTime = time.Now().Unix()
	err := database.DB.Table("cast").Create(&row).Error
	if err != nil {
		panic(err)
	}
	return row.ID
}

func GetCastInfoById(id uint, fields []string) Cast {
	var result Cast
	err := database.DB.Table("cast").Where("id=?", id).Select(fields).First(&result).Error
	if err != nil {
		panic(err)
	}
	return result
}

func UpdateCastById(id uint, row map[string]any) error {
	return database.DB.Table("cast").Where("id=?", id).Updates(row).Error
}
