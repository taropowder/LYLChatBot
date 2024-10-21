package task

import (
	"LYLChatBot/conf"
	"LYLChatBot/model"
	"LYLChatBot/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ShitStirrerConfig struct {
	GroupId            string `json:"group_id"`
	NextCheckTimeStamp int    `json:"next_check_time_stamp"`
	ShitStirrerRole    string `json:"shit_stirrer_role"`

	// 5分钟超过10条消息则触发

	ChatHistoryTime int `json:"chat_history_time"`

	ChatHistoryMaxSize int `json:"chat_history_max_size"`

	// 间隔秒数
	CheckInterval int `json:"check_interval"`
	SpeakInterval int `json:"speak_interval"`
}

// 触发的条件如下：
// 1. n 分钟聊天记录超过 m 条
// 2. 当前时间大于下次触发时间
//

func ShitStirrer(db *gorm.DB, record model.CronTaskRecord) {
	config := ShitStirrerConfig{}
	err := json.Unmarshal([]byte(record.TaskConfig), &config)
	if err != nil {
		logrus.Error(err)
	}

	now := time.Now().Unix()
	if now < int64(config.NextCheckTimeStamp) {
		return
	}

	record.Status = false
	db.Save(&record)

	speak := false

	defer func() {

		config.NextCheckTimeStamp = int(now) + config.CheckInterval
		if speak {
			t := int(now) + config.SpeakInterval
			if t > config.NextCheckTimeStamp {
				config.NextCheckTimeStamp = t
			}
		}

		newConfigStr, err := json.Marshal(config)
		if err != nil {
			logrus.Error(err)
			return
		}
		record.TaskConfig = string(newConfigStr)
		record.Status = true

		db.Save(&record)
	}()

	//	获取聊天记录
	recordHistory := utils.GetRecordsContentByGroupId(config.GroupId, config.ChatHistoryTime)
	if recordHistory != "" && strings.Count(recordHistory, "发言人")+1 > config.ChatHistoryMaxSize {
		logrus.Infof("start ShitStirrer for %v", config.GroupId)
		//	 开始搅屎
		systemPrompt := utils.GetRolePrompt(config.ShitStirrerRole, "现在请你扮演一个在社交群聊中热爱聊天的人，"+
			"我将会给你发送一段聊天记录，请你根据这些聊天记录上下文,发表跟上下文相关的一段话，你可以发表对这段话的评论，或者抛出一个新的问题，"+
			"无论如何，请你直接开始你的评论和问题就可以，请千万不要让别人看出来你是一个AI，回答中也不要提到聊天记录相关的内容，")
		gptFunc := utils.GetUserModuleFuncById(config.GroupId)

		historyMsgs := make([]utils.GptMessage, 0)
		resp, err := gptFunc(recordHistory, systemPrompt, historyMsgs)
		if err == nil {
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
				speak = true
				_, err := self.SendTextToGroup(findGroups[0], resp)
				//err := findGroups.SendText(resp, time.Second*5)
				if err != nil {
					logrus.Error(err)
				}
			}
		}
	}

}
