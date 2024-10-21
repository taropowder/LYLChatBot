package handlers

import (
	"LYLChatBot/constant"
	"LYLChatBot/utils"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

type ChatSummaryHandler struct {
}

func (h *ChatSummaryHandler) Match(message *openwechat.Message) bool {

	if strings.Contains(message.Content, constant.ChatSummaryKey) {
		return true

	}

	return false
}

func (h *ChatSummaryHandler) Helper(u *openwechat.User) string {
	return "您可使用 '@我 结合历史聊天 指令'的方式，让我结合群聊中的上下文来回答您的问题，比如:\n@我 群友过去1小时的聊天总结 或是\n @我 刚刚xxx群友说的是什么(不加时间描述默认查询1小时内记录)"
}

func (h *ChatSummaryHandler) Name() string {
	return "结合历史聊天"
}
func (h *ChatSummaryHandler) Handle(ctx *openwechat.MessageContext) {

	logrus.Debugf("ChatSummary handle %s", ctx.Message.MsgId)

	content := utils.RemoveAt(ctx.Message.Content)

	u, err := ctx.Message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return
	}

	systemPrompt := utils.GetRolePrompt(constant.ChatSummaryPromptName, constant.ChatSummaryDefaultPrompt)

	gptFunc := utils.GetUserModuleFunc(u)

	historyMsgs := make([]utils.GptMessage, 0)

	re := regexp.MustCompile(`(\d+)小时`)

	hours := 1
	match := re.FindStringSubmatch(ctx.Message.Content)
	if match != nil {
		fmt.Println("匹配到的数字: ", match[1])
		num, err := strconv.Atoi(match[1])
		if err == nil {
			hours = num
		}
	}

	content = "聊天记录如下:\n------" + utils.GetRecordsContentByGroupId(u.AvatarID(), 60*60*hours) +
		"-------\n 我的指令如下" + content

	if ctx.Message.IsSendByGroup() {

		if ctx.Message.IsAt() {
			go func() {
				resp, err := gptFunc(content, systemPrompt, historyMsgs)
				if err != nil {
					ctx.Message.ReplyText("不好意思，你们说的太乱了，给我CPU干烧了，恕在下无能为力")
					logrus.Errorf("ChatSummaryHandler error %s", err)
				} else {
					ctx.Message.ReplyText(resp)
				}
			}()
			ctx.Abort()
		}

	} else {
		go func() {
			resp, err := gptFunc(content, systemPrompt, historyMsgs)
			if err != nil {
				logrus.Errorf("ChatSummaryHandler error %s", err)
			} else {
				ctx.Message.ReplyText(resp)
			}
		}()
		ctx.Abort()

	}

	return
}
