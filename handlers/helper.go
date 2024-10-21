package handlers

import (
	"LYLChatBot/constant"
	"LYLChatBot/utils"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"strings"
)

type HelperHandler struct {
}

func (h *HelperHandler) Match(message *openwechat.Message) bool {

	content := utils.RemoveAt(message.Content)

	if message.IsText() && strings.HasPrefix(content, constant.HelperKey) {
		if message.IsSendByGroup() {
			if message.IsAt() {
				return true
			} else {
				return false
			}
		}

		return true
	}

	return false
}

func (h *HelperHandler) Helper(u *openwechat.User) string {
	return ""
}

func (h *HelperHandler) Name() string {
	return ""
}

func (h *HelperHandler) Handle(ctx *openwechat.MessageContext) {

	logrus.Debugf("help handle %s", ctx.Message.MsgId)

	content := utils.RemoveAt(ctx.Message.Content)

	if content == constant.HelperKey {
		res := "输入help 模块名，获取详细说明。\n当前注册模块:\n"
		for _, handler := range MessagesHandlers {
			if handler.Name() != "" {
				res = res + handler.Name() + "\n"
			}
		}
		ctx.Message.ReplyText(res)
		ctx.Abort()
	}

	helperSlice := strings.Split(content, " ")
	if len(helperSlice) == 2 {
		for _, handler := range MessagesHandlers {
			u, err := ctx.Message.Sender()
			if err != nil {
				logrus.Errorf("log error %s", err)
				return
			}

			if handler.Name() == helperSlice[1] {
				ctx.Message.ReplyText(handler.Helper(u))
				ctx.Abort()
			}
		}
	}
}
