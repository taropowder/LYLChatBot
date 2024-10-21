package task

import (
	"LYLChatBot/conf"
	"LYLChatBot/model"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type RemindTaskConfig struct {
	GroupId             string `json:"group_id"`
	NextRemindTimeStamp int    `json:"next_remind_time_stamp"`
	RemindContent       string `json:"remind_content"`
	// 秒数
	Interval int `json:"interval"`

	BlackStartTime int `json:"black_start_time"`
	BlackEndTime   int `json:"black_end_time"`
}

func remindRobot(db *gorm.DB, record model.CronTaskRecord) {
	config := RemindTaskConfig{}
	err := json.Unmarshal([]byte(record.TaskConfig), &config)
	if err != nil {
		logrus.Error(err)
	}

	now := time.Now().Unix()
	if now < int64(config.NextRemindTimeStamp) {
		return
	}

	record.Status = false
	db.Save(&record)

	defer func() {
		config.NextRemindTimeStamp = int(now) + config.Interval
		newConfigStr, err := json.Marshal(config)
		if err != nil {
			logrus.Error(err)
			return
		}
		record.TaskConfig = string(newConfigStr)
		record.Status = true

		db.Save(&record)
	}()

	self, err := conf.BotInstance.GetCurrentUser()
	if err != nil {
		logrus.Error(err)
		return
	}
	groups, err := self.Groups()

	findGroups := groups.SearchByID(config.GroupId)
	if len(findGroups) == 0 {
		logrus.Errorf("no such gid %v", config.GroupId)
		return
	} else {
		_, err := self.SendTextToGroup(findGroups[0], config.RemindContent)
		//err := findGroups.SendText(config.RemindContent, time.Second*5)
		if err != nil {
			logrus.Error(err)
		}
	}

}
