package models

import (
	"go-dress/database"
)

type Tags struct {
	ID         uint   `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name       string `gorm:"column:tag;"`
	CreatorUID uint   `gorm:"column:creator_uid;"`
}

func AddTag(tag Tags) uint {
	database.Insert(&tag)
	return tag.ID
}
