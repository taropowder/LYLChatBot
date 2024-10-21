package model

import "gorm.io/gorm"

type CronTaskRecord struct {
	gorm.Model

	TaskType   string `gorm:"type:varchar(128);"`
	TaskConfig string `gorm:"type:text;"`
	Status     bool   `gorm:"type:boolean;"`
}
