package handlers

import (
	"LYLChatBot/constant"
	"LYLChatBot/utils"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type NewsHandler struct {
}

func (h *NewsHandler) Match(message *openwechat.Message) bool {

	content := utils.RemoveAt(message.Content)

	regex := regexp.MustCompile(constant.HotNewsTypeKey)

	matches := regex.FindStringSubmatch(content)

	if message.IsText() && (strings.HasPrefix(content, constant.HotNewsKey) ||
		strings.HasPrefix(content, constant.RandomHotNewsKey) ||
		(matches != nil && len(matches) > 0)) {
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

func (h *NewsHandler) Helper(u *openwechat.User) string {
	return fmt.Sprintf("您可以输入\"您可艾特我后输入如下指令：\n1.今日新闻(将会为您拉取微博、B站、知乎的热点)\n"+
		"2.随便看看(将会为随机您拉取一些最近的热点新闻)"+
		"3.$news热榜 ($news 可选项：%v)", utils.NewsTypes)
}

func (h *NewsHandler) Name() string {
	return "热点新闻"
}

func (h *NewsHandler) Handle(ctx *openwechat.MessageContext) {

	log.Debugf("news handle %s", ctx.Message.MsgId)

	content := utils.RemoveAt(ctx.Message.Content)

	if content == constant.RandomHotNewsKey {
		ctx.Message.ReplyText(utils.RandomNews(3, 5))
		ctx.Abort()
	} else if content == constant.HotNewsKey {
		ctx.Message.ReplyText(utils.GetNewsByTypes([]string{"zhihu", "bilibili", "weibo"}, 8))
		ctx.Abort()
	}

	regex := regexp.MustCompile(constant.HotNewsTypeKey)

	matches := regex.FindStringSubmatch(content)

	if matches != nil && len(matches) > 0 {
		newsType := matches[1]
		newsResp, err := utils.GetHostNews(newsType)
		if err != nil {
			log.Error(err)
			return
		}
		s := utils.FormatHostNews(newsResp, newsType, 10)
		ctx.Message.ReplyText(s)
		ctx.Abort()
	}

}
