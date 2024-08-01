package models

import (
	"go-dress/database"
)

type Image struct {
	Id          uint   `gorm:"column:id;type:bigint unsigned;primary_key;AUTO_INCREMENT"`
	ImageName   string `gorm:"column:image_name;type:varchar(255);not null"`
	OwnerUid    uint   `gorm:"column:owner_uid;type:bigint unsigned;not null"`
	Category    uint   `gorm:"column:category;type:tinyint unsigned;not null"`
	ExtName     string `gorm:"column:ext_name;type:varchar(255);not null"`
	Size        uint   `gorm:"column:size;type:int unsigned;not null;default:0"`
	Url         string `gorm:"column:url;type:varchar(255);"`
	OperatorUid uint   `gorm:"column:operator_uid;type:bigint unsigned;not null"`
}

func CreateImageId(image Image) uint {
	row := Image{ImageName: image.ImageName, OwnerUid: image.OwnerUid, Category: image.Category, ExtName: image.ExtName, Size: image.Size, OperatorUid: image.OperatorUid}
	db := database.GetDb()
	db.Table("image").Create(&row)
	return row.Id
}
