package handlers

import (
	"LYLChatBot/constant"
	"LYLChatBot/utils"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
)

type BiliBiliSummaryHandler struct {
}

func (h *BiliBiliSummaryHandler) Match(message *openwechat.Message) bool {

	if message.IsMedia() {
		return true
	}

	return false
}

func (h *BiliBiliSummaryHandler) Helper(u *openwechat.User) string {
	return "当您在群中或是与我的对话中存在B站视频卡片，我会帮您分析其中的内容并进行总结。"
}

func (h *BiliBiliSummaryHandler) Name() string {
	return "视频总结"
}

func (h *BiliBiliSummaryHandler) Handle(ctx *openwechat.MessageContext) {

	content := utils.RemoveAt(ctx.Message.Content)

	u, err := ctx.Message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return
	}

	systemPrompt := utils.GetRolePrompt(constant.VideoPromptName, constant.VideoDefaultPrompt)

	bilbiliBid := utils.GetBliBliBvInMessageContent(ctx.Message.Content)
	if bilbiliBid == "" {
		return
	}

	gptFunc := utils.GetUserModuleFunc(u)
	historyMsgs := make([]utils.GptMessage, 0)

	copywriting := utils.DealWithBliBli(bilbiliBid)

	if copywriting == "" {
		return
	}

	content = copywriting

	ctx.Abort()

	if ctx.Message.IsSendByGroup() {

		go func() {
			resp, err := gptFunc(content, systemPrompt, historyMsgs)
			if err != nil {
				logrus.Errorf("BiliBiliSummaryHandler error %s", err)
			} else {
				ctx.Message.ReplyText(resp)
			}
		}()

	} else {
		go func() {
			resp, err := gptFunc(content, systemPrompt, historyMsgs)
			if err != nil {
				logrus.Errorf("BiliBiliSummaryHandler error %s", err)
				ctx.Message.ReplyText("报错啦")
			} else {
				ctx.Message.ReplyText(resp)
			}
		}()

	}

	return
}
