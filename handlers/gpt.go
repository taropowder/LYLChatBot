package handlers

import (
	"LYLChatBot/constant"
	"LYLChatBot/utils"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type GptHandler struct {
}

func (h *GptHandler) Match(message *openwechat.Message) bool {

	if message.IsText() {
		if message.IsSendByGroup() {
			if message.IsAt() {
				return true
			} else {
				me, err := message.Bot().GetCurrentUser()

				if strings.Contains(message.Content, me.NickName) {
					return true
				}

				reply := utils.ParseReplyMsgText(message.Content)
				if err == nil {
					return reply.OriginalUser == me.NickName
				}

			}
		} else {
			return true
		}

	}

	return false
}

func (h *GptHandler) Helper(u *openwechat.User) string {

	return fmt.Sprintf("您可以私聊或者在群里艾特我，提出您的问题，您可与我连续对话，注意在群内的时候，我只会记得与您本人在群内的聊天记录。\n"+
		"您也可以引用聊天记录中的图片，艾特我并提出您的问题，不过只能是每位群友发的最新一张图片，并且不能是表情包等信息。"+
		"当我拒绝跟您聊天的时候，我会清空与您的所有记忆，您也可以手动输入 %s 来主动清空记忆\n", constant.ClearHistoryKey)
}

func (h *GptHandler) Name() string {
	return "GPT描述"
}

func (h *GptHandler) Handle(ctx *openwechat.MessageContext) {

	content := utils.RemoveAt(ctx.Message.Content)

	u, err := ctx.Message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return
	}

	gptFunc := utils.GetUserModuleFunc(u)

	systemPrompt := utils.GetUserSystem(u)
	scoreKey := utils.GetImgCacheKeyImageByReplyMessage(ctx.Message)

	cacheKey := utils.GetGptCacheKeyByMessage(ctx.Message)
	historyMsgs := utils.GetHistoryMessageByKey(cacheKey)

	logrus.Debugf("gpt handle %s", ctx.Message.MsgId)
	if scoreKey != "" {
		imgData := utils.GetImagesBytesByKey(scoreKey)
		replyText := utils.ParseReplyMsgText(ctx.Message.Content)
		content = utils.RemoveAt(replyText.ReplyText)
		go func() {
			resp, err := utils.NewGeminitGptWithImg(content, imgData)
			if err != nil {
				logrus.Errorf("NewGeminitGptWithImg error %s", err)
				ctx.Message.ReplyText("弔图勿扰")
			} else {
				ctx.Message.ReplyText(resp)
			}
		}()

	} else {
		if ctx.Message.IsSendByGroup() {
			groupUser, err := ctx.Message.SenderInGroup()
			if err != nil {
				logrus.Errorf("SenderInGroup error %s", err)
			} else {

				reply := utils.ParseReplyMsgText(ctx.Message.Content)
				me, err := ctx.Message.Bot().GetCurrentUser()
				if err == nil && reply.OriginalUser == me.NickName {
					content = reply.ReplyText
				}

				if utils.IsGroupMemory(u) {
					systemPrompt = utils.GetRolePrompt(constant.GroupMemoryPromptName, constant.GroupMemoryDefaultPrompt)
					name := groupUser.NickName
					if groupUser.DisplayName != "" {
						name = groupUser.DisplayName
					}
					content = fmt.Sprintf("[%s]：%s", name, content)
				}

			}
		}

		go func() {
			for i := 0; i < 3; i++ {
				resp, err := gptFunc(content, systemPrompt, historyMsgs)
				if err != nil || resp == "" {
					logrus.Errorf("gpt error %s", err)
				} else {
					utils.SetHistoryMessageByKey(cacheKey, content, resp, historyMsgs)
					ctx.Message.ReplyText(resp)
					return
				}
				time.Sleep(time.Duration(i) * time.Second)
			}

			utils.SetHistoryMessageByKey(cacheKey, content, constant.ErrorReply, historyMsgs)
			ctx.Message.ReplyText(constant.ErrorReply)

		}()

	}
}
