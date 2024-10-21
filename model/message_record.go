package model

import "gorm.io/gorm"

type MessageRecord struct {
	//ID        uint       `gorm:"primary_key"`
	gorm.Model

	Content   string `gorm:"type:text;"`
	GroupId   string `gorm:"type:varchar(128);"`
	GroupName string `gorm:"type:varchar(128);"`
	UserId    string `gorm:"type:varchar(128);not null;"`
	UserName  string `gorm:"type:varchar(128);not null;"`
	NickName  string `gorm:"type:varchar(128);"`
	MessageId string `gorm:"type:varchar(128);not null;"`
}
