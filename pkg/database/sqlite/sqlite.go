package sqlite

import (
	"LYLChatBot/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// var _db *gorm.DB

func InitDatabase(filePath string) (_db *gorm.DB) {
	var err error
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.MessageRecord{})
	db.AutoMigrate(&model.CronTaskRecord{})
	db.AutoMigrate(&model.HandlerConfigRecord{})
	db.AutoMigrate(&model.GptRoleRecord{})
	return db

}

//func GetDB() *gorm.DB {
//	return _db
//}
