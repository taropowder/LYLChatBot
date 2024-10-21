package handlers

import (
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ChrysanthemumHandler struct {
}

func (h *ChrysanthemumHandler) Match(message *openwechat.Message) bool {
	if strings.Contains(message.Content, "提肛") {
		return true
	}
	return false
}

func (h *ChrysanthemumHandler) Helper(u *openwechat.User) string {
	return "您可无需艾特我直接输入如下指令：\n  1.今日提肛  (将会为您记录一次提肛)\n " + "2.提肛英雄榜 (将会为您展示今日提肛的英雄榜)\n"
}

func (h *ChrysanthemumHandler) Name() string {
	return "提肛模块"
}
func (h *ChrysanthemumHandler) Handle(ctx *openwechat.MessageContext) {

	db := database.GetDB()

	if ctx.Message.IsText() {

		content := ctx.Message.Content

		u, err := ctx.Message.Sender()
		logrus.Debugf("ChrysanthemumHandler %s", ctx.Message.MsgId)
		if err != nil {
			logrus.Errorf("log error %s", err)
			return
		}

		if group, success := u.AsGroup(); success {
			groupUser, err := ctx.Message.SenderInGroup()
			if err != nil {
				logrus.Error(err)
				return

			}

			if strings.Contains(content, "今日提肛") {

				replyText := fmt.Sprintf("%v 今日提肛次数 + 1", groupUser.NickName)

				record := model.LevatorAnusRecord{
					GroupId:  group.AvatarID(),
					NickName: groupUser.NickName,
				}

				db.Save(&record)

				_, err := ctx.Message.ReplyText(replyText)
				if err != nil {
					logrus.Error(err)
					return
				}
			} else if strings.Contains(content, "提肛英雄榜") {
				today := time.Now().Format("2006-01-02")

				// 按照 nickname 进行分组，并按照数量进行排序
				var queryResult []struct {
					NickName string
					Count    int
				}
				db.Model(&model.LevatorAnusRecord{}).Select("nick_name, COUNT(*) as count").
					Where("group_id = ? AND created_at >= ?", group.AvatarID(), today).
					Group("nick_name").Order("count desc").Scan(&queryResult)

				replyText := ""
				// 打印结果
				for _, r := range queryResult {
					replyText = replyText + fmt.Sprintf("群友: %s, 提肛次数: %d\n", r.NickName, r.Count)
				}
				_, err := ctx.Message.ReplyText(replyText)
				if err != nil {
					logrus.Error(err)
					return
				}
			}

		}
	}

}
