package handlers

import (
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	"LYLChatBot/pkg/redis_conn"
	"LYLChatBot/utils"
	"errors"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
)

type RecordHandler struct {
}

func (h *RecordHandler) Match(message *openwechat.Message) bool {
	return true
}

func (h *RecordHandler) Helper(u *openwechat.User) string {
	return ""
}

func (h *RecordHandler) Name() string {
	return ""
}

func (h *RecordHandler) Handle(ctx *openwechat.MessageContext) {

	db := database.GetDB()

	r := &model.MessageRecord{}

	if ctx.Message.IsText() {
		u, err := ctx.Message.Sender()
		logrus.Debugf("log hash %s", ctx.Message.MsgId)
		if err != nil {
			logrus.Errorf("log error %s", err)
			return
		}

		err = db.Where("message_id = ?", ctx.Message.MsgId).First(&r).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.MessageId = ctx.Message.MsgId
			if group, success := u.AsGroup(); success {
				logrus.Debugf("log group %s %s", group.NickName, group.AvatarID())
				r.GroupName = group.NickName
				r.GroupId = group.AvatarID()
				groupUser, err := ctx.Message.SenderInGroup()
				if err != nil {
					logrus.Errorf("SenderInGroup error %s", err)
				} else {
					r.UserName = groupUser.NickName
					r.NickName = groupUser.DisplayName
					//r.UserId = groupUser.AvatarID()
					//err := groupUser.Detail()
					//if err != nil {
					//	logrus.Errorf("user detail error %s", err)
					//}
					//if groupUser.AvatarID() != "" {
					//	r.UserName = groupUser.NickName
					//	r.NickName = groupUser.RemarkName
					//	r.UserId = groupUser.AvatarID()
					//}
					logrus.Debugf("log user %s %s", r.NickName, r.UserName)

					if err != nil {
						return
					}

				}
			} else {
				logrus.Debugf("log user %s %s", u.NickName, u.AvatarID())
				r.UserName = u.NickName
				r.UserId = u.AvatarID()

			}

			logrus.Debugf("log content %s", ctx.Content)
			logrus.Debugf("log is at %v", ctx.Message.IsAt())
			r.Content = ctx.Content
			db.Save(&r)
		}

	} else if ctx.Message.IsPicture() {
		redisPool := redis_conn.RedisConnPool.Get()
		defer redisPool.Close()
		k := utils.GetImgCacheKeyImageByPicture(ctx.Message)
		resp, _ := ctx.Message.GetPicture()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Error("读取响应失败:", err)
		}

		_, err = redisPool.Do("SET", k, bodyBytes)
		if err != nil {
			logrus.Error("读取响应失败:", err)
		}
	}

}
