package utils

import (
	"LYLChatBot/constant"
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	"LYLChatBot/pkg/redis_conn"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gomodule/redigo/redis"
)

func GetUserSystem(u *openwechat.User) string {

	res := "你是一个非常有用的助手, 能帮我解决各种问题"

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	cacheRes, _ := redis.String(redisPool.Do("GET", fmt.Sprintf(constant.CacheSystemRoleKey, u.AvatarID())))
	if cacheRes == "" {
		return GetRolePrompt(cacheRes, res)
	}

	return res

}

func GetUserModuleFunc(u *openwechat.User) func(prompt string, systemPrompt string, historyMsgs []GptMessage) (reply string, err error) {

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	cacheRes, _ := redis.String(redisPool.Do("GET", fmt.Sprintf(constant.CacheGptModuleKey, u.AvatarID())))
	if cacheRes != "" {
		switch cacheRes {
		case constant.BingModule:
			return NewBingGpt
		case constant.QwenModule:
			return NewQwenGpt
		case constant.GeminitModule:
			return NewGeminitGpt
		}
	}
	return NewGeminitGpt

}

func GetUserModuleFuncById(uId string) func(prompt string, systemPrompt string, historyMsgs []GptMessage) (reply string, err error) {

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	cacheRes, _ := redis.String(redisPool.Do("GET", fmt.Sprintf(constant.CacheGptModuleKey, uId)))
	if cacheRes == "" {
		switch cacheRes {
		case constant.BingModule:
			return NewBingGpt
		case constant.QwenModule:
			return NewQwenGpt
		case constant.GeminitModule:
			return NewGeminitGpt
		}
	}
	return NewGeminitGpt

}

func GetRolePrompt(roleName string, defaultPrompt string) string {
	db := database.GetDB()
	role := model.GptRoleRecord{}
	err := db.Where("role_name LIKE ?", roleName).First(&role).Error
	if err != nil {
		return defaultPrompt
	} else {
		return role.Prompt
	}

}
