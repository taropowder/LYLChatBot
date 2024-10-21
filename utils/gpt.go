package utils

import (
	"LYLChatBot/constant"
	"LYLChatBot/pkg/redis_conn"
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"strings"
)

const gptCacheKey = "gpt_cache_%s_%s"

type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func GetHistoryMessageByKey(scoreKey string) []GptMessage {

	msgs := make([]GptMessage, 0)

	if scoreKey != "" {

		redisPool := redis_conn.RedisConnPool.Get()
		defer redisPool.Close()

		res, err := redis.String(redisPool.Do("GET", scoreKey))
		if err != nil {
			logrus.Debug("获取 value 失败:", err)
			return msgs
		} else {
			err := json.Unmarshal([]byte(res), &msgs)
			if err != nil {
				logrus.Error(err)
			} else {
				return msgs
			}
		}

		return msgs
	}

	return msgs

}

func GetGptCacheKeyByMessage(message *openwechat.Message) string {
	u, err := message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return ""
	}

	scoreKey := ""
	if group, success := u.AsGroup(); success {

		if IsGroupMemory(u) {
			scoreKey = fmt.Sprintf(gptCacheKey, u.AvatarID(), u.AvatarID())
		} else {
			groupUser, err := message.SenderInGroup()
			if err != nil {
				logrus.Errorf("SenderInGroup error %s", err)
			} else {
				scoreKey = fmt.Sprintf(gptCacheKey, group.AvatarID(), groupUser.NickName)

			}
		}

	} else {
		scoreKey = fmt.Sprintf(gptCacheKey, u.AvatarID(), u.AvatarID())

	}
	return scoreKey
}

func IsGroupMemory(u *openwechat.User) bool {
	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	res, err := redis.Bool(redisPool.Do("GET", fmt.Sprintf(constant.CacheGroupMemoryKey, u.AvatarID())))
	if err != nil {
		logrus.Debug("获取 value 失败:", err)
		return false
	} else {
		return res
	}
}

func SetHistoryMessageByKey(cacheKey, userText, respText string, historyMsgs []GptMessage) {

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	if respText == constant.ErrorReply {
		logrus.Errorf("stop words %v", respText)
		_, err := redisPool.Do("DEL", cacheKey)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	if strings.Contains(userText, constant.ClearHistoryKey) {
		logrus.Errorf("stop words %v", respText)
		_, err := redisPool.Do("DEL", cacheKey)
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	historyMsgs = append(historyMsgs, GptMessage{
		Role:    "user",
		Content: userText,
	})
	historyMsgs = append(historyMsgs, GptMessage{
		Role:    "assistant",
		Content: respText,
	})
	cacheValue, err := json.Marshal(historyMsgs)
	if err != nil {
		logrus.Error(err)
	} else {
		_, err = redisPool.Do("SET", cacheKey, string(cacheValue))
	}
}

func IsGlobalModel(uid, prefixWord string) bool {

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	res, err := redis.String(redisPool.Do("GET", fmt.Sprintf("%s_prefix_word", uid)))
	if err != nil {
		logrus.Debug("获取 value 失败:", err)
		return false
	} else {
		if prefixWord == res {
			return true
		}
	}

	return false
}
