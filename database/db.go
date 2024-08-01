package database

import (
	"errors"
	"fmt"
	"go-dress/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

var dsn string

func getDsn() string {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Dbname)
	return dsn
}

var DB *gorm.DB
var once sync.Once

func GetDb() *gorm.DB {
	var err error
	once.Do(func() {
		dsn = getDsn()
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	})
	return DB
}
func Table(tbl string) *gorm.DB {
	return GetDb().Table(tbl)
}
func Insert(row any) {
	result := DB.Create(row)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func Update(model interface{}, updates map[string]interface{}, conditions map[string]interface{}) (int64, error) {
	if len(conditions) == 0 {
		return 0, errors.New("no conditions provided")
	}
	result := DB.Model(model).Where(conditions).Updates(updates)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func GetRow(table string, fields []string, conditions map[string]interface{}, group string, having string, order string) (map[string]interface{}, error) {
	row := DB.Table(table)
	if fields != nil {
		row = row.Select(fields)
	}
	if conditions != nil {
		row = row.Where(conditions)
	}
	if group != "" {
		row = row.Group(group)
	}
	if having != "" {
		row = row.Having(having)
	}
	if order != "" {
		row = row.Order(order)
	}
	row = row.Limit(1)
	var result map[string]interface{}
	row = row.First(&result)
	if row.Error != nil {
		return nil, row.Error
	}
	if errors.Is(row.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return result, nil
}
