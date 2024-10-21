package handlers

import (
	"LYLChatBot/constant"
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	"LYLChatBot/task"
	"LYLChatBot/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

type ScoreGGHandler struct {
}

func (h *ScoreGGHandler) Match(message *openwechat.Message) bool {
	if message.IsText() {
		if strings.Contains(message.Content, constant.StopLiveKey) ||
			strings.Contains(message.Content, constant.StopLiveKey) ||
			strings.Contains(message.Content, constant.MatchesQueryKey) {
			return true
		}
		regex := regexp.MustCompile(constant.LiveReportKeyReg)

		matches := regex.FindStringSubmatch(message.Content)

		if matches != nil {
			return true
		}

		return false
	}
	return false
}

func (h *ScoreGGHandler) Helper(u *openwechat.User) string {
	return "您可艾特我后输入如下指令：\n  1.近期赛事  (将会为您拉取最近一周的赛事,赛事后面的数字为matchId) \n " +
		"2.文字直播:$matchId (将会为您持续播报比赛)\n" +
		"3.停止直播 (将会关闭当前群的直播)\n"
}

func (h *ScoreGGHandler) Name() string {
	return "LOL比赛播报"
}

func (h *ScoreGGHandler) Handle(ctx *openwechat.MessageContext) {
	//content := utils.FormatMessage(ctx.Message.Content, conf.ConfigureInstance.Handlers.Qwen.KeyWords)

	logrus.Debugf("ScoreGG handle %s", ctx.Message.MsgId)

	regex := regexp.MustCompile(constant.LiveReportKeyReg)

	matches := regex.FindStringSubmatch(ctx.Message.Content)

	if ctx.Message.IsSendByGroup() {
		if ctx.Message.IsAt() {

			if len(matches) == 2 {
				matchId := matches[1]

				logrus.Debugf("start report %v", matchId)
				db := database.GetDB()

				record := model.CronTaskRecord{}

				now := time.Now().Unix()

				group, err := ctx.Message.Sender()
				if err != nil {
					logrus.Errorf("SenderInGroup error %s", err)
				} else {

					config := task.ScoreGGLiveTextTaskConfig{}

					old_records := make([]model.CronTaskRecord, 0)

					err = db.Where("task_type = 'LOL' and  task_config LIKE ? ", fmt.Sprintf("%%%s%%", group.AvatarID())).Find(&old_records).Error

					if !errors.Is(err, gorm.ErrRecordNotFound) {

						for _, od := range old_records {
							config := task.ScoreGGLiveTextTaskConfig{}
							err := json.Unmarshal([]byte(od.TaskConfig), &config)
							if err != nil {
								logrus.Error(err)
							}
							if int(now) < config.EndTimestamp {
								_, err = ctx.Message.ReplyText(fmt.Sprintf("本群已有直播 %s 进行，请先关闭之前直播", config.MatchId))
								if err != nil {
									return
								}
								return
							}
						}

						err = db.Where("task_type = 'LOL' and task_config LIKE ? ", fmt.Sprintf("%%%s%%", matchId)).First(&record).Error

						if errors.Is(err, gorm.ErrRecordNotFound) {
							//	新建
							config.LastKey = "1"
							config.MatchId = matchId
							config.Interval = 30
							config.GroupId = []string{group.AvatarID()}
							config.EndTimestamp = int(now + 60*60)
						} else {
							err := json.Unmarshal([]byte(record.TaskConfig), &config)
							if err != nil {
								logrus.Error(err)
							}
							if !strings.Contains(strings.Join(config.GroupId, " "), group.AvatarID()) {
								config.GroupId = append(config.GroupId, group.AvatarID())
							}
							config.EndTimestamp = int(now + 60*60)
						}

						newConfigStr, err := json.Marshal(config)
						if err != nil {
							logrus.Error(err)
							return
						}
						record.TaskConfig = string(newConfigStr)
						record.Status = true
						record.TaskType = "LOL"
						db.Save(&record)

						_, err = ctx.Message.ReplyText(fmt.Sprintf("已经开始播报 %s 比赛文字直播", matchId))
						if err != nil {
							return
						}

					}
					ctx.Abort()

				}
			} else if strings.Contains(ctx.Message.Content, constant.MatchesQueryKey) {

				logrus.Debugf("start get all matches %s %s", ctx.Message.Content, constant.MatchesQueryKey)

				g := utils.NewGames()
				//fmt.Println(g.GetGamesByTournamentId("629"))
				if strings.Contains(ctx.Message.Content, "lpl") {
					_, err := ctx.Message.ReplyText(g.GetGamesByTournamentId("629"))
					if err != nil {
						return
					}
				}
				if strings.Contains(ctx.Message.Content, "msi") {
					_, err := ctx.Message.ReplyText(g.GetGamesByTournamentId("697"))
					if err != nil {
						return
					}
				}
				ctx.Abort()

			} else if strings.Contains(ctx.Message.Content, constant.StopLiveKey) {
				db := database.GetDB()
				tasks := make([]model.CronTaskRecord, 0)
				group, err := ctx.Message.Sender()
				if err != nil {
					logrus.Errorf("SenderInGroup error %s", err)
				} else {
					err := db.Where("task_type = 'LOL' and task_config LIKE ? ", fmt.Sprintf("%%%s%%", group.AvatarID())).Find(&tasks).Error
					if errors.Is(err, gorm.ErrRecordNotFound) {
						_, err = ctx.Message.ReplyText("本群无直播进行")
					} else {

						_, err = ctx.Message.ReplyText(fmt.Sprintf("已关闭直播"))

						now := time.Now().Unix()

						for _, tk := range tasks {
							tk.Status = false
							config := task.ScoreGGLiveTextTaskConfig{}
							err := json.Unmarshal([]byte(tk.TaskConfig), &config)
							if err != nil {
								logrus.Error(err)
							}
							config.EndTimestamp = int(now)
							newConfigStr, err := json.Marshal(config)
							if err != nil {
								logrus.Error(err)
								continue
							}
							tk.TaskConfig = string(newConfigStr)
							logrus.Debugf("new config %v", string(newConfigStr))
							db.Save(&tk)
						}
					}
				}
				ctx.Abort()
			}

		}
	}
}
