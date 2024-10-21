package task

import (
	"LYLChatBot/model"
	"LYLChatBot/pkg/database"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.Infof("Starting Cron Jobs")

	c := cron.New()
	err := c.AddFunc("*/5 * * * *", taskRunner)
	if err != nil {
		log.Error(err)
	}

	c.Start()
}

type baseTaskConfig struct {
	GroupId       string `json:"group_id"`
	NextTimeStamp int    `json:"next_time_stamp"`
	Role          string `json:"role"`

	// 间隔秒数
	CheckInterval int `json:"check_interval"`
	SpeakInterval int `json:"speak_interval"`
}

func taskRunner() {

	db := database.GetDB()
	tasks := make([]model.CronTaskRecord, 0)

	db.Where("status = true").Find(&tasks)

	for _, task := range tasks {
		//log.Debugf("run task %v", task.TaskType)

		if task.TaskType == "LOL" {
			reportLiveTextTask(db, task)
		} else if task.TaskType == "Remind" {
			remindRobot(db, task)
		} else if task.TaskType == "Treat" {
			ShitStirrer(db, task)
		} else if task.TaskType == "Hot" {
			hotNews(db, task)
		}

	}

}
