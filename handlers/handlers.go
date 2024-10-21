package handlers

import (
	"LYLChatBot/handlers/config"
	"github.com/eatmoreapple/openwechat"
)

type Handlers interface {
	Match(message *openwechat.Message) bool
	Handle(ctx *openwechat.MessageContext)
	Name() string
	Helper(u *openwechat.User) string
}

var MessagesHandlers = []Handlers{
	&HelperHandler{},
	&RecordHandler{},
	&config.GptConfigHandler{},

	&ImgHandler{},
	&NewsHandler{},
	&BiliBiliSummaryHandler{},

	&KnowledgeHandler{},
	&UrlSummaryHandler{},
	&ChatSummaryHandler{},

	&ScoreGGHandler{},
	&ChrysanthemumHandler{},

	&GptHandler{},
}
