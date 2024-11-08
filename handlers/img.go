package handlers

import (
	"LYLChatBot/constant"
	"LYLChatBot/utils"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type ImgHandler struct {
}

func (h *ImgHandler) Match(message *openwechat.Message) bool {

	content := utils.RemoveAt(message.Content)

	if message.IsTickledMe() {
		return true
	}

	regex := regexp.MustCompile(constant.ImgKey)

	matches := regex.FindStringSubmatch(content)

	if message.IsText() && (matches != nil && len(matches) > 1) {
		if message.IsSendByGroup() {
			if message.IsAt() {
				return true
			}
		}

		return true
	}

	return false
}

func (h *ImgHandler) Helper(u *openwechat.User) string {
	return "拍一拍我，会随机掉落表情包哦"
}

func (h *ImgHandler) Name() string {
	return "斗图"
}

func (h *ImgHandler) Handle(ctx *openwechat.MessageContext) {

	logrus.Debugf("img handle %s", ctx.Message.MsgId)

	content := utils.RemoveAt(ctx.Message.Content)

	if ctx.IsTickledMe() {

		rand.Seed(time.Now().UnixNano())
		if rand.Intn(15) == 0 {

			_, err := ctx.Message.ReplyText("好好好，勤奋者的奖励！")
			rand.Seed(time.Now().UnixNano())
			start := rand.Intn(80) + 1
			img, err := utils.GetImagesByGoogleSearch(start, "性感美女")
			if err != nil {
				logrus.Error(err)
				return
			}

			go func() {
				_, err = ctx.Message.ReplyImage(img)
				if err != nil {
					logrus.Error("发送图片失败 , %v", err)
				}

			}()

		} else {
			go func() {
				//img, err := utils.GetDouTuImagesByParameters(0, "")
				img, err := utils.GetFaBiaoQingImagesByParameters(0, "")
				if err != nil {
					logrus.Error(err)
					return
				}

				_, err = ctx.Message.ReplyImage(img)
				if err != nil {
					logrus.Error("发送图片失败 , %v", err)
				}
			}()

		}

		logrus.Info("IsTickledMe")
		ctx.Abort()
		return
	}

	regex := regexp.MustCompile(constant.ImgKey)

	matches := regex.FindStringSubmatch(content)

	if matches != nil && len(matches) > 1 {

		//imgFun := utils.GetDouTuImagesByParameters
		imgFun := utils.GetFaBiaoQingImagesByParameters

		search := ""

		start := 0

		imgType := matches[1]

		if imgType == "表情包" {
			//imgFun = utils.GetDouTuImagesByParameters
			imgFun = utils.GetFaBiaoQingImagesByParameters
		} else if imgType == "图片" {
			imgFun = utils.GetImagesByGoogleSearch
		} else {
			return
		}

		search = matches[2]

		start = 1
		if matches[3] == "" {
			rand.Seed(time.Now().UnixNano())
			start = rand.Intn(10) + 1
		} else {
			var err error
			start, err = strconv.Atoi(matches[3])
			if err != nil {
				logrus.Error(err)
				return
			}

		}

		img, err := imgFun(start, search)
		if err != nil {
			logrus.Error(err)
			return
		}

		go func() {
			_, err = ctx.Message.ReplyImage(img)
			if err != nil {
				logrus.Error("发送图片失败, %v", err)
			}
		}()

		logrus.Infof("img %s %s", imgType, search)
		ctx.Abort()

	}

	return
}
