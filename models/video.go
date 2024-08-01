package models

import (
	"go-dress/database"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Video struct {
	ID                uint   `gorm:"column:id,primary_key;AUTO_INCREMENT"`
	Name              string `gorm:"column:name"`
	Cover             string `gorm:"column:cover"`
	ProducerId        uint   `gorm:"column:producer_id"`
	ProducerName      string `gorm:"column:producer_name"`
	ProducerAvatar    string `gorm:"column:producer_avatar"`
	Synopsis          string `gorm:"column:synopsis"`
	Cast              string `gorm:"column:cast;type:json;"`
	DramaId           uint   `gorm:"column:drama_id"`
	DramaName         string `gorm:"column:drama_name"`
	DramaRoyalty      string `gorm:"column:drama_royalty;default:0"`
	PlanId            uint   `gorm:"column:plan_id"`
	StudioId          uint   `gorm:"column:studio_id"`
	StudioName        string `gorm:"column:studio_name"`
	StudioMinutePrice uint   `gorm:"column:studio_minute_price"`
	ShootDate         string `gorm:"column:shoot_date;type:data;"`
	ShootStart        int64  `gorm:"column:shoot_start"`
	ShootEnd          int64  `gorm:"column:shoot_end"`
	Address           string `gorm:"column:address"`
	Cost              uint   `gorm:"column:cost"`
	Price             uint   `gorm:"column:price"`
	Royalty           uint   `gorm:"column:royalty"`
	RemakeNum         uint   `gorm:"column:remake_num"`
	PublishTime       int64  `gorm:"column:publish_time"`
	State             uint8  `gorm:"column:state"`
}

func AddVideo(row *Video) uint {
	row.PublishTime = time.Now().Unix()
	err := database.DB.Table("video").Create(&row).Error
	if err != nil {
		panic(err)
	}
	return row.ID
}

func SetRoleAvatar(videoId uint, roleNo uint8, avatar string) error {
	return database.DB.Table("video").Where("id=?", videoId).Update("cast", gorm.Expr("JSON_SET(`cast`, '$["+strconv.Itoa(int(roleNo))+"].avatar', ?)", avatar)).Error

}
