package postgresql

import (
	"LYLChatBot/conf"
	"LYLChatBot/model"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// var _db *gorm.DB
const dsnFormat = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai"

func InitDatabase(host, user, password, dbname string, port int) (_db *gorm.DB) {
	dsn := fmt.Sprintf(dsnFormat, host, user, password, dbname, port)
	var err error
	_db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(conf.ConfigureInstance.LogLevel - 1)),
	})
	if err != nil {
		log.Fatal(err)
	}
	// 不打印所有sql
	_db.Logger = _db.Logger.LogMode(logger.Silent)
	_db.AutoMigrate(&model.MessageRecord{})
	_db.AutoMigrate(&model.CronTaskRecord{})
	_db.AutoMigrate(&model.HandlerConfigRecord{})
	_db.AutoMigrate(&model.GptRoleRecord{})
	_db.AutoMigrate(&model.LevatorAnusRecord{})
	_db.AutoMigrate(&model.KnowledgeRecord{})
	_db.AutoMigrate(&model.KnowledgeTypeRecord{})
	_db.AutoMigrate(&model.KnowledgeQARecord{})

	return

}

//func GetDB() *gorm.DB {
//	return _db
//}
