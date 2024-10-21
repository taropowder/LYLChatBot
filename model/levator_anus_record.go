package model

import "gorm.io/gorm"

type LevatorAnusRecord struct {
	gorm.Model

	GroupId  string `gorm:"type:varchar(128);"`
	NickName string `gorm:"type:varchar(128);"`
}
