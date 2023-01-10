package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"jam3.com/user/config"
)

var _db *gorm.DB

func DB() *gorm.DB {

	return _db
}
func init() {

	dsn := config.C.DB.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("connect db fail, error=" + err.Error())
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	_db = db
}
