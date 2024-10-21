package utils

import (
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	"LYLChatBot/pkg/redis_conn"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"regexp"
	"time"
)

const imgCacheKey = "img_cache_%s_%s"

type WechatReplyMsg struct {
	OriginalUser string
	OriginalText string
	ReplyText    string
}

func GetRecordsContentByGroupId(groupId string, rangeTime int) string {
	db := database.GetDB()

	earliestTime := time.Now().Add(time.Duration(-rangeTime) * time.Second)

	res := ""
	records := []model.MessageRecord{}
	err := db.Where("(group_id = ?  or user_id = ?) and updated_at > ?", groupId, groupId, earliestTime).Find(&records).Error
	if err != nil {
		return ""
	} else {
		for _, record := range records {
			res = res + fmt.Sprintf("time : %v 发言人: %v(%v) 发言内容 : '%v'  \n", record.CreatedAt,
				record.NickName, record.UserName, record.Content)
		}
	}

	return res
}

func GetImgCacheKeyImageByReplyMessage(message *openwechat.Message) string {
	u, err := message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return ""
	}
	scoreKey := ""
	replyText := ParseReplyMsgText(message.Content)

	if replyText.OriginalUser != "" && replyText.OriginalText == "[图片]" {
		if group, success := u.AsGroup(); success {
			scoreKey = fmt.Sprintf(imgCacheKey, group.AvatarID(), replyText.OriginalUser)
		} else {
			scoreKey = fmt.Sprintf(imgCacheKey, u.AvatarID(), u.AvatarID())

		}
	}

	return scoreKey
}

func GetImgCacheKeyImageByPicture(message *openwechat.Message) string {
	u, err := message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return ""
	}
	scoreKey := ""

	if group, success := u.AsGroup(); success {
		groupUser, err := message.SenderInGroup()
		if err != nil {
			logrus.Errorf("SenderInGroup error %s", err)
		} else {
			name := groupUser.DisplayName
			if name == "" {
				name = groupUser.NickName
			}
			scoreKey = fmt.Sprintf(imgCacheKey, group.AvatarID(), name)

		}
	} else {
		scoreKey = fmt.Sprintf(imgCacheKey, u.AvatarID(), u.AvatarID())

	}

	return scoreKey
}

func GetImagesBytesByKey(scoreKey string) []byte {

	if scoreKey != "" {

		redisPool := redis_conn.RedisConnPool.Get()
		defer redisPool.Close()

		res, err := redis.Bytes(redisPool.Do("GET", scoreKey))
		if err != nil {
			logrus.Debug("获取 value 失败:", err)
			return nil
		} else {
			return res
		}

	}

	return nil

}

func ParseReplyMsgText(text string) (msg WechatReplyMsg) {
	// 定义一个正则表达式
	pattern := "(?s)「(.+?)：(.+)」\\n- - - - - - - - - - - - - - -\\n(.*)"

	// 编译正则表达式
	reg := regexp.MustCompile(pattern)

	matches := reg.FindStringSubmatch(text)

	if len(matches) == 4 {
		msg.OriginalUser = matches[1]
		msg.OriginalText = matches[2]
		msg.ReplyText = matches[3]
	}

	return
}
