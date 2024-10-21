package config

import (
	"LYLChatBot/constant"
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	"LYLChatBot/pkg/redis_conn"
	"LYLChatBot/utils"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type GptConfigHandler struct {
}

func (h *GptConfigHandler) Match(message *openwechat.Message) bool {

	content := utils.RemoveAt(message.Content)

	// 使用正则表达式匹配 PlayerRoleMenu
	if utils.MessageMatchInstruct(constant.PlayerRoleMenu, message) != nil {
		return true
	}

	if utils.MessageMatchInstruct(constant.GptModuleMenu, message) != nil {
		return true
	}

	if content == constant.GptGroupOneselfMemory || content == constant.GptGroupMemory {
		return true
	}

	if content == constant.KnowledgeStart || content == constant.KnowledgeEnd {
		return true
	}

	return false
}

func (h *GptConfigHandler) Helper(u *openwechat.User) string {
	db := database.GetDB()
	roles := make([]model.GptRoleRecord, 0)
	err := db.Find(&roles).Error
	if err != nil {
		logrus.Error(err)
		return "啊，有些乏了，可能是哪里出了问题"
	} else {
		rolesStr := ""
		for _, role := range roles {
			rolesStr = rolesStr + role.RoleName + "/"
		}

		redisPool := redis_conn.RedisConnPool.Get()
		defer redisPool.Close()

		module, _ := redis.String(redisPool.Do("GET", fmt.Sprintf(constant.CacheGptModuleKey, u.AvatarID())))
		if module == "" {
			module = "geminit"
		}

		role, _ := redis.String(redisPool.Do("GET", fmt.Sprintf(constant.CacheSystemRoleKey, u.AvatarID())))
		if role == "" {
			role = "普通模式"
		}

		return fmt.Sprintf("您可输入以下指令:\n扮演角色:%s\n设置GPT模块:bing/qwen/geminit\n当前GPT模块:%s ; 角色:%s ;\n 是否为全群记忆模式:%v\n", rolesStr, module, role, utils.IsGroupMemory(u))
	}
}

func (h *GptConfigHandler) Name() string {
	return "GPT配置"
}

func (h *GptConfigHandler) Handle(ctx *openwechat.MessageContext) {

	logrus.Debugf("gpt config handle %s", ctx.Message.MsgId)

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	u, err := ctx.Message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return
	}

	content := utils.RemoveAt(ctx.Content)

	if matches := utils.MessageMatchInstruct(constant.PlayerRoleMenu, ctx.Message); matches != nil && len(matches) > 0 {
		_, err = redisPool.Do("SET", fmt.Sprintf(constant.CacheSystemRoleKey, u.AvatarID()), matches[1])
		if err == nil {
			_, err := ctx.Message.ReplyText(fmt.Sprintf("更新 system role 为 %s", matches[1]))
			if err != nil {
				return
			}
		}

	}

	if matches := utils.MessageMatchInstruct(constant.GptModuleMenu, ctx.Message); matches != nil && len(matches) > 0 {
		_, err = redisPool.Do("SET", fmt.Sprintf(constant.CacheGptModuleKey, u.AvatarID()), matches[1])
		if err == nil {
			_, err := ctx.Message.ReplyText(fmt.Sprintf("更新 model 为 %s", matches[1]))
			if err != nil {
				return
			}
		}
	}

	if content == constant.GptGroupMemory {
		_, err = redisPool.Do("SET", fmt.Sprintf(constant.CacheGroupMemoryKey, u.AvatarID()), true)
		if err == nil {
			_, err := ctx.Message.ReplyText("已进入全群记忆模式")
			if err != nil {
				return
			}
		}
	}

	if content == constant.GptGroupOneselfMemory {
		_, err = redisPool.Do("SET", fmt.Sprintf(constant.CacheGroupMemoryKey, u.AvatarID()), false)
		if err == nil {
			_, err := ctx.Message.ReplyText("已进入独立记忆模式")
			if err != nil {
				return
			}
		}
	}

	switch content {
	case constant.KnowledgeStart:
		_, err = redisPool.Do("SET", fmt.Sprintf(constant.CacheKnowledgeKey, u.AvatarID()), true)
		if err == nil {
			_, err := ctx.Message.ReplyText("已开启知识库")
			if err != nil {
				return
			}
		}
	case constant.KnowledgeEnd:
		_, err = redisPool.Do("DEL", fmt.Sprintf(constant.CacheKnowledgeKey, u.AvatarID()))
		if err == nil {
			_, err := ctx.Message.ReplyText("已关闭知识库")
			if err != nil {
				return
			}
		}
	}

	ctx.Abort()

	//if ctx.Message.IsText() {
	//
	//	content := ctx.Message.Content
	//

	//
	//	if ctx.IsSendByGroup() {
	//		if !ctx.Message.IsAt() {
	//			return
	//		} else {
	//			content = utils.RemoveAt(content)
	//		}
	//	}
	//
	//	// update_system_role :  狂暴模式
	//	// xxxxx
	//	// play_role : 狂暴模式
	//	// list_system_role :
	//	// desc_system_role :
	//
	//	result := strings.Split(content, ":")
	//
	//	if len(result) == 2 {
	//
	//		instructions := strings.TrimSpace(result[0])
	//
	//		switch instructions {
	//		case "update_system_role":
	//			systemRoleStrSlice := strings.Split(result[1], "\n")
	//			if len(systemRoleStrSlice) == 2 {
	//				role := model.GptRoleRecord{}
	//				err = db.Where("role_name = ?", systemRoleStrSlice[0]).First(&role).Error
	//				if errors.Is(err, gorm.ErrRecordNotFound) {
	//					role.RoleName = strings.TrimSpace(systemRoleStrSlice[0])
	//					role.Prompt = systemRoleStrSlice[1]
	//				} else {
	//					role.Prompt = systemRoleStrSlice[1]
	//				}
	//				db.Save(&role)
	//				_, err := ctx.Message.ReplyText(fmt.Sprintf("%s 更新成功", role.RoleName))
	//				if err != nil {
	//					return
	//				}
	//
	//			}
	//		case "list_system_role":
	//			roles := make([]model.GptRoleRecord, 0)
	//			db.Find(&roles)
	//			roleStr := ""
	//			for _, role := range roles {
	//				roleStr = roleStr + role.RoleName + "\n"
	//			}
	//			_, err := ctx.Message.ReplyText(roleStr)
	//			if err != nil {
	//				return
	//			}
	//		case "desc_system_role":
	//
	//			role := model.GptRoleRecord{}
	//			err = db.Where("role_name = ?", strings.TrimSpace(strings.TrimSpace(result[1]))).First(&role).Error
	//			if errors.Is(err, gorm.ErrRecordNotFound) {
	//				_, err := ctx.Message.ReplyText("no such system role")
	//				if err != nil {
	//					return
	//				}
	//			} else {
	//				_, err := ctx.Message.ReplyText(role.Prompt)
	//				if err != nil {
	//					return
	//				}
	//			}
	//		case "play_role":
	//
	//			roleName := strings.TrimSpace(result[1])
	//			roleConfigRecord := model.HandlerConfigRecord{}
	//
	//			var tx *gorm.DB
	//			if u.AvatarID() != "" {
	//				tx = db.Where("user_id LIKE ?", fmt.Sprintf("%%%s%%", u.AvatarID()))
	//			}
	//			if u.NickName != "" {
	//				tx = tx.Or("user_name LIKE ?", fmt.Sprintf("%%%s%%", u.NickName))
	//			}
	//
	//			err = tx.First(&roleConfigRecord).Error
	//			if errors.Is(err, gorm.ErrRecordNotFound) {
	//				c := model.HandlerConfig{}
	//				c.SystemRole = roleName
	//				roleConfigRecord.SetConfig(c)
	//				roleConfigRecord.UserId = u.AvatarID()
	//				roleConfigRecord.UserName = u.NickName
	//			} else {
	//				c := roleConfigRecord.GetConfig()
	//				c.SystemRole = roleName
	//				roleConfigRecord.SetConfig(c)
	//			}
	//			db.Save(&roleConfigRecord)
	//			_, err := ctx.Message.ReplyText(fmt.Sprintf("已经进入: %s", roleName))
	//			if err != nil {
	//				return
	//			}
	//		case "set_prefix_word":
	//
	//			prefix_word := result[1]
	//
	//			redisPool := redis_conn.RedisConnPool.Get()
	//			defer redisPool.Close()
	//
	//			_, err := redisPool.Do("SET", fmt.Sprintf("%s_prefix_word", u.AvatarID()), strings.TrimSpace(prefix_word))
	//			if err != nil {
	//				logrus.Error(err)
	//			}
	//
	//			ctx.Message.ReplyText(fmt.Sprintf("prefix_word set Successful: %s", prefix_word))
	//			ctx.Message.MsgId = ""
	//
	//		case "remove_prefix_word":
	//
	//			redisPool := redis_conn.RedisConnPool.Get()
	//			defer redisPool.Close()
	//
	//			_, err := redisPool.Do("DEL", fmt.Sprintf("%s_prefix_word", u.AvatarID()))
	//			if err != nil {
	//				logrus.Error(err)
	//			}
	//
	//			ctx.Message.ReplyText(fmt.Sprintf("remove_prefix_word set Successful"))
	//			ctx.Message.MsgId = ""
	//
	//		case "set_config_yaml":
	//
	//			configStrSlice := strings.Split(result[1], " ")
	//
	//			configName := strings.TrimSpace(configStrSlice[1])
	//			configValue := strings.TrimSpace(configStrSlice[2])
	//
	//			switch configName {
	//			case "bilibili_cookie":
	//				conf.ConfigureInstance.Handlers.Abstract.Cookie = configValue
	//				ctx.Message.ReplyText(fmt.Sprintf("bilibili_cookie set success"))
	//
	//			}
	//
	//		}
	//
	//		ctx.Abort()
	//
	//	}
	//
	//}

}
