package database

import "gorm.io/gorm"

var _db *gorm.DB

func GetDB() *gorm.DB {
	return _db
}

func SetDB(db *gorm.DB) {
	_db = db
}
