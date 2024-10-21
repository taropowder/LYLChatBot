package model

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HandlerConfigRecord struct {
	gorm.Model

	UserId   string `gorm:"type:varchar(128);"`
	UserName string `gorm:"type:varchar(256);"`
	Config   string `gorm:"type:text;"`
}

type HandlerConfig struct {
	SystemRole string `json:"system_role"`
}

func (r *HandlerConfigRecord) GetConfig() HandlerConfig {

	if r.Config != "" {
		c := HandlerConfig{}
		err := json.Unmarshal([]byte(r.Config), &c)
		if err != nil {
			log.Error(err)
		} else {
			return c
		}
	}

	return HandlerConfig{SystemRole: ""}
}

func (r *HandlerConfigRecord) SetConfig(c HandlerConfig) {

	configStr, err := json.Marshal(c)
	if err == nil {
		r.Config = string(configStr)
	}

}
