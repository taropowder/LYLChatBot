package web

import (
	"LYLChatBot/web/router"
	log "github.com/sirupsen/logrus"
)

func RunServer(address string) {
	r := router.SetupRouter()
	if err := r.Run(address); err != nil {
		log.Infof("startup service failed, err:%v\n", err)
	}
}
