package task

import (
	"LYLChatBot/conf"
	"LYLChatBot/constant"
	"LYLChatBot/model"
	"LYLChatBot/pkg/redis_conn"
	"LYLChatBot/utils"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type HotNewsConfig struct {
	baseTaskConfig
	HotNewsType string `json:"hot_news_type"`
}

const CacheHotNewsKey = "%s_hot_news"

func hotNews(db *gorm.DB, record model.CronTaskRecord) {
	config := HotNewsConfig{}
	err := json.Unmarshal([]byte(record.TaskConfig), &config)
	if err != nil {
		logrus.Error(err)
	}

	now := time.Now().Unix()
	if now < int64(config.NextTimeStamp) {
		return
	}

	record.Status = false
	db.Save(&record)

	defer func() {
		config.NextTimeStamp = int(now) + config.CheckInterval
		newConfigStr, err := json.Marshal(config)
		if err != nil {
			logrus.Error(err)
			return
		}
		record.TaskConfig = string(newConfigStr)
		record.Status = true

		db.Save(&record)
	}()

	hostNews, err := utils.GetHostNews(config.HotNewsType)
	if err != nil {
		logrus.Error(err)
		return
	}

	if len(hostNews.Data) > 0 {
		url := hostNews.Data[0].Url
		title := hostNews.Data[0].Title
		redisPool := redis_conn.RedisConnPool.Get()
		defer redisPool.Close()

		logrus.Debugf("hot %s news check %s", config.HotNewsType, title)
		oldTitle, _ := redis.String(redisPool.Do("GET", fmt.Sprintf(CacheHotNewsKey, config.HotNewsType)))
		if oldTitle == "" || oldTitle != title {

			htmlText, err := utils.RequestWithChrome(url)
			if err != nil {
				logrus.Errorf("CollectWechatArticle RequestWithChrome error %s", err)
				return
			}
			chinese := utils.GetHtmlText(htmlText)
			resp := ""
			if chinese != "" {
				gptFunc := utils.GetUserModuleFuncById(config.GroupId)
				if err != nil {
					logrus.Error(err)
				}
				resp, err = gptFunc(chinese, utils.GetRolePrompt(config.Role, constant.UrlDefaultPrompt), []utils.GptMessage{})
				if err != nil {
					logrus.Error(err)
				}

			}

			self, err := conf.BotInstance.GetCurrentUser()
			groups, err := self.Groups()
			findGroups := groups.SearchByID(config.GroupId)

			logrus.Infof("hot %s news update %s", config.HotNewsType, title)
			if len(findGroups) == 0 {
				logrus.Errorf("no such gid %v", config.GroupId)
			} else {
				_, err := self.SendTextToGroup(findGroups[0], fmt.Sprintf("[%s]%s热榜登顶: \n%s\n%s\n%s",
					time.Now().Format("2006-01-02 15:04:05"), config.HotNewsType,
					title, url, resp))

				//err := findGroups.SendText(, 0)
				if err != nil {
					logrus.Error(err)
				}
			}

		}

		_, err = redisPool.Do("SET", fmt.Sprintf(CacheHotNewsKey, config.HotNewsType), title)

	}

}
