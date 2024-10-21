package task

import (
	"LYLChatBot/conf"
	"LYLChatBot/model"
	"LYLChatBot/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strings"
	"time"
)

type ScoreGGLiveTextTaskConfig struct {
	GroupId      []string `json:"group_id"`
	EndTimestamp int      `json:"end_timestamp"`
	Interval     int      `json:"interval"`
	LastKey      string   `json:"last_key"`
	MatchId      string   `json:"match_id"`
}

func reportLiveTextTask(db *gorm.DB, record model.CronTaskRecord) {

	config := ScoreGGLiveTextTaskConfig{}
	err := json.Unmarshal([]byte(record.TaskConfig), &config)
	if err != nil {
		logrus.Error(err)
	}
	now := time.Now().Unix()
	if now > int64(config.EndTimestamp) {
		record.Status = false
		db.Save(&record)
		return
	}

	lastTime := now - record.UpdatedAt.Unix()
	gap := time.Now().Unix() - lastTime

	if gap > int64(config.Interval) {
		messages, lastKey := utils.GetLastLiveTextMessageByMatchIdAndLastKey(config.MatchId, config.LastKey)

		if lastKey == "" {
			new_record := model.CronTaskRecord{}
			err = db.Where("id = ?", record.ID).Find(&new_record).Error
			if !errors.Is(err, gorm.ErrRecordNotFound) && new_record.TaskConfig == record.TaskConfig {
				record.Status = true
				db.Save(&record)
			}
			return
		}

		self, err := conf.BotInstance.GetCurrentUser()
		if err != nil {
			logrus.Error(err)
			return
		}
		groups, err := self.Groups()

		reportGroups := make(openwechat.Groups, 0)

		for _, gId := range config.GroupId {

			findGroups := groups.SearchByID(gId)
			if len(findGroups) == 0 {
				logrus.Errorf("no such gid %v", gId)
			} else {
				reportGroups = append(reportGroups, groups.SearchByID(gId)...)

			}
		}

		if lastKey == "error" {
			record.Status = false
			db.Save(&record)
			err := reportGroups.SendText(fmt.Sprintf("match %s 有误，请确认是否正确", config.MatchId), time.Second*5)
			if err != nil {
				logrus.Error(err)
			}
			return
		}

		config.LastKey = lastKey
		newConfigStr, err := json.Marshal(config)
		if err != nil {
			logrus.Error(err)
			return
		}

		record.TaskConfig = string(newConfigStr)
		record.Status = false
		db.Save(&record)

		for _, message := range messages {

			if reportGroups.Count() > 0 {

				logrus.Debugf("get message %v", message)

				if !strings.HasPrefix(message, "http") && !strings.Contains(message, "分享图片") {

					for _, rp := range reportGroups {
						_, err := self.SendTextToGroup(rp, message)
						if err != nil {
							logrus.Error(err)
							continue
						}
						time.Sleep(3 * time.Second)

					}

					//err := reportGroups.SendText(message, time.Second*5)
					//if err != nil {
					//	logrus.Error(err)
					//	continue
					//}
				} else {
					client := &http.Client{
						Timeout: 30 * time.Second,
					}

					// 下载文件
					resp, err := client.Get(message)
					if err != nil {
						logrus.Error(err)
						continue
					}
					defer resp.Body.Close()

					// 将 Body 属性转换为 io.Reader 类型
					reader := io.Reader(resp.Body)

					for _, rp := range reportGroups {
						self.SendImageToGroup(rp, reader)
						if err != nil {
							logrus.Error(err)
							continue
						}
						time.Sleep(3 * time.Second)
					}
					// TODO: FIX ME
					//err = reportGroups.SendImage(reader, time.Second*5)

				}

			}
			time.Sleep(time.Second * 10)

		}

		new_record := model.CronTaskRecord{}
		err = db.Where("id = ?", record.ID).Find(&new_record).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) && new_record.TaskConfig == record.TaskConfig {
			record.Status = true
			db.Save(&record)
		}

	} else {
		logrus.Debugf("time gap %v", gap)
	}

}
