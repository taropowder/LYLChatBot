package handlers

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
	"regexp"
	"strconv"
	"time"
)

type KnowledgeHandler struct {
}

func (h *KnowledgeHandler) Match(message *openwechat.Message) bool {

	u, err := message.Sender()
	if err != nil {
		logrus.Error(err)
		return false
	}

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	_, err = redis.Bool(redisPool.Do("GET", fmt.Sprintf(constant.CacheKnowledgeKey, u.AvatarID())))
	if err != nil {
		return false
	}

	if message.IsText() {

		me, err := message.Bot().GetCurrentUser()
		if err != nil {
			logrus.Error(err)
			return false
		}

		content := message.Content
		reply := utils.ParseReplyMsgText(content)
		if reply.OriginalText == "" || reply.OriginalUser != me.NickName {
			return false
		}

		regex := regexp.MustCompile(constant.KnowledgeKey)

		matches := regex.FindStringSubmatch(reply.OriginalText)

		if matches != nil {
			return true
		}

	}

	return false
}

func (h *KnowledgeHandler) Helper(u *openwechat.User) string {

	return fmt.Sprintf("记录知识库")
}

func (h *KnowledgeHandler) Name() string {
	return "知识库"
}

func (h *KnowledgeHandler) Handle(ctx *openwechat.MessageContext) {

	content := utils.RemoveAt(ctx.Message.Content)

	u, err := ctx.Message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return
	}

	gptFunc := utils.GetUserModuleFunc(u)

	systemPrompt := utils.GetUserSystem(u)

	reply := utils.ParseReplyMsgText(content)

	regex := regexp.MustCompile(constant.KnowledgeKey)

	content = reply.ReplyText

	matches := regex.FindStringSubmatch(reply.OriginalText)

	kId := 0

	historyMsgs := make([]utils.GptMessage, 0)

	if matches != nil {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			kId = num
		}
		historyMsgs = utils.GetHistoryByKnowledgeId(kId)
	}

	if len(historyMsgs) == 0 {
		return
	}

	ctx.Abort()

	// ----------------
	go func() {
		for i := 0; i < 3; i++ {
			resp, err := gptFunc(content, systemPrompt, historyMsgs)
			if err != nil || resp == "" {
				logrus.Errorf("gpt error %s", err)
			} else {

				db := database.GetDB()

				sharer := ""

				if ctx.Message.IsSendByGroup() {
					groupUser, err := ctx.Message.SenderInGroup()
					if err != nil {
						logrus.Errorf("SenderInGroup error %s", err)
					} else {
						if groupUser.DisplayName != "" {
							sharer = groupUser.DisplayName
						} else {
							sharer = groupUser.NickName
						}
					}
				} else {
					if u.DisplayName != "" {
						sharer = u.DisplayName
					} else {
						sharer = u.NickName
					}
				}

				k := model.KnowledgeRecord{}
				db.Where("id = ?", kId).First(&k)
				qa := model.KnowledgeQARecord{
					Knowledge:  k,
					Question:   content,
					Answer:     resp,
					Questioner: sharer,
					IsFirst:    false,
				}
				db.Save(&qa)
				resp = fmt.Sprintf(constant.KnowledgeKeyFormat, kId) + resp
				ctx.Message.ReplyText(resp)
				return
			}
			time.Sleep(time.Duration(i) * time.Second)
		}

		ctx.Message.ReplyText(constant.ErrorReply)

	}()

}
