package handlers

import (
	"LYLChatBot/constant"
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	"LYLChatBot/pkg/redis_conn"
	"LYLChatBot/utils"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

type UrlSummaryHandler struct {
}

func (h *UrlSummaryHandler) Match(message *openwechat.Message) bool {
	if message.IsArticle() {
		return true
	}

	if message.IsText() {
		re := regexp.MustCompile(`https?://\S+`)
		matches := re.FindString(message.Content)

		if matches != "" {
			return true
		}
	}

	return false
}

func (h *UrlSummaryHandler) Helper(u *openwechat.User) string {
	return "当您在群中或是与我的对话中存在微信文字链接，我会帮您分析其中的内容并进行总结。"
}

func (h *UrlSummaryHandler) Name() string {
	return "URL总结"
}
func (h *UrlSummaryHandler) Handle(ctx *openwechat.MessageContext) {

	logrus.Debugf("UrlSummary  handle %s", ctx.Message.MsgId)

	content := utils.RemoveAt(ctx.Message.Content)

	u, err := ctx.Message.Sender()
	if err != nil {
		logrus.Errorf("log error %s", err)
		return
	}

	gptFunc := utils.GetUserModuleFunc(u)
	systemPrompt := utils.GetRolePrompt(constant.ArticleContentPromptName, constant.ArticleContentDefaultPrompt)

	redisPool := redis_conn.RedisConnPool.Get()
	defer redisPool.Close()

	knowledgeType, err := redis.Bool(redisPool.Do("GET", fmt.Sprintf(constant.CacheKnowledgeKey, u.AvatarID())))
	if err == nil && knowledgeType {
		systemPrompt = utils.GetRolePrompt(constant.KnowledgePromptName, constant.KnowledgeDefaultPrompt)
		systemPrompt = fmt.Sprintf(systemPrompt, utils.GetAllKnowledgeTypes())
	} else {
		knowledgeType = false
	}

	messageUrl := ""

	if ctx.Message.IsArticle() {
		messageUrl = ctx.Message.Url
	} else {
		re := regexp.MustCompile(`https?://\S+`)
		matches := re.FindString(ctx.Content)

		if matches != "" {
			messageUrl = matches
			ctx.Abort()
		} else {
			return
		}
	}

	go func() {
		if strings.Contains(messageUrl, "weixin.qq.com/s?__biz=") {
			logrus.Debugf("url %s", messageUrl)
			a, err := utils.CollectWechatArticle(messageUrl)
			if err != nil {
				logrus.Errorf("CollectWechatArticle error %s", err)
				return
			}
			content = fmt.Sprintf("文章标题是：%s\n内容如下: %s\n", a.Title, a.TextContent)
		} else {
			htmlText, err := utils.RequestWithChrome(messageUrl)
			if err != nil {
				logrus.Errorf("CollectWechatArticle RequestWithChrome error %s", err)
				return
			}
			chinese := utils.GetHtmlText(htmlText)
			if chinese != "" {
				gptFunc = utils.GetUserModuleFunc(u)
				content = chinese
			} else {
				gptFunc = utils.NewBingGpt
				systemPrompt = utils.GetRolePrompt(constant.UrlPromptName, constant.UrlDefaultPrompt)
				content = fmt.Sprintf("我的链接是 %s ， 请你帮我总结这个链接的内容\n", messageUrl)
			}

		}

		historyMsgs := make([]utils.GptMessage, 0)

		resp, err := gptFunc(content, systemPrompt, historyMsgs)
		if err != nil {
			logrus.Errorf("UrlSummaryHandler error %s", err)
			ctx.Message.ReplyText("报错啦")
		} else {

			if knowledgeType {

				if strings.Contains(resp, "非预期文章") {
					systemPrompt = utils.GetRolePrompt(constant.ArticleContentPromptName, constant.ArticleContentDefaultPrompt)
					historyMsgs := make([]utils.GptMessage, 0)
					resp, err = gptFunc(content, systemPrompt, historyMsgs)
					if err == nil {
						ctx.Message.ReplyText(resp)
					}
					return
				}

				t, a := utils.GetKnowledgeFromResp(resp)
				if t != "" && a != "" {

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

					err = db.Where("knowledge_url = ?", messageUrl).First(&k).Error

					if errors.Is(err, gorm.ErrRecordNotFound) {
						k = model.KnowledgeRecord{
							KnowledgeType:     t,
							KnowledgeUrl:      messageUrl,
							KnowledgeAbstract: a,
							KnowledgeSharer:   sharer,
						}

						db.Save(&k)
						db.Create(&k) // 传递切片以插入多行数据

						qa := model.KnowledgeQARecord{
							KnowledgeId: k.ID,
							Question:    systemPrompt + "\n" + content,
							Answer:      resp,
							Questioner:  sharer,
							IsFirst:     true,
						}

						db.Save(&qa)
					}

					resp = fmt.Sprintf(constant.KnowledgeKeyFormat, k.ID) + resp

				}
			}

			ctx.Message.ReplyText(resp)
		}

	}()

	return
}
