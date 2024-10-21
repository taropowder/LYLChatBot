package api

import (
	"LYLChatBot/conf"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type SendMessage struct {
	GroupId string `json:"group_id"`
	UserId  string `json:"user_id"`
	Content string `json:"content"`
	ImgUrl  string `json:"img_url"`
	FileUrl string `json:"file_url"`
}

var limitPool map[string][]int64

func Send(c *gin.Context) {

	authHeader := c.Request.Header.Get("Authorization")

	if limitPool == nil {
		limitPool = make(map[string][]int64, 0)
	}

	now := time.Now().Unix()

	if _, ok := limitPool[authHeader]; ok {
		if len(limitPool[authHeader]) < 5 {
			limitPool[authHeader] = append(limitPool[authHeader], now)
		} else {
			limitPool[authHeader] = append(limitPool[authHeader], now)
			limitPool[authHeader] = limitPool[authHeader][1:]
			oldestReq := limitPool[authHeader][0]
			if (now - oldestReq) < 60*30 {

				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "failure",
					"data":   "speed limit",
				})
				return

			}
		}
	} else {
		limitPool[authHeader] = make([]int64, 0)
		limitPool[authHeader] = append(limitPool[authHeader], now)
	}

	sendData := SendMessage{}

	if err := c.BindJSON(&sendData); err != nil {
		log.Errorf("error in BindJSON %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failure",
			"data":   err.Error(),
		})
		return
	}

	self, err := conf.BotInstance.GetCurrentUser()
	if err != nil {
		log.Error(err)
		return
	}

	if sendData.UserId != "" {

		friends, err := self.Friends()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failure",
				"data":   err.Error(),
			})
			return

		}

		findFriends := friends.SearchByID(sendData.UserId)
		if len(findFriends) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failure",
				"data":   fmt.Sprintf("no such uid %v", sendData.GroupId),
			})
			return
		} else {
			err := findFriends.SendText(sendData.Content, time.Second*5)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "failure",
					"data":   err.Error(),
				})
				return
			}
		}

	}

	if sendData.GroupId != "" {

		groups, err := self.Groups()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failure",
				"data":   err,
			})
			return

		}

		findGroups := groups.SearchByID(sendData.GroupId)
		if len(findGroups) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failure",
				"data":   fmt.Sprintf("no such gid %v", sendData.GroupId),
			})
			return
		} else {
			err := findGroups.SendText(sendData.Content, time.Second*5)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": "failure",
					"data":   err,
				})
				return
			}
		}
	}

	//c.String(http.StatusOK, string(content))
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   sendData.Content,
	})

}
